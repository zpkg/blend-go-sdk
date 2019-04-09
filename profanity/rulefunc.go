package profanity

// RuleFunc is a function that evaluates a corpus.
type RuleFunc func(string, []byte) RuleResult
