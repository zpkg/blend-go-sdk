package stats

// Tagger is an interface for specifying and retrieving default stats tags
type Tagger interface {
	AddDefaultTag(string, string)
	DefaultTags() []string
}
