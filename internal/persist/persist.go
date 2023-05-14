package persist

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

type Event struct {
	LastCheckTime time.Time `json:"last_check_time"`
	Issues        int       `json:"issue_num"`
}

// persist records last write time
func Persist(persistPath string, issue int) (err error) {
	f, err := os.Create(persistPath)
	if err != nil {
		return err
	}
	defer func() {
		closeErr := f.Close()
		if err == nil {
			err = closeErr
		}
	}()

	newEvent := Event{
		LastCheckTime: time.Now().UTC(),
		Issues:        issue,
	}
	data, err := json.Marshal(&newEvent)
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	if err != nil {
		return err
	}

	return nil
}

// lastPersist retrieves information from persistent storage
func LastPersist(persistPath string) (isRecent, newDay bool, lastIssue int, err error) {
	_, err = os.Stat(persistPath)
	if err != nil {
		return false, false, -1, err
	}

	data, err := os.ReadFile(persistPath)
	if err != nil {
		return false, false, -1, err
	}
	event := Event{}
	json.Unmarshal(data, &event)
	log.Printf("Last Persistent time is %s\n", event.LastCheckTime)

	if isNewDay(event.LastCheckTime) {
		return false, true, -1, nil
	}

	truncateTime := event.LastCheckTime.Truncate(time.Hour)
	if time.Since(truncateTime) > time.Hour {
		return false, false, event.Issues, nil
	}

	return true, false, event.Issues, nil
}

func isNewDay(timestamp time.Time) bool {
	day := 24 * time.Hour
	// Get the current date
	currentDate := time.Now().UTC().Truncate(day)

	// Get the date of the timestamp
	timestampDate := timestamp.UTC().Truncate(day)

	// Compare the dates
	return currentDate.After(timestampDate)
}
