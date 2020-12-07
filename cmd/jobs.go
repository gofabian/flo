package cmd

type JobType string

const (
	Refresh JobType = "refresh"
	Build   JobType = "build"
	All     JobType = "all"
)
