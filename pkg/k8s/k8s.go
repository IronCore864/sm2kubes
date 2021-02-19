package k8s

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/eks"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"sigs.k8s.io/aws-iam-authenticator/pkg/token"
)

// NewClientset returns a pointer to a kubernetes.Clientset object
// based on the output of describe eks cluster
func NewClientset(cluster *eks.Cluster) (*kubernetes.Clientset, error) {
	gen, err := token.NewGenerator(true, false)
	if err != nil {
		return nil, err
	}

	opts := &token.GetTokenOptions{
		ClusterID: aws.StringValue(cluster.Name),
	}

	tok, err := gen.GetWithOptions(opts)
	if err != nil {
		return nil, err
	}

	ca, err := base64.StdEncoding.DecodeString(aws.StringValue(cluster.CertificateAuthority.Data))
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(
		&rest.Config{
			Host:        aws.StringValue(cluster.Endpoint),
			BearerToken: tok.Token,
			TLSClientConfig: rest.TLSClientConfig{
				CAData: ca,
			},
		},
	)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

// BuildSecret creates a corev1.Secret object
func BuildSecret(secret, secretName string) *corev1.Secret {
	var sec map[string]interface{}
	err := json.Unmarshal([]byte(secret), &sec)
	if err != nil {
		panic(err)
	}

	data := make(map[string][]byte)
	for k, v := range sec {
		data[k] = []byte(fmt.Sprintf("%v", v))
	}

	result := corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   secretName,
			Labels: map[string]string{"managed-by-aws-secrets-manager": "true"},
		},
		Data: data,
		Type: "Opaque",
	}
	return &result
}

// UpsertSecret creates or updates a secret in K8s
func UpsertSecret(ctx context.Context, c kubernetes.Interface, name, namespace string, secret *corev1.Secret) error {
	_, err := c.CoreV1().Secrets(namespace).Get(ctx, name, metav1.GetOptions{})

	if err != nil {
		log.Println(err)
		if strings.Index(err.Error(), "not found") >= 0 {
			log.Printf("Secret %s doesn't exist, creating ...\n", name)
			_, err = c.CoreV1().Secrets(namespace).Create(ctx, secret, metav1.CreateOptions{})
			if err != nil {
				return err
			}
			log.Printf("Secret %s created!\n", name)
		}
	} else {
		log.Printf("Secret %s exist, updating ...\n", name)
		_, err = c.CoreV1().Secrets(namespace).Update(ctx, secret, metav1.UpdateOptions{})
		if err != nil {
			return err
		}
		log.Printf("Secret %s updated!\n", name)
	}

	return nil
}

// DeleteSecret deletes a secret from K8s
func DeleteSecret(ctx context.Context, c kubernetes.Interface, name, namespace string) error {
	_, err := c.CoreV1().Secrets(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if strings.Index(err.Error(), "not found") >= 0 {
			log.Printf("Secret %s doesn't exist\n", name)
			return nil
		}
		return err
	}

	err = c.CoreV1().Secrets(namespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		log.Println(err)
	}
	log.Printf("Secret %s deleted!\n", name)

	return nil
}
