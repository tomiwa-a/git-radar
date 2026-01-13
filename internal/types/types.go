package types

type FileChange struct {
	Status    string
	Path      string
	Additions int
	Deletions int
}

type DiffLine struct {
	Type    string // "add", "del", "equal", "collapse"
	Content string
}

type ParentInfo struct {
	Hash    string
	Message string
	Branch  string
}

// GraphCommit represents a commit in the visual graph
type GraphCommit struct {
	Hash        string
	FullHash    string
	Message     string
	Author      string
	Date        string
	Branches    []string
	Parents     []string
	ParentInfos []ParentInfo
	GraphChars  string
	IsMerge     bool
	Lane        int
	Files       []FileChange
}

type Branch struct {
	Name     string
	FullName string
	Hash     string
	IsRemote bool
	IsHead   bool
}
