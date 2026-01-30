package gitcontrol

import (
	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
)

type Repository struct {
	repo *git.Repository
}

func OpenRepository(path string) (*Repository, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, err
	}
	return &Repository{repo: repo}, nil
}

// DeleteBranch deletes the branch with the given name from the repository. (local & remote)
func (r *Repository) DeleteBranch(bName string) error {
	err := r.repo.Storer.RemoveReference(plumbing.NewBranchReferenceName(bName))
	return err
}

func (r *Repository) GetCurrentBranch() (string, error) {
	head, err := r.repo.Head()
	if err != nil {
		return "", err
	}
	return head.Name().Short(), nil
}
