// Colors miru commands in PowerShell code blocks the way the bash grammar
// colors them in Unix blocks. Mintlify's PowerShell grammar leaves native
// commands and their flags uncolored, and its custom-language support does
// not load on hosted sites (https://github.com/orgs/mintlify/discussions/3401).
(() => {
  // [light, dark] colors Mintlify's themes give the matching bash tokens.
  const COLORS = {
    plain: ["#1F2328", "#D4D4D4"],
    command: ["#953800", "#DCDCAA"],
    argument: ["#0A3069", "#CE9178"],
    option: ["#0550AE", "#569CD6"],
    variable: ["#1F2328", "#9CDCFE"],
    continuation: ["#CF222E", "#D7BA7D"],
  };
  const TOKEN =
    /(\s+)|(`(?=\s*$))|(--?[A-Za-z][\w-]*)|("[^"]*"|'[^']*')|(\$(?:env:)?\w+)|([^\s`"'$]+)|(.)/g;

  function tokenize(text, startsCommand) {
    const tokens = [];
    let isFirst = startsCommand;
    for (const [tok, space, cont, option, string, variable, word] of text.matchAll(TOKEN)) {
      let kind = "plain";
      if (cont) kind = "continuation";
      else if (option) kind = "option";
      else if (string) kind = "argument";
      else if (variable) kind = "variable";
      else if (word) kind = isFirst ? "command" : "argument";
      tokens.push([tok, kind]);
      if (!space) isFirst = false;
    }
    return tokens;
  }

  function toSpan([text, kind]) {
    const span = document.createElement("span");
    span.style.color = COLORS[kind][0];
    span.style.setProperty("--shiki-dark", COLORS[kind][1]);
    span.textContent = text;
    return span;
  }

  function highlight(code) {
    let inCommand = false;
    for (const line of code.querySelectorAll(".line")) {
      const text = line.textContent;
      const startsCommand = /^\s*miru(\s|$)/.test(text);
      const isCommand = startsCommand || inCommand;
      if (isCommand && !line.dataset.miruHighlighted) {
        line.replaceChildren(...tokenize(text, startsCommand).map(toSpan));
        line.dataset.miruHighlighted = "true";
      }
      inCommand = isCommand && /\s`\s*$/.test(text);
    }
  }

  // Mintlify renders pages and tabs client-side, so recheck after DOM changes.
  let isQueued = false;
  function schedule() {
    if (isQueued) return;
    isQueued = true;
    requestAnimationFrame(() => {
      isQueued = false;
      document.querySelectorAll('code[language="powershell"]').forEach(highlight);
    });
  }
  new MutationObserver(schedule).observe(document.documentElement, {
    childList: true,
    subtree: true,
  });
  schedule();
})();
