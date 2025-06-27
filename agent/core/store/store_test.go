package store_test

import (
	"testing"

	"github.com/clover0/issue-agent/core/store"
	"github.com/clover0/issue-agent/test/assert"
)

func TestAddChangedFile(t *testing.T) {
	s := store.NewStore()
	f := store.File{}

	s.AddChangedFile(f)

	assert.Equal(t, len(s.ChangedFiles()), 1)
}

func TestChangedFiles(t *testing.T) {
	s := store.NewStore()
	f := store.File{}

	s.AddChangedFile(f)
	files := s.ChangedFiles()

	assert.Equal(t, len(files), 1)
	assert.Equal(t, files[0], f)
}

func TestAddAndGetSubmittedWork(t *testing.T) {
	s := store.NewStore()
	key := "test_key"
	work := store.SubmittedWork{}

	result := s.GetSubmittedWork("non_existent")

	assert.Nil(t, result)

	s.AddSubmittedWork(key, work)
	got := s.GetSubmittedWork(key)

	assert.Equal(t, *got, work)
}
