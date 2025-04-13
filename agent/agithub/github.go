package agithub

import (
	"fmt"
	"os"

	"github.com/google/go-github/v70/github"
)

func NewGitHub() (*github.Client, error) {
	token, ok := os.LookupEnv("GITHUB_TOKEN")
	if !ok {
		return nil, fmt.Errorf("GITHUB_TOKEN is not set")
	}
	return github.NewClient(nil).WithAuthToken(token), nil
}
