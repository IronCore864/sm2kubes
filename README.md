# Secrets Manager Sync to EKS Cluster

A small app written in Golang that:

1. runs as a piece of Lambda function,
2. which is triggered automatically when there is a change in your AWS Secrets Manager,
3. so that the secrets in your Secrets Manager are synchronized to your Kubernetes cluster as native K8s Secrets.

[![Go](https://github.com/IronCore864/sm2kubes/actions/workflows/master.yaml/badge.svg)](https://github.com/IronCore864/sm2kubes/actions/workflows/master.yaml)

## Background

No matter if you use Hashicorp Vault, or AWS Secrets Manager, if your pod in K8s consumes these secrets, there are basically only two ways to do so at the moment:

- Inject secrets using annotations into the pod as a file (see [here](https://www.hashicorp.com/blog/injecting-vault-secrets-into-kubernetes-pods-via-a-sidecar) and [here](https://aws.amazon.com/blogs/containers/aws-secrets-controller-poc/)), then update your app code to read secrets from that file instead of from env vars,
- or update your application to call Vault/Secrets Manager directly, instead of reading files or env vars.

### The Downsides of the above Two Methods

For the first injection/annotation/file method, first of all, you need to update your YAML manifest. Then, you need to update your app code so that instead of reading from ENV vars, it reads a file.

For the second method, you don't have to update any manifest, but you must refactor your app to interact with the secrets manager APIs.

If you just got started and only have a limited number of apps, or you are creating a new app from scratch, the above two methods aren't that bad (except it counters [12-factor app](https://12factor.net/) which loves to store config as ENV vars).

But in reality, in any sizable project, chances are, you already have like 50 microservices up and running, and they are all reading configs from ConfigMaps and Secrets. You don't want to refactor all of them just because you want to use a state-of-the-art secret manager.

### Possible Improvement

If you are using HashiCorp Vault, it doesn't have webhooks to notify you when a secret is changed. The best you can do is poll. You can create a cronjob polling the Vault periodically, and when there is a change, you apply it to your K8s.

In this way, you don't need to change your app, it still reads from ENV vars, and you can use your Vault as the single source of truth. I made two solutions for this before, one for older versions of OpenShift, one for K8s, see [here](https://github.com/IronCore864/vs2yaml) and [here](https://github.com/IronCore864/vs2kubes).

### Downside of this Improvement

This solution's downside is clear: because there is no "hook" mechanism, you must rely on polling.

The downside of polling is the interval. If you set the interval too high, you waste compute resources. If you set the interval too low, you must wait for a more extended time when you did some change. There is no way to set the interval "just right."

If only there is a way to know when a secret is changed, we can immediately read it from the secrets manager and then write it into K8s as a native secret.

## AWS Secrets Manager Coming to the Rescue

Truth be told, AWS Secrets Manager doesn't have webhooks to notify you either.

But AWS does provide you a whole set of services that can be integrated together to basically do whatever you want. This is where AWS shines. It doesn't just provide you with a bunch of services. It gives you the possibility to create an automated solution.

Here, we can use CloudTrail + CloudTrail Log Event + CloudWatch Rule + Lambda + Secrets Manager API + K8s API to achieve our goal.

*NOTE* I don't want to get into details of Vault or Secrets Manager, which is better, etc.; they basically do the same thing. For me, I don't care much about which tool to use; I care more about if the tool fits in my 100% automation setup.

## Architecture

CloudTrail: when you do an AWS API call (even operation in AWS web console, which essentially uses API to interact with those services), like it's recorded in CloudTrail.

->

CloudTrail writes log events.

->

Cloudwatch rules can listen on CloudTrail logs.

->

If the rule is triggered, CloudWatch can trigger a Lambda function.

->

In Lambda, we use Secrets Manager API to fetch the changed secret, then use K8s API to change the secret in the K8s cluster.


## This Repo

Is the Lambda part.

GitHub Actions are used to build, pack and upload the binary into an S3 bucket then update the Lambda function code from the S3 bucket.

Naming: **s**(ecrets)**m**(anager)**2kube**(rnetes)**s**(ecret).

## The Cloud Part

See this repo: https://github.com/IronCore864/terraform-sm2kubes.

Created with Terraform, which creates the CloudWatch rule, the Lambda function definition
