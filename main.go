package main

import (
	"context"
	"flag"

	client "github.com/jgensler8/openfaas-controller/pkg/client/clientset/versioned"
	"github.com/jgensler8/openfaas-controller/pkg/controller"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"github.com/jgensler8/openfaas-client-go"
)

func main() {
	kubeconfig := flag.String("kubeconfig", "", "Path to a kube config. Only required if out-of-cluster.")
	openfaasBasePath := flag.String("openfaas-basepath", "http://openfaas.default.svc.cluster.local", "Base Path to use when contacting the openfaas controller")
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

	api := swagger.NewDefaultApiWithBasePath(*openfaasBasePath)

	// start a controller on instances of our custom resource
	crcontroller := controller.FunctionController{
		//FunctionClient: crclient.RESTClient(),
		KubernetesFunctionInterface: crclient.Cr().RESTClient(),
		OpenFaaSFunctionAPIClient: api,
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
