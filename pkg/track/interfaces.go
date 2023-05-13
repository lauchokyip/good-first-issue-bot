package track

type TrackInputFiles interface {
	Update() error
	Add(filename string, urls []string) error
	Delete(filename string, urls []string) error
	PushUpdate() error
}
