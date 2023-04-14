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
		flagsmith.WithRequestTimeout(30*time.Second),
		flagsmith.WithDefaultHandler(badDefaultHandler),
	)
	// The client will now start its update goroutine in the background, but
	// let's say that right now there is something wrong with the DNS so all
	// requests will hang.

	for {
		flags, err := client.GetEnvironmentFlags() // This will now block every iteration until network access is restored.
		if err != nil {
			fmt.Printf("Failed getting flags: %s\n", err)
		} else {
			v, err := flags.GetFeatureValue("example")
			if err != nil {
				fmt.Printf("An error occurred: %s\n", err)
			} else {
				fmt.Printf("We read this value: %s", v)
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func badDefaultHandler(feature string) flagsmith.Flag {
	switch feature {
	case "example":
		return flagsmith.Flag{
			Enabled: true,
			Value:   "hello",
		}
	default:
		// Related issue: https://github.com/Flagsmith/flagsmith/issues/2025
		// Here we have no idea if the flag returned is good. The best thing
		// we probably could do is to set the name to something whack.
		return flagsmith.Flag{
			FeatureName: "This is not the feature you are looking for",
		}
	}
}
