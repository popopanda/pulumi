package main

import (
	"log"

	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metadatav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	appName := "empty"
	envName := "dev"
	teamName := "devops"
	namespaceName := appName

	metaData := &metadatav1.ObjectMetaArgs{
		Name: pulumi.String(appName),
		Labels: pulumi.StringMap{
			"app":  pulumi.String(appName),
			"env":  pulumi.String(envName),
			"team": pulumi.String(teamName),
		},
	}

	err := pulumiRun(namespaceName, appName, metaData)

	if err != nil {
		log.Fatal(err)
	}

}

func pulumiRun(namespaceName, appName string, metaData *metadatav1.ObjectMetaArgs) error {
	pulumi.Run(func(ctx *pulumi.Context) error {
		NS, err := createNS(ctx, namespaceName, metaData)

		if err != nil {
			return err
		}

		DP, err := createDeployment(ctx, appName, metaData)

		if err != nil {
			return err
		}

		ctx.Export("Namespace Name: ", NS.Metadata.Name())
		ctx.Export("deployment export: ", DP.Metadata.Name())

		return nil
	})
}

func createNS(ctx *pulumi.Context, namespaceName string, metaData *metadatav1.ObjectMetaArgs) (*corev1.Namespace, error) {
	myNamespace, err := corev1.NewNamespace(ctx, namespaceName, &corev1.NamespaceArgs{
		Metadata: metaData,
	})
	if err != nil {
		return nil, err
	}
	return myNamespace, nil
}

func createDeployment(ctx *pulumi.Context, appName string, metaData *metadatav1.ObjectMetaArgs) (*appsv1.Deployment, error) {
	deployment, err := appsv1.NewDeployment(ctx, "deployment", &appsv1.DeploymentArgs{
		Metadata: metaData,
		Spec: &appsv1.DeploymentSpecArgs{
			MinReadySeconds: pulumi.Int(30),
			Replicas:        pulumi.Int(3),
			Selector: &metadatav1.LabelSelectorArgs{
				MatchLabels: pulumi.StringMap{
					"app": pulumi.String(appName),
				},
			},
			Template: &corev1.PodTemplateSpecArgs{
				Metadata: metaData,
				Spec: &corev1.PodSpecArgs{
					Containers: &corev1.ContainerArray{
						&corev1.ContainerArgs{
							Image:           pulumi.String("nginx:latest"),
							Name:            pulumi.String("main"),
							ImagePullPolicy: pulumi.String("Always"),
							Ports: &corev1.ContainerPortArray{
								&corev1.ContainerPortArgs{
									ContainerPort: pulumi.Int(80),
								},
							},
						},
					},
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return deployment, nil
}
