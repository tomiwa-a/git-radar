package git

import (
	"sort"
	"strings"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/format/diff"
	"github.com/go-git/go-git/v6/plumbing/object"
	"github.com/tomiwa-a/git-radar/internal/types"
	"github.com/tomiwa-a/git-radar/utils"
)

type Service struct {
	repo *git.Repository
}

func NewService(path string) (*Service, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, err
	}
	return &Service{repo: repo}, nil
}

func (s *Service) GetBranches() ([]types.Branch, error) {
	var branches []types.Branch

	head, err := s.repo.Head()
	if err != nil {
		return nil, err
	}
	headBranch := head.Name().String()

	// Get local branches
	branchIter, err := s.repo.Branches()
	if err != nil {
		return nil, err
	}
	branchIter.ForEach(func(ref *plumbing.Reference) error {
		branch := types.Branch{
			Name:     ref.Name().Short(),
			FullName: ref.Name().String(),
			Hash:     ref.Hash().String(),
			IsRemote: false,
			IsHead:   ref.Name().String() == headBranch,
		}
		branches = append(branches, branch)
		return nil
	})

	// Get remote branches
	remoteIter, err := s.repo.References()
	if err != nil {
		return nil, err
	}
	remoteIter.ForEach(func(ref *plumbing.Reference) error {
		if ref.Name().IsRemote() {
			branch := types.Branch{
				Name:     ref.Name().Short(),
				FullName: ref.Name().String(),
				Hash:     ref.Hash().String(),
				IsRemote: true,
				IsHead:   false,
			}
			branches = append(branches, branch)
		}
		return nil
	})

	// Sort: local first, then remote, alphabetically within each group
	sort.Slice(branches, func(i, j int) bool {
		if branches[i].IsRemote != branches[j].IsRemote {
			return !branches[i].IsRemote
		}
		return branches[i].Name < branches[j].Name
	})

	return branches, nil
}

func (s *Service) GetCurrentBranch() (string, error) {
	head, err := s.repo.Head()
	if err != nil {
		return "", err
	}
	return head.Name().Short(), nil
}

func (s *Service) GetCommits(branch string, limit int) ([]types.GraphCommit, error) {
	// Build hash → branch names map
	branchMap := make(map[string][]string)
	branches, _ := s.GetBranches()
	for _, b := range branches {
		branchMap[b.Hash] = append(branchMap[b.Hash], b.Name)
	}

	// Resolve branch to reference
	var fromHash plumbing.Hash
	if branch != "" {
		ref, err := s.repo.Reference(plumbing.NewBranchReferenceName(branch), true)
		if err != nil {
			// Try remote branch
			ref, err = s.repo.Reference(plumbing.NewRemoteReferenceName("origin", branch), true)
			if err != nil {
				// Fall back to HEAD
				head, headErr := s.repo.Head()
				if headErr != nil {
					return nil, headErr
				}
				fromHash = head.Hash()
			} else {
				fromHash = ref.Hash()
			}
		} else {
			fromHash = ref.Hash()
		}
	} else {
		// Use HEAD if no branch specified
		head, err := s.repo.Head()
		if err != nil {
			return nil, err
		}
		fromHash = head.Hash()
	}

	// Get commit log
	commitIter, err := s.repo.Log(&git.LogOptions{
		From:  fromHash,
		Order: git.LogOrderCommitterTime,
	})
	if err != nil {
		return nil, err
	}

	var commits []types.GraphCommit
	count := 0

	err = commitIter.ForEach(func(c *object.Commit) error {
		if count >= limit {
			return nil
		}

		// Get parent hashes
		var parents []string
		for _, p := range c.ParentHashes {
			parents = append(parents, p.String())
		}

		// Determine graph character
		graphChars := "* "
		isMerge := len(parents) > 1
		if isMerge {
			graphChars = "*─┐"
		}

		// Get branch labels for this commit
		branchLabels := branchMap[c.Hash.String()]

		// Format message (first line only)
		message := strings.Split(strings.TrimSpace(c.Message), "\n")[0]

		var files []types.FileChange
		if len(c.ParentHashes) > 0 {
			parent, perr := c.Parent(0)
			if perr == nil {
				patch, perr := parent.Patch(c)
				if perr == nil {
					for _, fp := range patch.FilePatches() {
						from, to := fp.Files()
						name := ""
						if to != nil {
							name = to.Path()
						} else if from != nil {
							name = from.Path()
						}
						adds, dels := 0, 0
						for _, chunk := range fp.Chunks() {
							content := chunk.Content()
							lines := strings.Count(content, "\n")
							if len(content) > 0 && content[len(content)-1] != '\n' {
								lines++
							}
							switch chunk.Type() {
							case diff.Add:
								adds += lines
							case diff.Delete:
								dels += lines
							}
						}
						files = append(files, types.FileChange{Status: "M", Path: name, Additions: adds, Deletions: dels})
					}
				}
			}
		} else {
			tree, terr := c.Tree()
			if terr == nil {
				tree.Files().ForEach(func(f *object.File) error {
					additions := 0
					if contents, cerr := f.Contents(); cerr == nil {
						additions = strings.Count(contents, "\n")
						if len(contents) > 0 && contents[len(contents)-1] != '\n' {
							additions++
						}
					}
					files = append(files, types.FileChange{Status: "A", Path: f.Name, Additions: additions, Deletions: 0})
					return nil
				})
			}
		}

		commit := types.GraphCommit{
			Hash:       c.Hash.String()[:7],
			FullHash:   c.Hash.String(),
			Message:    message,
			Author:     c.Author.Name,
			Date:       utils.FormatRelativeTime(c.Author.When),
			Parents:    parents,
			Branches:   branchLabels,
			GraphChars: graphChars,
			IsMerge:    isMerge,
			Files:      files,
		}

		commits = append(commits, commit)
		count++
		return nil
	})

	if err != nil {
		return nil, err
	}

	return commits, nil
}
