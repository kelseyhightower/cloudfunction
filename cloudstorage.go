// Copyright 2018 Google Inc. All Rights Reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.

package cloudfunction

import (
	"context"
	"io"
	"io/ioutil"

	"cloud.google.com/go/storage"
)

// ObjectToTempFile fetches the given Google Cloud Storage object
// from the given GCS bucket and writes the contents to a temp file.
//
// Returns the full path to tempfile.
func ObjectToTempFile(client *storage.Client, bucketName, objectName string) (string, error) {
	var err error
	ctx := context.Background()

	if client == nil {
		client, err = storage.NewClient(ctx)
		if err != nil {
			return "", err
		}
	}

	o, err := client.Bucket(bucketName).Object(objectName).NewReader(ctx)
	if err != nil {
		return "", err
	}
	defer o.Close()

	t, err := ioutil.TempFile("", "")
	if err != nil {
		return "", err
	}
	defer t.Close()

	if _, err := io.Copy(t, o); err != nil {
		return "", err
	}

	return t.Name(), nil
}
