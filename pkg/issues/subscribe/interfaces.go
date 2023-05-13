package subscribe

import "github.com/google/go-github/github"

type SubscribedIssues interface {
	GetAll() ([]*github.Issue, error)
	GetInactive(days int) ([]*github.Issue, error)
	Delete(url string) error
	Update(url string) error
}
