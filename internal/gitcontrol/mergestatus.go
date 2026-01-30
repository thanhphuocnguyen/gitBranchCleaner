package gitcontrol

import (
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/object"
)

func (r *Repository) IsMerged(bName, mainBranch string) (bool, error) {
	brandRef, err := r.repo.Reference(plumbing.NewBranchReferenceName(bName), true)
	if err != nil {
		return false, err
	}
	mainRef, err := r.repo.Reference(plumbing.NewBranchReferenceName(mainBranch), true)
	if err != nil {
		return false, err
	}

	return r.isAncestor(brandRef.Hash(), mainRef.Hash())
}

func (r *Repository) isAncestor(ancHash, desHash plumbing.Hash) (bool, error) {
	commit, err := r.repo.CommitObject(desHash)
	if err != nil {
		return false, err
	}
	commitIter := object.NewCommitPreorderIter(commit, nil, nil)
	defer commitIter.Close()

	for {
		c, err := commitIter.Next()
		if err != nil {
			break
		}
		if c.Hash == ancHash {
			return true, nil
		}
	}
	return false, nil
}
