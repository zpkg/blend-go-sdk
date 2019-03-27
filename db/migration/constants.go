package migration

// Migration Stats
const (
	StatApplied = "applied"
	StatFailed  = "failed"
	StatSkipped = "skipped"
	StatTotal   = "total"
)

// Verbs and Nouns
const (
	VerbCreate = "create"
	VerbAlter  = "alter"
	verbRun    = "run"

	NounColumn     = "column"
	NounTable      = "table"
	NounIndex      = "index"
	NounConstraint = "constraint"
	NounRole       = "role"

	AdverbAlways    = "always"
	AdverbExists    = "exists"
	AdverbNotExists = "not exists"
)
