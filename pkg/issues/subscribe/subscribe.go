package subscribe

import (
	"bufio"
	"context"
	"log"
	"os"
	"time"

	"github.com/google/go-github/github"
	"github.com/lauchokyip/good-first-issue-bot/pkg/issues/types"
	"github.com/lauchokyip/good-first-issue-bot/pkg/issues/util"
)

type FileSubscribed struct {
	cache    []*github.Issue
	filepath string
	client   *github.Client
}

func NewFileSubscribed(filepath string, client *github.Client) *FileSubscribed {
	return &FileSubscribed{
		filepath: filepath,
		client:   client,
	}
}
func fromStringsToNumberedGithubIssues(ctx context.Context, client *github.Client, queries []types.IssueQueryWithNumber) []*github.Issue {
	issues := []*github.Issue{}
	for _, q := range queries {
		issue, _, err := client.Issues.Get(ctx, q.Owner, q.Repo, q.Number)
		if err != nil {
			log.Println(err)
			continue
		} else {
			// TODO if closed, delete it from entry
			if issue.ClosedAt == nil {
				issues = append(issues, issue)
			}
		}
	}

	return issues
}

func (fs *FileSubscribed) GetAll() ([]*github.Issue, error) {
	f, err := os.Open(fs.filepath)
	if err != nil {
		return nil, err
	}

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(f)

	// Iterate over each line and store it in a slice of strings
	lines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	// Check for any errors that occurred during scanning
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	err = f.Close()
	if err != nil {
		return nil, err
	}

	queries, err := util.ConvertToAPIEndpoint(lines)
	if err != nil {
		return nil, err
	}

	issues := fromStringsToNumberedGithubIssues(
		context.TODO(),
		fs.client,
		queries,
	)

	// update cache
	fs.cache = issues

	return issues, nil
}

func (fs *FileSubscribed) GetInactive(days int) ([]*github.Issue, error) {
	if fs.cache == nil {
		_, err := fs.GetAll()
		if err != nil {
			return nil, err
		}
	}
	issues := fs.cache

	ret := []*github.Issue{}
	for _, issue := range issues {
		if time.Since(*issue.UpdatedAt) > (time.Duration(days) * 24 * time.Hour) {
			ret = append(ret, issue)
		}
	}

	return ret, nil
}
