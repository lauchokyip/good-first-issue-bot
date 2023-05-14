package track

import (
	"os"
	"testing"
)

const (
	basePath = "/home/lau/Projects/good-first-issue-bot/good-first-issue-repo/"
	testFile = "test"
)

func TestFileOperation(t *testing.T) {
	submodule := NewTrackWithGitSubModule(basePath)
	err := submodule.Add("test", []string{"haha", "haha", "hoho"})
	if err != nil {
		t.Fatal(err)
	}

	err = submodule.Delete("test", []string{"haha"})
	if err != nil {
		t.Fatal(err)
	}
}

func TestMain(m *testing.M) {
	m.Run()
	os.Remove(basePath + testFile)
}
