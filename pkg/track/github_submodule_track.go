package track

import (
	"bufio"
	"bytes"
	"log"
	"os"
	"os/exec"

	"github.com/lauchokyip/good-first-issue-bot/pkg/issues/util"
)

type GitSubmodule struct {
	basePath string
}

func NewTrackWithGitSubModule(basePath string) *GitSubmodule {
	return &GitSubmodule{
		basePath: basePath,
	}
}

func (gs *GitSubmodule) Update() error {
	log.Println("Running git submodule pull")
	pullCmd := exec.Command("git", "submodule", "foreach", "git", "pull", "origin", "master")
	err := pullCmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func (gs *GitSubmodule) PushUpdate() error {
	log.Println("Running git submodule push")
	commitCmd := exec.Command("git", "submodule", "foreach", "git", "commit", "-am", "update")
	err := commitCmd.Run()
	if err != nil {
		return err
	}
	pushCmd := exec.Command("git", "submodule", "foreach", "git", "push", "origin", "master")
	err = pushCmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func (gs *GitSubmodule) Add(filename string, urls []string) (err error) {
	gs.Update()

	f, err := os.OpenFile(gs.basePath+filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer func() {
		closeErr := f.Close()
		if err == nil {
			err = closeErr
		}
	}()

	for _, url := range urls {
		_, err := f.WriteString(url + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}

func (gs *GitSubmodule) Delete(filename string, urls []string) (err error) {
	gs.Update()

	// Read the file content
	fileContent, err := os.ReadFile(gs.basePath + filename)
	if err != nil {
		log.Printf("Failed to read file: %v", err)
	}

	// Create a temporary file
	tempFile, err := os.CreateTemp("", "")
	if err != nil {
		log.Printf("Failed to create temporary file: %v", err)
	}
	defer func() {
		closeErr := tempFile.Close()
		if err == nil {
			err = closeErr
		}
	}()

	// Create a writer for the temporary file
	writer := bufio.NewWriter(tempFile)

	// Create a scanner to read the file content
	scanner := bufio.NewScanner(bytes.NewReader(fileContent))

	lookup := util.SlicesToMap(urls)
	// Iterate over each line in the file
	for scanner.Scan() {
		line := scanner.Text()

		// Check if the line matches any string in the slice of strings to delete
		if lookup[line] {
			continue
		}

		// If the line does not match, write it to the temporary file
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			log.Printf("Failed to write to temporary file: %v\n", err)
		}

	}

	// Check for any errors during the scanning process
	if err := scanner.Err(); err != nil {
		log.Printf("Error encountered while scanning file: %v\n", err)
	}

	// Flush the writer to ensure all data is written to the temporary file
	err = writer.Flush()
	if err != nil {
		log.Printf("Failed to flush writer: %v\n", err)
	}

	// Replace the original file with the temporary file
	err = os.Rename(tempFile.Name(), gs.basePath+filename)
	if err != nil {
		log.Printf("Failed to replace file: %v\n", err)
	}

	return nil
}
