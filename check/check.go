package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/gorhill/cronexpr"
	"github.com/pivotal-cf-experimental/croncourse-resource/models"
)

func getVersions(request models.CheckRequest) (models.CheckResponse, error) {

	expr, err := cronexpr.Parse(request.Source.Expression)
	if err != nil {
		return nil, err
	}

	versions := []models.Version{}

	shouldFire := false

	loc, err := time.LoadLocation(request.Source.Location)
	if err != nil {
		return nil, err
	}

	previouslyFiredAt := request.Version.Time
	if previouslyFiredAt.IsZero() {
		previouslyFiredAt = time.Now().In(loc).Add(-1 * time.Hour)
	}

	shouldFireAt := expr.Next(previouslyFiredAt)
	if previouslyFiredAt.Before(shouldFireAt) && time.Now().After(shouldFireAt) {
		shouldFire = true
	}

	if shouldFire {
		versions = append(versions, models.Version{
			Time: time.Now().In(loc),
		})
	}

	return versions, nil
}

func main() {
	var request models.CheckRequest

	err := json.NewDecoder(os.Stdin).Decode(&request)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error decoding payload: "+err.Error())
		os.Exit(1)
	}

	versions, err := getVersions(request)
	if err != nil {
		fmt.Fprintln(os.Stderr, "invalid crontab expression: "+err.Error())
		os.Exit(1)
	}

	json.NewEncoder(os.Stdout).Encode(versions)
}
