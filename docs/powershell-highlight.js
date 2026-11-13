// Colors native commands in PowerShell code blocks, such as miru, icacls and
// aws, the way the bash grammar colors them in Unix blocks. Mintlify's
// PowerShell grammar leaves native commands and their flags uncolored, and its
// custom-language support does not load on hosted sites
// (https://github.com/orgs/mintlify/discussions/3401).
(() => {
  // [light, dark] colors the code block themes pinned in docs.json give the
  // matching bash tokens.
  const COLORS = {
    plain: ["#1F2328", "#D4D4D4"],
    command: ["#953800", "#DCDCAA"],
    argument: ["#0A3069", "#CE9178"],
    option: ["#0550AE", "#569CD6"],
    variable: ["#1F2328", "#9CDCFE"],
    continuation: ["#CF222E", "#D7BA7D"],
  };
  // A line that runs a native command, optionally assigning its output. Cmdlets
  // (Verb-Noun) and keywords are left to Mintlify's grammar, which colors them.
  const COMMAND_LINE =
    /^\s*(\$[\w:]+\s+=\s+)?(?!(if|elseif|else|foreach|for|while|do|switch|function|param|return|try|catch|finally|throw)\b)[A-Za-z][\w.]*(\s|$)/;
  const TOKEN =
    /(\s+)|(`(?=\s*$))|((?<=\s)=(?=\s)|\|)|(--?[A-Za-z][\w-]*)|("[^"]*"|'[^']*')|(\$(?:env:)?\w+)|([^\s`"'$|]+)|(.)/g;

  function tokenize(text, startsCommand) {
    const tokens = [];
    let isFirst = startsCommand;
    for (const match of text.matchAll(TOKEN)) {
      const [tok, space, cont, operator, option, string, variable, word] = match;
      let kind = "plain";
      if (cont) kind = "continuation";
      else if (option) kind = "option";
      else if (string) kind = "argument";
      else if (variable) kind = "variable";
      else if (word) kind = isFirst ? "command" : "argument";
      tokens.push([tok, kind]);
      // The word after "=" or "|" is the next command, as in bash.
      if (operator) isFirst = true;
      else if (!space) isFirst = false;
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
      const startsCommand = COMMAND_LINE.test(text);
      const isCommand = startsCommand || inCommand;
      if (isCommand && !line.dataset.highlighted) {
        line.replaceChildren(...tokenize(text, startsCommand).map(toSpan));
        line.dataset.highlighted = "true";
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
