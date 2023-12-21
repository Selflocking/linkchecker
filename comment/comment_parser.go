package comment

import "regexp"

type Parser interface {
	GetComments(source string) []string
}

func ParserFactory(ext string) Parser {
	switch ext {
	case ".c":
		return C{}
	case ".go":
		return CLike{}
	default:
		return nil
	}
}

func GetCommentsByPatterns(source string, patterns []*regexp.Regexp) []string {
	var comments []string

	for _, pattern := range patterns {
		comments = append(comments, pattern.FindAllString(source, -1)...)
	}

	return comments
}

// GetCommentsByPatternsWithLookbehind
/*
 * @brief get comments by patterns with lookbehind
 * lookbehind: https://prismjs.com/extending.html#object-notation
 * Simply put, the first character of the match is removed
 */
func GetCommentsByPatternsWithLookbehind(source string, patterns []*regexp.Regexp) []string {
	var comments []string

	for _, pattern := range patterns {
		comments = append(comments, pattern.FindAllString(source, -1)...)
	}

	for i, comment := range comments {
		comments[i] = comment[1:]
	}

	return comments
}
