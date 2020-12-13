package concourse

import "github.com/gofabian/flo/drone"

type Config struct {
	SelfUpdateJob bool
	BuildJob      bool
	Branches      []string
	DronePipeline *drone.Pipeline
}
