package main

import (
	"context"
	"flag"

	client "github.com/jgensler8/openfaas-controller/pkg/client/clientset/versioned"
	"github.com/jgensler8/openfaas-controller/pkg/controller"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	kubeconfig := flag.String("kubeconfig", "", "Path to a kube config. Only required if out-of-cluster.")
	flag.Parse()

	// Create the client config. Use kubeconfig if given, otherwise assume in-cluster.
	config, err := buildConfig(*kubeconfig)
	if err != nil {
		panic(err)
	}

	crclient, err := client.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	// start a controller on instances of our custom resource
	crcontroller := controller.FunctionController{
		//FunctionClient: crclient.RESTClient(),
		FunctionInterface: crclient.RESTClient(),
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	//go controller.Run(ctx)
	crcontroller.Run(ctx)
}

func buildConfig(kubeconfig string) (*rest.Config, error) {
	if kubeconfig != "" {
		return clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	return rest.InClusterConfig()
}
