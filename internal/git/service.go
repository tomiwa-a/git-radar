package git

import (
	"sort"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/tomiwa-a/git-radar/internal/types"
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
