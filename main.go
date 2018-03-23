package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	opkit "github.com/rook/operator-kit"
	mysql "github.com/tonya11en/mysql-operator/pkg/apis/myproject/v1alpha1"
	mysqlclient "github.com/tonya11en/mysql-operator/pkg/client/clientset/versioned/typed/myproject/v1alpha1"
	"k8s.io/api/core/v1"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	fmt.Println("Getting kubernetes context")
	context, mySqlClientset, err := createContext()
	if err != nil {
		fmt.Printf("failed to create context. %+v\n", err)
		os.Exit(1)
	}

	// Create and wait for CRD resources.
	fmt.Println("Registering the mysql resource")
	resources := []opkit.CustomResource{mysql.MySqlResource}
	err = opkit.CreateCustomResources(*context, resources)
	if err != nil {
		fmt.Printf("failed to create custom resource. %+v\n", err)
		os.Exit(1)
	}

	// Create signals to stop watching the resources.
	signalChan := make(chan os.Signal, 1)
	stopChan := make(chan struct{})
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// Start watching the mysql resource.
	fmt.Println("Watching the mysql resource")
	controller := newMySqlController(context, mySqlClientset)
	controller.StartWatch(v1.NamespaceAll, stopChan)

	for {
		select {
		case <-signalChan:
			fmt.Println("shutdown signal received, exiting...")
			close(stopChan)
			return
		}
	}
}

func createContext() (*opkit.Context, mysqlclient.MyprojectV1alpha1Interface, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get k8s config. %+v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get k8s client. %+v", err)
	}

	apiExtClientset, err := apiextensionsclient.NewForConfig(config)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create k8s API extension clientset. %+v", err)
	}

	mySqlClientset, err := mysqlclient.NewForConfig(config)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create mysql clientset. %+v", err)
	}

	context := &opkit.Context{
		Clientset:             clientset,
		APIExtensionClientset: apiExtClientset,
		Interval:              500 * time.Millisecond,
		Timeout:               60 * time.Second,
	}
	return context, mySqlClientset, nil

}
