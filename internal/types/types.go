package types

type Commit struct {
	Hash    string
	Message string
	Author  string
	Email   string
	Date    string
	Files   []FileChange
}

type FileChange struct {
	Status    string
	Path      string
	Additions int
	Deletions int
}

type GraphCommit struct {
	Hash       string
	Message    string
	Author     string
	Date       string
	Branches   []string
	Parents    []string
	GraphChars string
	IsMerge    bool
}
