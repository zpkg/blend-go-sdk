package profanity

import "strings"

// MergeRules merges variadic rule sets.
func MergeRules(ruleSets ...Rules) Rules {
	output := make(Rules)
	for _, rules := range ruleSets {
		for key, rule := range rules {
			output[key] = rule
		}
	}
	return output
}

// Rules is a map of string id to rule.
type Rules map[string]Rule

// String
func (r Rules) String() string {
	if len(r) == 0 {
		return "<empty>"
	}
	var output []string
	for _, rule := range r {
		output = append(output, rule.String())
	}
	return strings.Join(output, "\n")
}
