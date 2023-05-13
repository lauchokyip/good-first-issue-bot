package persist

import (
	"encoding/json"
	"os"
	"time"
)

type Event struct {
	LastCheckTime time.Time `json:"last_check_time"`
}

// persist records last write time
func Persist(persistPath string) (err error) {
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
		LastCheckTime: time.Now(),
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

// lastPersist returns isRecent and error
func LastPersist(persistPath string) (bool, error) {
	_, err := os.Stat(persistPath)
	if err != nil {
		return false, nil
	}

	data, err := os.ReadFile(persistPath)
	if err != nil {
		return false, err
	}
	event := Event{}
	json.Unmarshal(data, &event)

	truncateTime := event.LastCheckTime.Truncate(time.Hour)
	if time.Since(truncateTime) > time.Hour {
		return false, nil
	}

	return true, nil
}
