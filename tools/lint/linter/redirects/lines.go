package redirects

import (
	"bytes"
	"encoding/json"
)

func lineLookup(docsJSONText string, count int) []int {
	result := make([]int, count)
	for i := range result {
		result[i] = 1
	}
	if count == 0 {
		return result
	}

	dec := json.NewDecoder(bytes.NewReader([]byte(docsJSONText)))
	if !advanceToRedirectsArray(dec) {
		return result
	}
	scanRedirectEntries(dec, docsJSONText, result, count)
	return result
}

func advanceToRedirectsArray(dec *json.Decoder) bool {
	for {
		tok, err := dec.Token()
		if err != nil {
			return false
		}
		key, ok := tok.(string)
		if !ok || key != "redirects" {
			continue
		}
		next, err := dec.Token()
		if err != nil {
			return false
		}
		if delim, ok := next.(json.Delim); ok && delim == '[' {
			return true
		}
	}
}

func scanRedirectEntries(
	dec *json.Decoder,
	docsJSONText string,
	result []int,
	count int,
) {
	arrayDepth := 1
	entryIdx := 0
	for {
		offsetBefore := dec.InputOffset()
		tok, err := dec.Token()
		if err != nil {
			return
		}
		delta := depthDelta(tok)
		if delta >= 0 && arrayDepth == 1 && entryIdx < count {
			result[entryIdx] = lineForOffset(
				docsJSONText,
				tokenStart(docsJSONText, int(offsetBefore)),
			)
			entryIdx++
		}
		arrayDepth += delta
		if arrayDepth == 0 {
			return
		}
	}
}

func depthDelta(tok json.Token) int {
	switch tok {
	case json.Delim('{'), json.Delim('['):
		return 1
	case json.Delim('}'), json.Delim(']'):
		return -1
	}
	return 0
}

func tokenStart(text string, offset int) int {
	for offset < len(text) {
		char := text[offset]
		if char == ' ' || char == '\t' || char == '\n' || char == '\r' || char == ',' {
			offset++
			continue
		}
		break
	}
	return offset
}

func lineForOffset(text string, byteOffset int) int {
	line := 1
	for i := 0; i < byteOffset; i++ {
		if text[i] == '\n' {
			line++
		}
	}
	return line
}
