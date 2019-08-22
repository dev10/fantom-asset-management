package types

import "strings"

// QueryResultTokens is a payload for a tokens query
type QueryResultTokens []Token

// String implements fmt.Stringer
func (r QueryResultTokens) String() string {
	array := make([]string, len(r))
	for _, t := range r {
		array = append(array, t.String())
	}
	return strings.Join(array, "\n")
}
