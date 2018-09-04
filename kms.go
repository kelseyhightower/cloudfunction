// Copyright 2018 Google Inc. All Rights Reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.

package cloudfunction

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"

	cloudkms "cloud.google.com/go/kms/apiv1"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
)

// Decryptenv retrieves and decrypts the value of the environment variable
// named by the key.
//
// If client is nil, a new KMS client will be created.
// If keyId is empty, the value of the KMS_KEY_ID environment variable will
// be used.
//
// key is must be a base64 encoded encrypted value.
//
//    gcloud kms encrypt \
//      --location global \
//      --keyring my-keyring \
//      --key my-key \
//      --plaintext-file ${HOME}/.secret \
//      --ciphertext-file - | base64 -w 0
//
// It returns the decrypted value, which will be empty if the variable is
// not present or if an error occurs during decryption.
func Decryptenv(client *cloudkms.KeyManagementClient, keyId, key string) (string, error) {
	var err error
	ctx := context.Background()

	encryptedEnv := os.Getenv(key)
	if encryptedEnv == "" {
		return "", fmt.Errorf("%s environment variable unset or missing", key)
	}

	if keyId == "" {
		keyId = os.Getenv("KMS_KEY_ID")
	}

	if keyId == "" {
		return "", fmt.Errorf("KMS_KEY_ID environment variable unset or missing")
	}

	if client == nil {
		client, err = cloudkms.NewKeyManagementClient(ctx)
		if err != nil {
			return "", err
		}
	}

	encryptedData, err := base64.StdEncoding.DecodeString(encryptedEnv)
	if err != nil {
		return "", err
	}

	decryptRequest := &kmspb.DecryptRequest{
		Name:       keyId,
		Ciphertext: encryptedData,
	}

	resp, err := client.Decrypt(ctx, decryptRequest)
	if err != nil {
		return "", err
	}

	return string(resp.Plaintext), nil
}
