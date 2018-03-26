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
	"k8s.io/api/apps/v1beta2"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
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

// Create a pod spec. Note that this is specific to the example found here:
// https://kubernetes.io/docs/tasks/run-application/run-single-instance-stateful-application/
func (c *MySqlController) makePodSpec(objName string, ctrName string, ctrImage string, port int32, podGroup string, envVars map[string]string) *v1.PodTemplateSpec {
	var env []v1.EnvVar
	for k, v := range envVars {
		env = append(env, v1.EnvVar{Name: k, Value: v})
	}

	volumeName := "mysql-persistent-storage"
	podSpec := &v1.PodTemplateSpec{
		ObjectMeta: meta_v1.ObjectMeta{
			Name:   objName,
			Labels: map[string]string{"app": "mysql"},
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  ctrName,
					Image: ctrImage,
					Env:   env,
					Ports: []v1.ContainerPort{
						{
							ContainerPort: port,
						},
					},
					VolumeMounts: []v1.VolumeMount{
						{
							Name:      volumeName,
							MountPath: "/var/lib/" + objName,
						},
					},
				},
			},
			Volumes: []v1.Volume{
				{
					Name: volumeName,
					VolumeSource: v1.VolumeSource{
						PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
							ClaimName: "mysql-pv-claim",
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
			Name:   name,
			Labels: map[string]string{"app": "mysql"},
		},
		Spec: v1.ServiceSpec{
			Selector: map[string]string{"app": "mysql"},
			Ports: []v1.ServicePort{
				{
					Port: port,
				},
			},
			ClusterIP: v1.ClusterIPNone,
		},
	})

	if err != nil {
		fmt.Println("failed to create service:", err)
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
			Name: getPvcName(name),
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
		fmt.Println("failed to create pvc:", err)
	}

	return pvc, err
}

// Make a deployment. Note that this is specific to the example found here:
// https://kubernetes.io/docs/tasks/run-application/run-single-instance-stateful-application/
func (c *MySqlController) makeDeployment(name string, podSpec v1.PodTemplateSpec) (*v1beta2.Deployment, error) {
	fmt.Println("Making deployment")
	appsClient := c.context.Clientset.AppsV1beta2()
	deployment, err := appsClient.Deployments(v1.NamespaceDefault).Create(&v1beta2.Deployment{
		ObjectMeta: meta_v1.ObjectMeta{
			Name: name,
		},
		Spec: v1beta2.DeploymentSpec{
			Template: podSpec,
			Selector: &meta_v1.LabelSelector{
				MatchLabels: map[string]string{"app": "mysql"},
			},
			Strategy: v1beta2.DeploymentStrategy{
				Type: v1beta2.RecreateDeploymentStrategyType,
			},
		},
	})

	if err != nil {
		fmt.Println("failed to create deployment:", err)
	}

	return deployment, err
}

func getPvcName(objName string) string {
	return objName + "-pv-claim"
}

func (c *MySqlController) onAdd(obj interface{}) {
	fmt.Println("Handling MySql add")

	s := obj.(*mysql.MySql).DeepCopy()

	_, err := c.makeService(s.Name, 3306)
	if !errors.IsAlreadyExists(err) && err != nil {
		return
	}

	_, err = c.makePVC(s.Name)
	if !errors.IsAlreadyExists(err) && err != nil {
		return
	}

	podEnvVars := map[string]string{
		"MYSQL_ROOT_PASSWORD": s.Spec.RootPassword,
	}
	podSpec := c.makePodSpec(s.Name, "mysql-ctr", s.Spec.Image, 3306, "mysql-pod-group", podEnvVars)
	c.makeDeployment(s.Name, *podSpec)
}

func (c *MySqlController) onUpdate(oldObj, newObj interface{}) {
	fmt.Println("Handling MySql update")

	// This is currently a no-op because the MySQL resource currently has nothing
	// in its spec that can be modified without being disruptive.
}

// This is a single-instance MySQL operator, so we can get away with deleting
// all objects related to the app. We have to do it this way also because
// cascading deletes (to specify all related items) aren't supported.
func (c *MySqlController) onDelete(obj interface{}) {
	fmt.Println("Handling MySql delete")

	s := obj.(*mysql.MySql).DeepCopy()
	var delOpts meta_v1.DeleteOptions
	listOpts := meta_v1.ListOptions{LabelSelector: "app=mysql"}

	// Delete deployments.
	appsClient := c.context.Clientset.AppsV1beta2()
	err := appsClient.Deployments(v1.NamespaceDefault).DeleteCollection(&delOpts, listOpts)
	if err != nil {
		fmt.Println("failed to delete deployment:", err)
	}

	// Delete service.
	coreV1Client := c.context.Clientset.CoreV1()
	err = coreV1Client.Services(v1.NamespaceDefault).Delete(s.Name, &delOpts)
	if err != nil {
		fmt.Println("failed to delete service:", err)
	}

	// Delete replica sets.
	err = appsClient.ReplicaSets(v1.NamespaceDefault).DeleteCollection(&delOpts, listOpts)
	if err != nil {
		fmt.Println("failed to delete replication controller:", err)
	}

	// Delete PVC.
	err = coreV1Client.PersistentVolumeClaims(v1.NamespaceDefault).Delete(getPvcName(s.Name), &delOpts)
	if err != nil {
		fmt.Println("failed to delete pvc:", err)
	}

	// Delete pods.
	err = coreV1Client.Pods(v1.NamespaceDefault).DeleteCollection(&delOpts, listOpts)
	if err != nil {
		fmt.Println("failed to delete pod:", err)
	}

}
