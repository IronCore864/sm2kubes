package eks

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/eks"
	"k8s.io/client-go/kubernetes"

	"github.com/ironcore864/sm2kubes/pkg/k8s"
)

// GetEKSClientSet describes an EKS cluster in a region
// then create a kubernetes.Clientset object
func GetEKSClientSet(name, region string) (*kubernetes.Clientset, error) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	}))
	eksSvc := eks.New(sess)

	input := &eks.DescribeClusterInput{
		Name: aws.String(name),
	}

	result, err := eksSvc.DescribeCluster(input)
	if err != nil {
		log.Fatalf("Error calling DescribeCluster: %v", err)
	}

	clientset, err := k8s.NewClientset(result.Cluster)
	if err != nil {
		panic(err.Error())
	}

	return clientset, nil
}
