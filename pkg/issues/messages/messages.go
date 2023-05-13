package messages

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/go-github/github"
)

const (
	summaryTemplate = `Summary from yesterday:
Lists of subscribed issues which are inactive for more than %d day(s):
%s

List of issues which might interest you from yesterday:
%s
`

	foundNewIssueTemplate = `Found new issues:
%s
`
)

func NewSummary(inactiveIssues []*github.Issue, goodIssues []*github.Issue, inactiveDays int) (title, body string) {
	inactiveIssuesSlices := make([]string, len(inactiveIssues))
	for i, inactiveIssue := range inactiveIssues {
		inactiveIssuesSlices[i] = "* [" + inactiveIssue.GetTitle() + "]" + "(" + inactiveIssue.GetHTMLURL() + ")"
	}

	goodIssuesSlices := make([]string, len(goodIssues))
	for i, goodIssue := range goodIssues {
		goodIssuesSlices[i] = "* [" + goodIssue.GetTitle() + "]" + "(" + goodIssue.GetHTMLURL() + ")"
	}

	return time.Now().Format(time.DateOnly),
		fmt.Sprintf(
			summaryTemplate,
			inactiveDays,
			strings.Join(inactiveIssuesSlices, "\n"),
			strings.Join(goodIssuesSlices, "\n"),
		)
}

func NewIssues(goodIssues []*github.Issue) string {
	goodIssuesSlices := make([]string, len(goodIssues))
	for i, goodIssue := range goodIssues {
		goodIssuesSlices[i] = "* [" + goodIssue.GetTitle() + "]" + "(" + goodIssue.GetHTMLURL() + ")"
	}
	return fmt.Sprintf(
		foundNewIssueTemplate,
		strings.Join(goodIssuesSlices, "\n"),
	)
}
