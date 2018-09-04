// Copyright 2018 Google Inc. All Rights Reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.

package cloudfunction

import (
	"fmt"
	"os"
)

// FunctionName returns the name of the function.
//
// Returns an error if the FUNCTION_NAME env var is missing or
// not set.
func FunctionName() (string, error) {
	functionName := os.Getenv("FUNCTION_NAME")
	if functionName == "" {
		return "", fmt.Errorf("FUNCTION_NAME environment variable unset or missing")
	}

	return functionName, nil
}

// ProjectID returns the project ID.
//
// Returns an error if the GCP_PROJECT env var is missing or
// not set.
func ProjectID() (string, error) {
	projectId := os.Getenv("GCP_PROJECT")
	if projectId == "" {
		return "", fmt.Errorf("GCP_PROJECT environment variable unset or missing")
	}

	return projectId, nil
}
