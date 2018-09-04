// Copyright 2018 Google Inc. All Rights Reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.

package cloudfunction

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/logging"
	"contrib.go.opencensus.io/exporter/stackdriver"
	"go.opencensus.io/trace"
	"google.golang.org/genproto/googleapis/api/monitoredres"
)

// EnableStackdriverTrace enables Stackdriver tracing for all requests.
func EnableStackdriverTrace() error {
	projectId := os.Getenv("GCP_PROJECT")
	if projectId == "" {
		return fmt.Errorf("GCP_PROJECT environment variable unset or missing")
	}

	stackdriverExporter, err := stackdriver.NewExporter(stackdriver.Options{ProjectID: projectId})
	if err != nil {
		return err
	}

	trace.RegisterExporter(stackdriverExporter)
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})

	return nil
}

// NewStackdriverLogger returns a new Stackdriver logger based on
// the functions execution environment.
func NewStackdriverLogger() (*logging.Logger, error) {
	projectId := os.Getenv("GCP_PROJECT")
	if projectId == "" {
		return nil, fmt.Errorf("GCP_PROJECT environment variable unset or missing")
	}

	functionName := os.Getenv("FUNCTION_NAME")
	if functionName == "" {
		return nil, fmt.Errorf("FUNCTION_NAME environment variable unset or missing")
	}

	region := os.Getenv("FUNCTION_REGION")
	if region == "" {
		return nil, fmt.Errorf("FUNCTION_REGION environment variable unset or missing")
	}

	client, err := logging.NewClient(context.Background(), projectId)
	if err != nil {
		return nil, err
	}

	monitoredResource := monitoredres.MonitoredResource{
		Type: "cloud_function",
		Labels: map[string]string{
			"function_name": functionName,
			"region":        region,
		},
	}
	commonResource := logging.CommonResource(&monitoredResource)

	return client.Logger(functionName, commonResource), nil
}
