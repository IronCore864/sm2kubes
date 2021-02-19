package util

import (
	"errors"
	"strings"
)

func IsSupportedEvent(eventName string) bool {
	switch eventName {
	case
		"CreateSecret",
		"PutSecretValue",
		"DeleteSecret":
		return true
	}
	return false
}

func GetClusterNamespaceAndSecret(secretName string) (string, string, string, error) {
	secretNameSections := strings.Split(secretName, "/")
	if len(secretNameSections) != 3 {
		return "", "", "", errors.New("secret naming format doesn't follow format cluster/namespace/name, ignore")
	}
	return secretNameSections[0], secretNameSections[1], secretNameSections[2], nil
}
