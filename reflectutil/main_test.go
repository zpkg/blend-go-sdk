package reflectutil

type testType struct {
	ID        int
	Name      string
	NotTagged string
	Tagged    string
	SubTypes  []subType
}

type subType struct {
	ID   int
	Name string
}
