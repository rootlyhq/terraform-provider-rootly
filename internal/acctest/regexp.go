package acctest

import (
	"regexp"
	"strings"

	"github.com/samber/lo"
)

// ExpectLiteralErrors joins multiple literal text snippets in order, matching any
// amount of whitespace between words and across the boundaries between strings.
func ExpectLiteralErrors(texts ...string) *regexp.Regexp {
	patterns := lo.Map(texts, func(text string, _ int) string {
		words := strings.Fields(text) // Splits on all whitespace (newlines, tabs, spaces)
		escaped := lo.Map(words, func(w string, _ int) string {
			return regexp.QuoteMeta(w)
		})
		return strings.Join(escaped, `\s+`)
	})

	// Join all snippet patterns with \s+ to ignore whitespace between snippets as well
	combinedPattern := strings.Join(patterns, `\s+`)
	return regexp.MustCompile(combinedPattern)
}
