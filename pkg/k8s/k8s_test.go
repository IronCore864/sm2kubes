package k8s

import (
	"context"
	"fmt"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	testclient "k8s.io/client-go/kubernetes/fake"
)

func assertEqual(t *testing.T, a interface{}, b interface{}, message string) {
	if a == b {
		return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("%v != %v", a, b)
	}
	t.Fatal(message)
}

func TestBuildSecret(t *testing.T) {
	secretName := "test"
	secretJSON := "{ \"name\": \"tiexin\" }"
	res := BuildSecret(secretJSON, secretName)
	assertEqual(t, res.ObjectMeta.Name, secretName, "")
	assertEqual(t, res.Labels["managed-by-aws-secrets-manager"], "true", "")
	assertEqual(t, string(res.Data["name"]), "tiexin", "")
}

func TestUpsertSecretAndDelete(t *testing.T) {
	clientset := testclient.NewSimpleClientset()

	ctx := context.TODO()

	namespace := "default"
	secretName := "test"

	// get, should't exist
	secret, err := clientset.CoreV1().Secrets(namespace).Get(ctx, secretName, metav1.GetOptions{})
	if secret != nil {
		fmt.Println("Error: secret should not exist yet")
	}

	// create
	err = UpsertSecret(ctx, clientset, secretName, namespace, BuildSecret("{ \"name\": \"tiexin\" }", secretName))
	if err != nil {
		t.Fatal(err.Error())
	}

	// update
	err = UpsertSecret(ctx, clientset, secretName, namespace, BuildSecret("{ \"name\": \"tiexin\" }", secretName))
	if err != nil {
		t.Fatal(err.Error())
	}

	// delete
	err = DeleteSecret(context.TODO(), clientset, secretName, namespace)
	if err != nil {
		t.Fatal(err.Error())
	}

	// get, should't exist
	secret, err = clientset.CoreV1().Secrets(namespace).Get(ctx, secretName, metav1.GetOptions{})
	if secret != nil {
		fmt.Println("Error: secret should not exist yet")
	}

	// delete again on non-existing secret, shouldn't error
	err = DeleteSecret(context.TODO(), clientset, secretName, namespace)
	if err != nil {
		t.Fatal(err.Error())
	}
}
