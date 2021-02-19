package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ironcore864/sm2kubes/pkg/eks"
	"github.com/ironcore864/sm2kubes/pkg/events"
	"github.com/ironcore864/sm2kubes/pkg/k8s"
	"github.com/ironcore864/sm2kubes/pkg/sm"
	"github.com/ironcore864/sm2kubes/pkg/util"
)

func handler(ctx context.Context, event events.SecretsManagerRequest) (string, error) {
	region := os.Getenv("region")
	if len(region) == 0 {
		errMsg := "env var region is not set"
		return errMsg, errors.New(errMsg)
	}

	eventName := events.GetEventNameFromEventStruct(&event)
	if util.IsSupportedEvent(eventName) {
		fmt.Printf("Event name: %s\n", eventName)
	} else {
		return "unsupported event type", nil
	}

	originalSecretName := events.GetRequestParametersName(&event)
	clusterName, namespace, k8sSecretName, err := util.GetClusterNamespaceAndSecret(originalSecretName)
	if err != nil {
		return err.Error(), nil
	}

	clientset, err := eks.GetEKSClientSet(clusterName, region)
	if err != nil {
		return fmt.Sprintf("error creating clientset: %v", err), err
	}

	switch eventName {
	case "DeleteSecret":
		err := k8s.DeleteSecret(ctx, clientset, k8sSecretName, namespace)
		if err != nil {
			return fmt.Sprintf("error deleting secret %v", err), err
		}
	default:
		// CreateSecret or PutSecretValue
		secret, err := sm.GetSecret(originalSecretName, region)
		if err != nil {
			return "error getting secret from Secrets Manager", err
		}

		err = k8s.UpsertSecret(ctx, clientset, k8sSecretName, namespace, k8s.BuildSecret(secret, k8sSecretName))
		if err != nil {
			return fmt.Sprintf("error creating secret %v", err), err
		}
	}

	return "Success", nil
}

func main() {
	lambda.Start(handler)
}
