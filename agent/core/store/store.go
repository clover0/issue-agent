package store

const LastSubmissionKey = "last_submission"

type Store struct {
	changedFiles []File
	submissions  map[string]Submission
}

func NewStore() Store {
	sub := make(map[string]Submission)
	return Store{
		changedFiles: []File{},
		submissions:  sub,
	}
}

func (s *Store) AddChangedFile(f File) {
	s.changedFiles = append(s.changedFiles, f)
}

func (s *Store) ChangedFiles() []File {
	return s.changedFiles
}

func (s *Store) AddSubmission(key string, sub Submission) {
	s.submissions[key] = sub
}

// TODO: not returning nil
func (s *Store) GetSubmission(key string) *Submission {
	if v, ok := s.submissions[key]; !ok {
		return nil
	} else {
		return &v
	}
}
