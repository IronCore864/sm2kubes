package util

import (
	"testing"
)

func TestIsSupportedEvent(t *testing.T) {
	tests := []string{"CreateSecret", "PutSecretValue", "DeleteSecret", "NotSupported"}

	expected := map[string]bool{
		"CreateSecret":   true,
		"PutSecretValue": true,
		"DeleteSecret":   true,
		"NotSupported":   false,
	}

	for _, name := range tests {
		res := IsSupportedEvent(name)
		if res != expected[name] {
			t.Errorf("Error verifying supported events %s, got: %t, want: %t.", name, res, expected[name])
		}
	}
}

func TestGetClusterNamespaceAndSecret(t *testing.T) {
	tests := map[string]string{
		"normal": "cluster/default/secret",
		"error":  "a-secret-doesnt-match-expected-format",
	}

	expected := map[string]map[string]string{
		"normal": map[string]string{
			"cluster":   "cluster",
			"namespace": "default",
			"secret":    "secret",
		},
		"error": map[string]string{
			"cluster":   "",
			"namespace": "",
			"secret":    "",
		},
	}

	for name, secret := range tests {
		cluster, namespace, sec, _ := GetClusterNamespaceAndSecret(secret)
		if cluster != expected[name]["cluster"] {
			t.Errorf("Error getting cluster from secret name, got: %s, want: %s.", cluster, expected[name]["cluster"])
		}
		if namespace != expected[name]["namespace"] {
			t.Errorf("Error getting namespace from secret name, got: %s, want: %s.", namespace, expected[name]["namespace"])
		}
		if sec != expected[name]["secret"] {
			t.Errorf("Error getting secret from secret name, got: %s, want: %s.", cluster, expected[name]["secret"])
		}
	}
}
