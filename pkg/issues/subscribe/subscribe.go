package subscribe

import (
	"bufio"
	"context"
	"log"
	"os"

	"github.com/google/go-github/github"
	"github.com/lauchokyip/good-first-issue-bot/pkg/issues/types"
	"github.com/lauchokyip/good-first-issue-bot/pkg/issues/util"
)

type FileSubscribed struct {
	cache    []*github.Issue
	filename string
	client   *github.Client
}

func NewFileSubscribed(filename string, client *github.Client) *FileSubscribed {
	return &FileSubscribed{
		filename: filename,
		client:   client,
	}
}

func (fs *FileSubscribed) GetAll() ([]*github.Issue, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	path := dir + "/" + fs.filename
	f, err := os.Open(path)
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

func fromStringsToNumberedGithubIssues(ctx context.Context, client *github.Client, queries []types.IssueQueryWithNumber) []*github.Issue {
	issues := make([]*github.Issue, len(queries))
	for _, q := range queries {
		issue, _, err := client.Issues.Get(ctx, q.Owner, q.Repo, q.Number)
		if err != nil {
			log.Println(err)
		}
		issues = append(issues, issue)
	}

	return issues
}
