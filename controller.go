/*
Copyright 2018 Tony Allen. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Some of the code below came from:
https://github.com/rook/rook
https://github.com/rook/operator-kit

which have the same license.
*/

package main

import (
	"fmt"

	opkit "github.com/rook/operator-kit"
	mysql "github.com/tonya11en/mysql-operator/pkg/apis/myproject/v1alpha1"
	mysqlclient "github.com/tonya11en/mysql-operator/pkg/client/clientset/versioned/typed/myproject/v1alpha1"
	"k8s.io/api/core/v1"
	extensions "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/resource"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
)

// MySqlController represents a controller object for mysql custom resources
type MySqlController struct {
	context        *opkit.Context
	mySqlClientset mysqlclient.MyprojectV1alpha1Interface
}

// Creates a controller watching for mysql custom resources.
func newMySqlController(context *opkit.Context, mySqlClientset mysqlclient.MyprojectV1alpha1Interface) *MySqlController {
	return &MySqlController{
		context:        context,
		mySqlClientset: mySqlClientset,
	}
}

// Watch watches for instances of MySql custom resources and acts on them
func (c *MySqlController) StartWatch(namespace string, stopCh chan struct{}) error {
	fmt.Println("Starting watch on the mysql resource")

	resourceHandlers := cache.ResourceEventHandlerFuncs{
		AddFunc:    c.onAdd,
		UpdateFunc: c.onUpdate,
		DeleteFunc: c.onDelete,
	}
	restClient := c.mySqlClientset.RESTClient()
	watcher := opkit.NewWatcher(mysql.MySqlResource, namespace, resourceHandlers, restClient)
	go watcher.Watch(&mysql.MySql{}, stopCh)
	return nil
}

// Create a pod spec.
func (c *MySqlController) makePodSpec(objName string, ctrName string, ctrImage string, port int32, podGroup string) *v1.PodTemplateSpec {
	podSpec := &v1.PodTemplateSpec{
		ObjectMeta: meta_v1.ObjectMeta{
			Name: objName,
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  ctrName,
					Image: ctrImage,
					Ports: []v1.ContainerPort{
						{
							ContainerPort: port,
						},
					},
				},
			},
		},
	}
	return podSpec
}

// Create a service.
func (c *MySqlController) makeService(name string, port int32) (*v1.Service, error) {
	fmt.Println("Making svc")
	coreV1Client := c.context.Clientset.CoreV1()
	svc, err := coreV1Client.Services(v1.NamespaceDefault).Create(&v1.Service{
		ObjectMeta: meta_v1.ObjectMeta{
			Name: name,
		},
		Spec: v1.ServiceSpec{
			Type:     v1.ServiceTypeNodePort,
			Selector: map[string]string{"app": "mysql"},
			Ports: []v1.ServicePort{
				{
					Port: port,
				},
			},
		},
	})

	if err != nil {
		fmt.Println("failed to create service.", err)
	}

	return svc, err
}

// Create a PVC. Note that this is specific to the example found here:
// https://kubernetes.io/docs/tasks/run-application/run-single-instance-stateful-application/
func (c *MySqlController) makePVC(name string) (*v1.PersistentVolumeClaim, error) {
	fmt.Println("Making pvc")
	coreV1Client := c.context.Clientset.CoreV1()
	pvc, err := coreV1Client.PersistentVolumeClaims(v1.NamespaceDefault).Create(&v1.PersistentVolumeClaim{
		ObjectMeta: meta_v1.ObjectMeta{
			Name: name,
		},
		Spec: v1.PersistentVolumeClaimSpec{
			AccessModes: []v1.PersistentVolumeAccessMode{"ReadWriteOnce"},
			Resources: v1.ResourceRequirements{
				Requests: v1.ResourceList{
					"storage": resource.MustParse("20G"),
				},
			},
		},
	})

	if err != nil {
		fmt.Println("failed to create pvc.", err)
	}

	return pvc, err
}

// Make a deployment. Note that this is specific to the example found here:
// https://kubernetes.io/docs/tasks/run-application/run-single-instance-stateful-application/
func (c *MySqlController) makeDeployment(name string, podSpec v1.PodTemplateSpec) *extensions.Deployment {
	fmt.Println("Making deployment")
	i := int32(1)
	return &extensions.Deployment{
		ObjectMeta: meta_v1.ObjectMeta{
			Name: name,
		},
		Spec: extensions.DeploymentSpec{
			Template: podSpec,
			Replicas: &i},
	}
}

func (c *MySqlController) onAdd(obj interface{}) {
	fmt.Println("Handling MySql add")
	s := obj.(*mysql.MySql).DeepCopy()

	podSpec := c.makePodSpec(s.Name, "mysql-ctr", s.Spec.Image, 3306, "mysql-pod-group")

	_, err := c.makeService(s.Name, 3306)
	if err != nil {
		return
	}

	_, err = c.makePVC("mysql-pv-claim")
	if err != nil {
		return
	}

	c.makeDeployment(s.Name, *podSpec)
}

func (c *MySqlController) onUpdate(oldObj, newObj interface{}) {
	fmt.Println("Handling MySql update")
}

func (c *MySqlController) onDelete(obj interface{}) {
	fmt.Println("Handling MySql delete")
}
