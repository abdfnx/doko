package shared

import "errors"

var (
	// ErrNoContainer no docker container error.
	ErrNoContainer = errors.New("no docker container")
	// ErrNoImage no image error.
	ErrNoImage = errors.New("no image")
	// ErrNoVolume no volume error.
	ErrNoVolume = errors.New("no volume")
	// ErrNoNetwork no network error.
	ErrNoNetwork = errors.New("no network")
	// ErrDockerConnect cannot connect to docker engine error.
	ErrDockerConnect = errors.New("unable to connect to `docker`")
	// ErrSmallTermWindowSize cannot run doko because of a small terminal window size
	ErrSmallTermWindowSize = errors.New("because of a small terminal window size cannot run `doko`")
)
