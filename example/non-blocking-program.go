package main

import (
	"context"
	"fmt"
	"os"
	"time"

	flagsmith "github.com/Flagsmith/flagsmith-go-client/v2"
)

func main() {
	ctx := context.Background()

	// Intialize the flagsmith client
	client := flagsmith.NewClient(os.Getenv("FLAGSMITH_ENVIRONMENT_KEY"),
		flagsmith.WithContext(ctx),
		flagsmith.WithLocalEvaluation(),
	)
	// The client will now start its update goroutine in the background.

	for {
		// The program continues and we want to resolve a feature immediately.
		flag, err := GetFeature(ctx, client, "example")
		if err != nil {
			fmt.Printf("An error occurred: %s\n")
		} else {
			fmt.Printf("We read this flag: %+v", flag)
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func GetFeature(ctx context.Context, client *flagsmith.Client, feature string) (flagsmith.Flag, error) {
	// We try to read the environment, and if the background task has completed
	// its update procedure we are happy with the latest info.
	env, err := client.ReadEnvironment()
	if err == nil {
		f, err := client.GetEnvironmentFlagsFromDocument(env)
		if err != nil {
			return flagsmith.Flag{}, err
		}
		return f.GetFlag(feature)
	}
	// Else we go directly to our default flag handler, which means we don't
	// run into any delays when starting our service.
	return customDefaultHandler(feature)
}

func customDefaultHandler(feature string) (flagsmith.Flag, error) {
	switch feature {
	case "example":
		return flagsmith.Flag{
			Enabled: true,
			Value:   "hello",
		}, nil
	default:
		// Related issue: https://github.com/Flagsmith/flagsmith/issues/2025
		// We probably want to be able to propagate an error back in case no
		// known flag exist.
		return flagsmith.Flag{}, fmt.Errorf("unknown feature flag %q", feature)
	}
}
