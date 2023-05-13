package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"github.com/google/go-github/github"

	"github.com/lauchokyip/good-first-issue-bot/internal/persist"
	"github.com/lauchokyip/good-first-issue-bot/pkg/issues/new_issue"
	"github.com/lauchokyip/good-first-issue-bot/pkg/issues/subscribe"
)

const (
	persistPath       = "/tmp/stupidbot"
	subscribeFilename = "subscribed_issues"
	repoFilename      = "repos"
)

func main() {
	isRecent, err := persist.LastPersist(persistPath)
	if err != nil {
		panic(err)
	}
	if isRecent {
		return
	}

	if !isRecent {
		fmt.Println("more than an hour passed, need to persist again")
		err := persist.Persist(persistPath)
		if err != nil {
			panic(err)
		}

		client := github.NewClient(nil)

		fs := subscribe.NewFileSubscribed(subscribeFilename, client)
		subscribedIssues, err := fs.GetAll()
		if err != nil {
			log.Fatal(err)
		}

		goodFirstIssues := new_issue.NewGoodIssues(repoFilename, client, []string{"good first issue"}, nil)
		helpWantedIssues := new_issue.NewGoodIssues(repoFilename, client, []string{"help wanted"}, nil)

		for {
			fmt.Println("Email sent successfully!")
			time.Sleep(time.Hour)
		}
	}
}

// Helper function to encode a raw email message
func encodeMessage(to, subject, body string) string {
	message := "To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body

	return base64.URLEncoding.EncodeToString([]byte(message))
}
