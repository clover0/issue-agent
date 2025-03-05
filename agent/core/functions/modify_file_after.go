package functions

import (
	store2 "github.com/clover0/issue-agent/core/store"
)

func StoreFileAfterModifyFile(s *store2.Store, file store2.File) {
	StoreFileAfterPutFile(s, file)
}
