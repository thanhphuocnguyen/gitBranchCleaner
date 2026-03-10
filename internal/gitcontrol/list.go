package gitcontrol

import "github.com/go-git/go-git/v6/plumbing"

type Branch struct {
	Name      string
	FullName  string
	IsCurrent bool
	Hash      string
}

func (r *Repository) ListLocalBranches() ([]Branch, error) {
	branches := []Branch{}
	head, err := r.repo.Head()
	if err != nil {
		return branches, err
	}

	currentBranch := head.Name().Short()

	branchIter, err := r.repo.Branches()
	if err != nil {
		return branches, err
	}

	err = branchIter.ForEach(func(ref *plumbing.Reference) error {
		branchName := ref.Name().Short()
		branchHash := ref.Hash().String()
		isCurrent := branchName == currentBranch

		branch := Branch{
			Name:      branchName,
			IsCurrent: isCurrent,
			FullName:  ref.Name().String(),
			Hash:      branchHash,
		}
		branches = append(branches, branch)
		return nil
	})
	return branches, err
}
