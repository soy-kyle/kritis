/*
Copyright 2018 Google LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package kritisconfig

import (
	"context"
	"fmt"

	"github.com/soy-kyle/kritis/pkg/kritis/apis/kritis/v1beta1"
	clientset "github.com/soy-kyle/kritis/pkg/kritis/client/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

// KritisConfigs returns all KritisConfig in the cluster
func KritisConfigs() ([]v1beta1.KritisConfig, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, fmt.Errorf("error building config: %v", err)
	}

	client, err := clientset.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("error building clientset: %v", err)
	}
	list, err := client.KritisV1beta1().KritisConfigs().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("error listing all kritis configs: %v", err)
	}
	return list.Items, nil
}
