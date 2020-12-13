package concourse

import "github.com/gofabian/flo/drone"

type Config struct {
	SelfUpdateJob bool
	BuildJob      bool
	GitURL        string
	Branches      []string
	DronePipeline *drone.Pipeline
}
