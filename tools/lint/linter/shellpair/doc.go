// Package shellpair keeps the Bash and PowerShell versions of a
// multi-line command correct and in step. A bash block must not continue
// lines with PowerShell's trailing "`", and a powershell block must not
// continue them with a trailing "\". Inside a <CodeGroup>, a fenced block
// titled "Bash" that continues lines with "\" needs a block titled
// "PowerShell" with the same command, and a "PowerShell" block that
// continues lines with "`" needs a "Bash" twin. The two match when they
// are equal once each line's continuation character and surrounding
// whitespace are removed, and PowerShell's "\" path separators read as "/".
package shellpair
