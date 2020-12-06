package concourse

import "strings"

type JobType string

const (
	Refresh JobType = "refresh"
	Build   JobType = "build"
	All     JobType = "all"
)

func createImageResource(image string) *ImageResource {
	return &ImageResource{
		Type:   "registry-image",
		Source: *createSourceFromImage(image),
	}
}

func createSourceFromImage(image string) *ImageSource {
	imageElements := strings.SplitN(image, ":", 2)
	repository := imageElements[0]
	var tag string
	if len(imageElements) > 1 {
		tag = imageElements[1]
	}
	return &ImageSource{
		Repository: repository,
		Tag:        tag,
	}
}
