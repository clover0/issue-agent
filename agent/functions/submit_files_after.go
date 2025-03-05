package functions

import (
	"github.com/clover0/issue-agent/core/store"
)

func SubmitFilesAfter(s *store.Store, storeKey string, storeValue SubmitFilesOutput) {
	s.AddSubmittedWork(storeKey, store.SubmittedWork{
		BaseBranch:        storeValue.PushedBranch,
		PullRequestNumber: storeValue.PullRequestNumber,
	})
}
