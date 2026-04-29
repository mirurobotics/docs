#!/usr/bin/env node
// Validate `redirects` in docs.json against the on-disk docs/ tree.
//
// Mintlify serves docs/foo/bar.mdx at URL /docs/foo/bar. The `redirects`
// array in docs.json rewrites URLs at the edge. This script catches:
//   - Dead redirects (source already serves a real page).
//   - Missing destinations (destination has no real page).
//   - Bad prefixes / unsupported schemes / malformed paths.
//
// Honors DOCS_LINT_ROOT (defaults to the repo root one level above this
// script). Reads ${root}/docs.json, resolves files under ${root}/docs/.
//
// Diagnostics are emitted on stderr in the form
//   docs.json:<line>: redirects[<i>] <field> "<value>": <message>
// matching the file:line:col: style emitted by tools/lint/main.go.

import fs from 'node:fs';
import path from 'node:path';
import process from 'node:process';
import { fileURLToPath } from 'node:url';

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const root = process.env.DOCS_LINT_ROOT || path.resolve(__dirname, '..');
const docsJsonPath = path.join(root, 'docs.json');

const WILDCARD_SEGMENT = /^:[A-Za-z][A-Za-z0-9]*\*?$/;

let warnedOnLineLookup = false;

function warnOnce(message) {
	if (warnedOnLineLookup) return;
	warnedOnLineLookup = true;
	process.stderr.write(`warning: ${message}\n`);
}

function readDocsJson() {
	let text;
	try {
		text = fs.readFileSync(docsJsonPath, 'utf8');
	} catch (err) {
		if (err && err.code === 'ENOENT') {
			process.stdout.write(`No docs.json at ${root}; nothing to check\n`);
			process.exit(0);
		}
		process.stderr.write(`docs.json: read error: ${err.message}\n`);
		process.exit(1);
	}

	let parsed;
	try {
		parsed = JSON.parse(text);
	} catch (err) {
		process.stderr.write(`docs.json: invalid JSON: ${err.message}\n`);
		process.exit(1);
	}
	return { text, parsed };
}

// Locate the 1-based line number of the n-th `"source":` literal inside
// the `"redirects"` array. Returns null if anchoring fails.
function buildSourceLineLookup(text, redirectCount) {
	const arrayKey = /"redirects"\s*:\s*\[/g;
	const arrayMatch = arrayKey.exec(text);
	if (!arrayMatch) return null;

	// Offset of the opening '['.
	const arrayStart = arrayMatch.index + arrayMatch[0].length - 1;
	const lookup = new Array(redirectCount).fill(null);

	const sourcePattern = /"source"\s*:/g;
	sourcePattern.lastIndex = arrayStart;

	let i = 0;
	let match;
	while (i < redirectCount && (match = sourcePattern.exec(text)) !== null) {
		// Compute 1-based line number of match.index.
		let line = 1;
		for (let j = 0; j < match.index; j += 1) {
			if (text.charCodeAt(j) === 10) line += 1;
		}
		lookup[i] = line;
		i += 1;
	}
	return lookup;
}

function lineFor(lookup, index) {
	if (!lookup) return null;
	const v = lookup[index];
	return typeof v === 'number' ? v : null;
}

function diag(line, index, field, value, message) {
	const prefix = line === null
		? `docs.json:`
		: `docs.json:${line}:`;
	const valueStr = value === undefined || value === null ? '' : String(value);
	process.stderr.write(
		`${prefix} redirects[${index}] ${field} "${valueStr}": ${message}\n`,
	);
}

// Strip leading '/', trailing '/', '?...', '#...'. Returns the cleaned path.
function cleanPath(p) {
	let s = p;
	const queryIdx = s.indexOf('?');
	if (queryIdx !== -1) s = s.slice(0, queryIdx);
	const hashIdx = s.indexOf('#');
	if (hashIdx !== -1) s = s.slice(0, hashIdx);
	if (s.startsWith('/')) s = s.slice(1);
	if (s.endsWith('/')) s = s.slice(0, -1);
	return s;
}

// Walk the segments and return the prefix array (segments before the
// first wildcard segment). If no wildcard, returns all segments.
function splitWildcard(cleaned) {
	const segments = cleaned.split('/').filter((s) => s.length > 0);
	const prefix = [];
	let hasWildcard = false;
	for (const seg of segments) {
		if (WILDCARD_SEGMENT.test(seg)) {
			hasWildcard = true;
			break;
		}
		prefix.push(seg);
	}
	return { prefix, hasWildcard };
}

function fileExists(p) {
	try {
		const st = fs.statSync(p);
		return st.isFile();
	} catch {
		return false;
	}
}

function dirExists(p) {
	try {
		const st = fs.statSync(p);
		return st.isDirectory();
	} catch {
		return false;
	}
}

// Recursively check whether a directory contains any .mdx or .md file.
function dirHasPages(dir) {
	let entries;
	try {
		entries = fs.readdirSync(dir, { withFileTypes: true });
	} catch {
		return false;
	}
	for (const ent of entries) {
		const full = path.join(dir, ent.name);
		if (ent.isFile()) {
			if (ent.name.endsWith('.mdx') || ent.name.endsWith('.md')) {
				return true;
			}
		} else if (ent.isDirectory()) {
			if (dirHasPages(full)) return true;
		}
	}
	return false;
}

function validateSource(prefixFs, hasWildcard) {
	if (!hasWildcard) {
		if (fileExists(`${prefixFs}.mdx`) || fileExists(`${prefixFs}.md`)) {
			return 'dead redirect (source resolves to a real page)';
		}
		return null;
	}
	// Wildcard source: prefix must not be a real page or contain pages.
	if (fileExists(`${prefixFs}.mdx`) || fileExists(`${prefixFs}.md`)) {
		return 'dead redirect (wildcard source prefix resolves to a real page)';
	}
	if (dirExists(prefixFs) && dirHasPages(prefixFs)) {
		return 'dead redirect (wildcard source prefix has real pages)';
	}
	return null;
}

// Collect all OpenAPI source paths declared in docs.json's nav. Mintlify
// generates pages from these yaml files at build time, so a wildcard
// destination may target a virtual directory that only exists at build
// time. We accept `${prefix}.yaml` as a valid destination prefix when it
// is referenced as an openapi.source somewhere in the nav.
function collectOpenApiSources(parsed) {
	const sources = new Set();
	const walk = (node) => {
		if (!node || typeof node !== 'object') return;
		if (Array.isArray(node)) {
			for (const item of node) walk(item);
			return;
		}
		if (node.openapi && typeof node.openapi === 'object'
			&& typeof node.openapi.source === 'string') {
			sources.add(node.openapi.source);
		}
		for (const key of Object.keys(node)) {
			walk(node[key]);
		}
	};
	walk(parsed);
	return sources;
}

function validateDestination(prefixFs, hasWildcard, prefixRel, openApiSources) {
	if (!hasWildcard) {
		if (fileExists(`${prefixFs}.mdx`) || fileExists(`${prefixFs}.md`)) {
			return null;
		}
		return 'missing destination (no .mdx or .md page exists)';
	}
	if (dirExists(prefixFs)) return null;
	// Mintlify-generated OpenAPI routes: accept ${prefix}.yaml when the
	// yaml is registered as a nav openapi source.
	const yamlRel = `${prefixRel}.yaml`;
	if (openApiSources.has(yamlRel) && fileExists(`${prefixFs}.yaml`)) {
		return null;
	}
	return 'wildcard prefix not a directory';
}

function main() {
	const { text, parsed } = readDocsJson();
	const redirects = parsed && Array.isArray(parsed.redirects) ? parsed.redirects : [];

	if (redirects.length === 0) {
		process.stdout.write('Checked 0 redirects: OK\n');
		process.exit(0);
	}

	const lineLookup = buildSourceLineLookup(text, redirects.length);
	if (lineLookup === null) {
		warnOnce('could not anchor "redirects" array; falling back to indexless diagnostics');
	}
	const openApiSources = collectOpenApiSources(parsed);

	let violations = 0;
	const reportLine = (i) => {
		const ln = lineFor(lineLookup, i);
		if (ln === null) {
			warnOnce(`could not locate "source" for redirects[${i}]; falling back to indexless diagnostic`);
		}
		return ln;
	};

	for (let i = 0; i < redirects.length; i += 1) {
		const entry = redirects[i];
		const line = reportLine(i);

		if (!entry || typeof entry !== 'object') {
			diag(line, i, 'entry', '', 'not an object');
			violations += 1;
			continue;
		}

		const { source, destination } = entry;

		// Rule (a): both source and destination must be non-empty strings.
		const sourceOk = typeof source === 'string' && source.length > 0;
		const destOk = typeof destination === 'string' && destination.length > 0;
		if (!sourceOk) {
			diag(line, i, 'source', source ?? '', 'must be a non-empty string');
			violations += 1;
		}
		if (!destOk) {
			diag(line, i, 'destination', destination ?? '', 'must be a non-empty string');
			violations += 1;
		}
		if (!sourceOk || !destOk) continue;

		// Rule (b): source must start with '/'; destination must start with
		// '/' or http(s)://.
		if (!source.startsWith('/')) {
			diag(line, i, 'source', source, "bad path: must start with '/'");
			violations += 1;
		}

		const destIsHttp =
			destination.startsWith('http://') || destination.startsWith('https://');
		if (!destination.startsWith('/') && !destIsHttp) {
			diag(line, i, 'destination', destination, "bad path: must start with '/' (or http(s)://)");
			violations += 1;
		}

		// Validate source filesystem rules if it has the right shape.
		if (source.startsWith('/')) {
			const cleaned = cleanPath(source);
			if (!cleaned.startsWith('docs/') && cleaned !== 'docs') {
				diag(line, i, 'source', source, 'bad prefix (must start with /docs/)');
				violations += 1;
			} else {
				const { prefix, hasWildcard } = splitWildcard(cleaned);
				const prefixFs = path.join(root, prefix.join('/'));
				const err = validateSource(prefixFs, hasWildcard);
				if (err) {
					diag(line, i, 'source', source, err);
					violations += 1;
				}
			}
		}

		// Validate destination filesystem rules unless it's external.
		if (!destIsHttp && destination.startsWith('/')) {
			const cleaned = cleanPath(destination);
			if (!cleaned.startsWith('docs/') && cleaned !== 'docs') {
				diag(line, i, 'destination', destination, 'bad prefix (must start with /docs/)');
				violations += 1;
			} else {
				const { prefix, hasWildcard } = splitWildcard(cleaned);
				const prefixRel = prefix.join('/');
				const prefixFs = path.join(root, prefixRel);
				const err = validateDestination(prefixFs, hasWildcard, prefixRel, openApiSources);
				if (err) {
					diag(line, i, 'destination', destination, err);
					violations += 1;
				}
			}
		}
	}

	if (violations > 0) {
		process.exit(1);
	}
	process.stdout.write(`Checked ${redirects.length} redirects: OK\n`);
	process.exit(0);
}

main();
