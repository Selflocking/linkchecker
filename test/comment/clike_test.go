package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/work4dev/linkchecker/comment"
)

func TestCLike_GetComments(t *testing.T) {
	var code = `
#include <stdio.h>

/**
 * @brief get the sum of two int
 *
 * @param a
 * @param b
 * @return int
 */
int add(int a, int b) {
    // cacl sum
    return a + b;
}

int main() {
    int a = 1;
    int b = 2;
    // output the sum of a and b
    printf("a + b = %d", add(a, b));
    return 0;
}
`

	var expt = []string{
		`/**
 * @brief get the sum of two int
 *
 * @param a
 * @param b
 * @return int
 */`,
		`// cacl sum`,
		`// output the sum of a and b`,
	}

	parser := comment.CLike{}
	res := parser.GetComments(code)
	assert.EqualValues(t, expt, res)
}
