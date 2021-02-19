package events

import (
	"testing"
)

func TestGetEventNameFromEventStruct(t *testing.T) {
	tests := map[string]*SecretsManagerRequest{
		"normal": {Detail: Detail{EventName: "CreateSecret"}},
		"empty":  {},
	}

	expected := map[string]string{
		"normal": "CreateSecret",
		"empty":  "",
	}

	for name, request := range tests {
		res := GetEventNameFromEventStruct(request)
		if res != expected[name] {
			t.Errorf("Error getting eventName, got: %s, want: %s.", res, expected[name])
		}
	}
}
func TestGetRequestParametersName(t *testing.T) {
	tests := map[string]*SecretsManagerRequest{
		"create": {Detail: Detail{RequestParameters: RequestParameters{Name: "name"}}},
		"update": {Detail: Detail{RequestParameters: RequestParameters{SecretID: "id"}}},
		"empty":  {},
	}

	expected := map[string]string{
		"create": "name",
		"update": "id",
		"empty":  "",
	}

	for name, request := range tests {
		res := GetRequestParametersName(request)
		if res != expected[name] {
			t.Errorf("Error getting eventName, got: %s, want: %s.", res, expected[name])
		}
	}
}
