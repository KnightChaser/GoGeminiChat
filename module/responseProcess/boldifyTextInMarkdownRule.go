package responseProcess

import (
	"regexp"
	"strings"

	"github.com/fatih/color"
)

func BoldifyTextInMarkdownRule(input string) string {
	// Define a regular expression to match double asterisks within words
	re := regexp.MustCompile(`\*\*(.+?)\*\*`)

	// Replace double asterisks with Markdown bold syntax
	modifiedInput := re.ReplaceAllStringFunc(input, func(match string) string {
		word := strings.Trim(match, "**")
		boldPrint := color.New(color.Bold)
		boldWord := boldPrint.Sprintf("%s", word)
		return boldWord
	})

	return modifiedInput
}
