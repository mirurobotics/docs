// Colors native commands in shell code blocks, such as miru, icacls and aws,
// the way Shiki's bash grammar colors them. Mintlify's PowerShell grammar
// leaves native commands and their flags uncolored, and its hosted bash
// highlighting leaves some blocks uncolored too. Its custom-language support
// does not load on hosted sites (https://github.com/orgs/mintlify/discussions/3401).
(() => {
  // [light, dark] colors Mintlify's default code themes (github-light-default,
  // dark-plus) give the matching bash tokens.
  const COLORS = {
    plain: ["#1F2328", "#D4D4D4"],
    command: ["#953800", "#DCDCAA"],
    cmdlet: ["#0550AE", "#DCDCAA"],
    argument: ["#0A3069", "#CE9178"],
    option: ["#0550AE", "#569CD6"],
    variable: ["#1F2328", "#9CDCFE"],
    continuation: ["#CF222E", "#D7BA7D"],
  };
  // A line that runs a native command, after an optional "$ " prompt or an
  // assignment of its output. Lines that start with a cmdlet (Verb-Noun) or
  // keyword are left to Mintlify's grammar, which colors them.
  const COMMAND_LINE =
    /^\s*(\$\s+|\$[\w:]+\s+=\s+)?(?!(if|elif|elseif|else|then|fi|foreach|for|in|while|do|done|until|case|esac|switch|function|filter|param|begin|process|end|return|exit|break|continue|try|catch|finally|throw|trap|using)\b)[A-Za-z][\w.]*(\s|$)/i;

  const CMDLET = /^[A-Za-z]+-[A-Za-z]+$/;
  const VARIABLE = /(\$\{?(?:env:)?\w+\}?)/;

  // Each shell's line continuation character, by Mintlify's code language.
  const SHELLS = {
    powershell: { cont: "`", token: tokenPattern("`"), hasCmdlets: true },
    shellscript: { cont: "\\", token: tokenPattern("\\\\"), hasCmdlets: false },
  };

  function tokenPattern(cont) {
    return new RegExp(
      `(\\s+)|(${cont}(?=\\s*$))|((?<=\\s)=(?=\\s)|\\|)|(--?[A-Za-z][\\w-]*)` +
        `|("[^"]*"|'[^']*')|${VARIABLE.source}|([^\\s${cont}"'$|]+)|(.)`,
      "g",
    );
  }

  function tokenize(text, startsCommand, shell) {
    const tokens = [];
    let isFirst = startsCommand;
    for (const match of text.matchAll(shell.token)) {
      const [tok, space, cont, operator, option, string, variable, word] = match;
      let kind = "plain";
      if (cont) kind = "continuation";
      else if (option) kind = "option";
      else if (variable) kind = "variable";
      else if (word && isFirst && shell.hasCmdlets && CMDLET.test(word)) kind = "cmdlet";
      else if (word) kind = isFirst ? "command" : "argument";
      // A "$ " prompt takes the command color, as in Shiki's bash grammar.
      else if (tok === "$" && isFirst && tokens.every(([, k]) => k === "plain"))
        kind = "command";
      if (string) tokens.push(...stringTokens(string));
      else tokens.push([tok, kind]);
      // The word after "=" or "|" is the next command, as in bash.
      if (operator) isFirst = true;
      else if (!space) isFirst = false;
    }
    return tokens;
  }

  // Variables inside double quotes keep their color, as in both grammars.
  function stringTokens(string) {
    if (string.startsWith("'")) return [[string, "argument"]];
    return string
      .split(VARIABLE)
      .filter(Boolean)
      .map((part) => [part, VARIABLE.test(part) ? "variable" : "argument"]);
  }

  function toSpan([text, kind]) {
    const span = document.createElement("span");
    span.style.color = COLORS[kind][0];
    span.style.setProperty("--shiki-dark", COLORS[kind][1]);
    span.textContent = text;
    return span;
  }

  // Mintlify's bash coloring is right on most blocks; only redo the blocks it
  // left without a command or flag color.
  function isUncolored(code) {
    const colors = [COLORS.command[1], COLORS.option[1]];
    return ![...code.querySelectorAll("span[style]")].some((span) =>
      colors.includes(span.style.getPropertyValue("--shiki-dark").toUpperCase()),
    );
  }

  function highlight(code) {
    const shell = SHELLS[code.getAttribute("language")];
    if (shell === SHELLS.shellscript && !isUncolored(code)) return;
    const endsContinued = new RegExp(`[^;\\s]\\s+\\${shell.cont}\\s*$`);
    let inCommand = false;
    for (const line of code.querySelectorAll(".line")) {
      const text = line.textContent;
      const startsCommand = COMMAND_LINE.test(text);
      const isCommand = startsCommand || inCommand;
      if (isCommand && !line.dataset.highlighted) {
        line.replaceChildren(...tokenize(text, startsCommand, shell).map(toSpan));
        line.dataset.highlighted = "true";
      }
      // A continuation after ";" starts a new statement, not a continued command.
      inCommand = isCommand && endsContinued.test(text);
    }
  }

  // Mintlify renders pages and tabs client-side, so recheck after DOM changes.
  let isQueued = false;
  function schedule() {
    if (isQueued) return;
    isQueued = true;
    requestAnimationFrame(() => {
      isQueued = false;
      const selector = Object.keys(SHELLS)
        .map((lang) => `code[language="${lang}"]`)
        .join(",");
      document.querySelectorAll(selector).forEach(highlight);
    });
  }
  new MutationObserver(schedule).observe(document.documentElement, {
    childList: true,
    subtree: true,
  });
  schedule();
})();
