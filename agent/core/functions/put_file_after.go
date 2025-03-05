package functions

import (
	"github.com/clover0/issue-agent/core/store"
)

// TODO: implement hook?
func StoreFileAfterPutFile(s *store.Store, file store.File) {
	s.AddChangedFile(file)
}
