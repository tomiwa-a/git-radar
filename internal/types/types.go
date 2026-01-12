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
