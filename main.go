package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"

	"github.com/lauchokyip/good-first-issue-bot/internal/persist"
	"github.com/lauchokyip/good-first-issue-bot/pkg/issues/messages"
	"github.com/lauchokyip/good-first-issue-bot/pkg/issues/new_issue"
	"github.com/lauchokyip/good-first-issue-bot/pkg/issues/subscribe"
	"github.com/lauchokyip/good-first-issue-bot/pkg/issues/types"
	"github.com/lauchokyip/good-first-issue-bot/pkg/issues/util"
	"github.com/lauchokyip/good-first-issue-bot/pkg/track"
)

const (
	persistPath       = "/tmp/stupidbot"
	subscribeFilename = "subscribed_issues"
	repoFilename      = "repos"
	inactiveDays      = 7

	owner = "lauchokyip"
	repo  = "good-first-issue-repo"
)

var (
	sourceOfTruth = getWD() + "good-first-issue-repo/"
)

func getWD() string {
	dir, _ := os.Getwd()
	return dir + "/"
}

func main() {
	var tracker track.TrackInputFiles = track.NewTrackWithGitSubModule(sourceOfTruth)
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("GITHUB_TOKEN environment variable is not set")
	}
	ctx := context.TODO()
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tokenClient := oauth2.NewClient(ctx, tokenSource)
	client := github.NewClient(tokenClient)

	go func() {
		for {
			err := tracker.Update()
			if err != nil {
				log.Println(err)
			}
			time.Sleep(15 * time.Minute)
		}
	}()

	for {
		isRecent, isNewDay, issueNum, err := persist.LastPersist(persistPath)
		log.Printf("isRecent %t, isNewDay %t, issueNum %d\n", isRecent, isNewDay, issueNum)

		if err != nil {
			if _, ok := err.(*os.PathError); !ok {
				panic(err)
			}
		}

		if isNewDay {
			doInNewDay(ctx, client)
		} else {
			// work with old issues
			if !isRecent {
				doInTheSameDay(ctx, client, issueNum)
			}
		}

		nextDuration := time.Now().Truncate(time.Hour).Add(time.Hour)
		sleepingTime := time.Until(nextDuration)
		log.Println("Sleep for ", sleepingTime.Minutes(), " minutes")
		time.Sleep(sleepingTime)
	}
}

func doInNewDay(ctx context.Context, client *github.Client) {
	log.Println("creating summary")

	// Found new issues
	subscribedHandler := subscribe.NewFileSubscribed(sourceOfTruth+subscribeFilename, client)
	// refresh the cache
	_, err := subscribedHandler.GetAll()
	if err != nil {
		panic(err)
	}
	inActiveIssues, err := subscribedHandler.GetInactive(inactiveDays)
	if err != nil {
		panic(err)
	}
	goodIssuesHandler := new_issue.NewGoodIssues(
		sourceOfTruth+repoFilename,
		client,
		[]string{
			"good first issue",
			"help wanted",
		},
		nil,
	)
	oneDay := 24 * time.Hour
	goodIssues, err := goodIssuesHandler.GetAll(persist.TruncateToDay(time.Now().Add(-oneDay)))
	if err != nil {
		panic(err)
	}

	title, body := messages.NewSummary(inActiveIssues, goodIssues, inactiveDays)
	issue, _, err := client.Issues.Create(ctx, owner, repo, &github.IssueRequest{
		Title: &title,
		Body:  &body,
	})
	if err != nil {
		panic(err)
	}

	err = persist.Persist(persistPath, int(*issue.Number))
	if err != nil {
		panic(err)
	}
}

func doInTheSameDay(ctx context.Context, client *github.Client, issueNum int) {
	log.Println("persisting")
	if issueNum == -1 {
		log.Println("the persistence file did not exists, getting the max issue number")
		num, err := util.GetLargestIssueNumber(ctx, client, types.IssueQuery{Repo: repo, Owner: owner})
		if err != nil {
			panic(err)
		}
		issueNum = num
	}

	// check if new issue exists

	err := persist.Persist(persistPath, issueNum)
	if err != nil {
		panic(err)
	}
	goodIssuesHandler := new_issue.NewGoodIssues(
		sourceOfTruth+repoFilename,
		client,
		[]string{
			"good first issue",
			"help wanted",
		},
		nil,
	)

	// offset the time a bit
	lastHour := time.Now().Add(-5 * time.Minute).Truncate(time.Hour)
	goodIssues, err := goodIssuesHandler.GetAll(lastHour)
	if err != nil {
		panic(err)
	}
	log.Println(goodIssues)
	msg := messages.NewIssues(goodIssues)
	if len(goodIssues) != 0 {
		log.Println("found new issues!")
		_, _, err := client.Issues.CreateComment(
			ctx,
			owner,
			repo,
			issueNum,
			&github.IssueComment{Body: &msg},
		)
		if err != nil {
			panic(err)
		}
	}

	rateLimit, _, err := client.RateLimits(ctx)
	if err != nil {
		panic(err)
	}
	log.Println(rateLimit.String())
}
