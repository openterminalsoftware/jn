package utils

import (
	"jn/colors"
	"regexp"
	"strings"
)

func HighlightMarkdown(line string) string {
	if strings.HasPrefix(line, "#") {
		return colors.Blue + line + colors.Reset
	}

	boldAsteriskRe := regexp.MustCompile(`\*\*([^\*]+)\*\*`)
	line = boldAsteriskRe.ReplaceAllString(line, colors.Bold+"$1"+colors.Reset)
	boldUnderscoreRe := regexp.MustCompile(`__([^_]+)__`)
	line = boldUnderscoreRe.ReplaceAllString(line, colors.Bold+"$1"+colors.Reset)

	italicAsteriskRe := regexp.MustCompile(`\*([^\*]+)\*`)
	line = italicAsteriskRe.ReplaceAllString(line, colors.Italic+"$1"+colors.Reset)
	italicUnderscoreRe := regexp.MustCompile(`_([^_]+)_`)
	line = italicUnderscoreRe.ReplaceAllString(line, colors.Italic+"$1"+colors.Reset)

	codeRe := regexp.MustCompile("`([^`]+)`")
	line = codeRe.ReplaceAllString(line, colors.Cyan+"$1"+colors.Reset)

	linkRe := regexp.MustCompile(`\[(.*?)](.*?)`)
	line = linkRe.ReplaceAllString(line,
		colors.Magenta+"[$1]"+colors.Reset+
			"("+colors.Blue+"$2"+colors.Reset+")")

	if strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "* ") || strings.HasPrefix(line, "+ ") {
		return colors.Green + line + colors.Reset
	}
	if strings.HasPrefix(line, ">") {
		return colors.Yellow + line + colors.Reset
	}
	if strings.HasPrefix(line, "```") {
		return colors.Cyan + line + colors.Reset
	}
	return line
}
