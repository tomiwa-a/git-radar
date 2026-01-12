package types

type FileChange struct {
	Status    string
	Path      string
	Additions int
	Deletions int
}

type DiffLine struct {
	Type    string // "add", "del", "equal"
	Content string
}

// GraphCommit represents a commit in the visual graph
type GraphCommit struct {
	Hash       string
	FullHash   string
	Message    string
	Author     string
	Date       string
	Branches   []string
	Parents    []string
	GraphChars string
	IsMerge    bool
	Lane       int
	Files      []FileChange
}

type Branch struct {
	Name     string
	FullName string
	Hash     string
	IsRemote bool
	IsHead   bool
}
