package num

import (
	"fmt"
	"strings"
)

// SliceToString returns a string of array values separated by the specified character
func SliceToString(a []float64, sep string) string {
	if len(a) == 0 {
		return ""
	}

	b := make([]string, len(a))
	for i, v := range a {
		b[i] = fmt.Sprintf("%f", v)
	}
	return strings.Join(b, sep)
}
