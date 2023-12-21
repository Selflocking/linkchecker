package comment

import "regexp"

type C struct {
}

func (C) GetComments(source string) []string {
	var patterns = []*regexp.Regexp{
		regexp.MustCompile(`//(?:[^\r\n\\]|\\(?:\r\n?|\n|(\?![\r\n])))*|/\*[\s\S]*?(?:\*/|$)`),
	}

	return GetCommentsByPatterns(source, patterns)
}
