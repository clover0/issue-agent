package store

const LastSubmissionKey = "last_submission"

// Store is a struct that holds the changed files and the submissions by agent
type Store struct {
	changedFiles   []File
	submittedWorks map[string]SubmittedWork
}

func NewStore() Store {
	sub := make(map[string]SubmittedWork)
	return Store{
		changedFiles:   []File{},
		submittedWorks: sub,
	}
}

func (s *Store) AddChangedFile(f File) {
	s.changedFiles = append(s.changedFiles, f)
}

func (s *Store) ChangedFiles() []File {
	return s.changedFiles
}

func (s *Store) AddSubmittedWork(key string, sub SubmittedWork) {
	s.submittedWorks[key] = sub
}

// TODO: not returning nil
func (s *Store) GetSubmittedWork(key string) *SubmittedWork {
	if v, ok := s.submittedWorks[key]; !ok {
		return nil
	} else {
		return &v
	}
}
