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
	repo      *git.Repository
	branchMap map[string][]string
}

func NewService(path string) (*Service, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, err
	}
	s := &Service{repo: repo}
	s.BuildBranchMap()
	return s, nil
}

func (s *Service) BuildBranchMap() {
	s.branchMap = make(map[string][]string)
	branches, _ := s.GetBranches()
	for _, b := range branches {
		s.branchMap[b.Hash] = append(s.branchMap[b.Hash], b.Name)
	}
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

func (s *Service) GetFileContent(commitHash, filePath string) (string, error) {
	hash := plumbing.NewHash(commitHash)
	commit, err := s.repo.CommitObject(hash)

	if err != nil {
		return "", err
	}

	file, err := commit.File(filePath)
	if err != nil {
		return "", err
	}

	content, err := file.Contents()
	if err != nil {
		return "", err
	}

	return content, nil
}

func (s *Service) GetFileDiff(commitHash, filePath string) ([]types.DiffLine, error) {
	hash := plumbing.NewHash(commitHash)
	commit, err := s.repo.CommitObject(hash)
	if err != nil {
		return nil, err
	}

	if len(commit.ParentHashes) == 0 {
		file, err := commit.File(filePath)
		if err != nil {
			return nil, err
		}
		content, err := file.Contents()
		if err != nil {
			return nil, err
		}
		var lines []types.DiffLine
		for _, line := range strings.Split(content, "\n") {
			lines = append(lines, types.DiffLine{Type: "add", Content: line})
		}
		return lines, nil
	}

	parent, err := commit.Parent(0)
	if err != nil {
		return nil, err
	}

	patch, err := parent.Patch(commit)
	if err != nil {
		return nil, err
	}

	for _, fp := range patch.FilePatches() {
		from, to := fp.Files()
		var patchPath string
		if to != nil {
			patchPath = to.Path()
		} else if from != nil {
			patchPath = from.Path()
		}

		if patchPath != filePath {
			continue
		}

		var lines []types.DiffLine
		for _, chunk := range fp.Chunks() {
			content := chunk.Content()
			chunkLines := strings.Split(content, "\n")

			for i, line := range chunkLines {
				if i == len(chunkLines)-1 && line == "" {
					continue
				}
				var lineType string
				switch chunk.Type() {
				case diff.Add:
					lineType = "add"
				case diff.Delete:
					lineType = "del"
				default:
					lineType = "equal"
				}
				lines = append(lines, types.DiffLine{Type: lineType, Content: line})
			}
		}
		return lines, nil
	}

	return nil, nil
}

func (s *Service) GetCommits(branch string, limit int) ([]types.GraphCommit, error) {
	var fromHash plumbing.Hash
	if branch != "" {
		ref, err := s.repo.Reference(plumbing.NewBranchReferenceName(branch), true)
		if err != nil {
			ref, err = s.repo.Reference(plumbing.NewRemoteReferenceName("origin", branch), true)
			if err != nil {
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
		head, err := s.repo.Head()
		if err != nil {
			return nil, err
		}
		fromHash = head.Hash()
	}

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

		var parents []string
		for _, p := range c.ParentHashes {
			parents = append(parents, p.String())
		}

		graphChars := "* "
		isMerge := len(parents) > 1
		if isMerge {
			graphChars = "*─┐"
		}

		branchLabels := s.branchMap[c.Hash.String()]
		message := strings.Split(strings.TrimSpace(c.Message), "\n")[0]

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

func (s *Service) GetCommitDetails(fullHash string) ([]types.ParentInfo, []types.FileChange, error) {
	hash := plumbing.NewHash(fullHash)
	c, err := s.repo.CommitObject(hash)
	if err != nil {
		return nil, nil, err
	}

	var parentInfos []types.ParentInfo
	if len(c.ParentHashes) > 1 {
		for _, ph := range c.ParentHashes {
			parentCommit, perr := s.repo.CommitObject(ph)
			if perr == nil {
				parentMsg := strings.Split(strings.TrimSpace(parentCommit.Message), "\n")[0]
				parentBranches := s.branchMap[ph.String()]
				branchName := ""
				if len(parentBranches) > 0 {
					branchName = parentBranches[0]
				}
				parentInfos = append(parentInfos, types.ParentInfo{
					Hash:    ph.String()[:7],
					Message: parentMsg,
					Branch:  branchName,
				})
			}
		}
	}

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

	return parentInfos, files, nil
}
