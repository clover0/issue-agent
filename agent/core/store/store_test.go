package store

import (
	"testing"

	"github.com/clover0/issue-agent/test/assert"
)

func TestNewStore(t *testing.T) {
	s := NewStore()

	assert.Equal(t, len(s.changedFiles), 0)
	assert.Equal(t, len(s.submittedWorks), 0)
}

func TestAddChangedFile(t *testing.T) {
	s := NewStore()
	f := File{}

	s.AddChangedFile(f)

	assert.Equal(t, len(s.changedFiles), 1)
}

func TestChangedFiles(t *testing.T) {
	s := NewStore()
	f := File{}

	s.AddChangedFile(f)
	files := s.ChangedFiles()

	assert.Equal(t, len(files), 1)
	assert.Equal(t, files[0], f)
}

func TestAddAndGetSubmittedWork(t *testing.T) {
	s := NewStore()
	key := "test_key"
	work := SubmittedWork{}

	result := s.GetSubmittedWork("non_existent")

	assert.Nil(t, result)

	s.AddSubmittedWork(key, work)
	got := s.GetSubmittedWork(key)

	assert.Equal(t, *got, work)
}
