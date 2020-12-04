package cli

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/gofabian/flo/concourse"
	"github.com/gofabian/flo/drone"
	"gopkg.in/yaml.v3"
)

func ConvertPipeline(flagSet *flag.FlagSet, args []string) error {
	inputPtr := flagSet.String("i", "", "path to Drone pipeline YAML (default: stdin)")
	outputPtr := flagSet.String("o", "", "path to Concourse pipeline YAML (default: stdout)")
	flagSet.Parse(args)
	fmt.Println("input:", *inputPtr)
	fmt.Println("output:", *outputPtr)

	// open file
	var inputFile *os.File
	if len(*inputPtr) > 0 {
		var err error
		inputFile, err = os.Open(*inputPtr)
		if err != nil {
			return fmt.Errorf("cannot open '%s': %w", *inputPtr, err)
		}
		defer inputFile.Close()
	} else {
		inputFile = os.Stdin
		fileInfo, err := inputFile.Stat()
		if err != nil {
			return fmt.Errorf("cannot get file info of '%s': %w", *inputPtr, err)
		}
		if (fileInfo.Mode() & os.ModeCharDevice) != 0 {
			return errors.New("provide Drone pipeline via stdin or -i option")
		}
	}

	// decode Drone pipeline
	reader := bufio.NewReader(inputFile)
	dronePipeline := drone.Pipeline{}
	decoder := yaml.NewDecoder(reader)
	decoder.KnownFields(true)
	err := decoder.Decode(&dronePipeline)
	if err != nil {
		return fmt.Errorf("cannot decode drone pipeline: %w", err)
	}

	// validate Drone pipeline
	errs := drone.ValidatePipeline(&dronePipeline)
	if len(errs) > 0 {
		msg := "Validation errors: "
		for _, e := range errs {
			msg += ", " + e.Error()
		}
		return errors.New(msg)
	}

	// create Concourse pipeline
	concoursePipeline := concourse.CreatePipeline(&dronePipeline)

	encoder := yaml.NewEncoder(os.Stdout)
	err = encoder.Encode(concoursePipeline)
	if err != nil {
		panic(err)
	}

	return nil
}
