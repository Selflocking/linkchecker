package comment

import "regexp"

type CLike struct {
}

func (CLike) GetComments(source string) []string {
	var patterns = []*regexp.Regexp{
		regexp.MustCompile(`(^|[^\\])/\*[\s\S]*?(?:\*/|$)`),
		regexp.MustCompile(`(^|[^\\:])//.*`),
	}

	return GetCommentsByPatternsWithLookbehind(source, patterns)
}
