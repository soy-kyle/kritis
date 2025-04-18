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

// Code generated by lister-gen. DO NOT EDIT.

package v1beta1

import (
	v1beta1 "github.com/soy-kyle/kritis/pkg/kritis/apis/kritis/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// BuildPolicyLister helps list BuildPolicies.
type BuildPolicyLister interface {
	// List lists all BuildPolicies in the indexer.
	List(selector labels.Selector) (ret []*v1beta1.BuildPolicy, err error)
	// BuildPolicies returns an object that can list and get BuildPolicies.
	BuildPolicies(namespace string) BuildPolicyNamespaceLister
	BuildPolicyListerExpansion
}

// buildPolicyLister implements the BuildPolicyLister interface.
type buildPolicyLister struct {
	indexer cache.Indexer
}

// NewBuildPolicyLister returns a new BuildPolicyLister.
func NewBuildPolicyLister(indexer cache.Indexer) BuildPolicyLister {
	return &buildPolicyLister{indexer: indexer}
}

// List lists all BuildPolicies in the indexer.
func (s *buildPolicyLister) List(selector labels.Selector) (ret []*v1beta1.BuildPolicy, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1beta1.BuildPolicy))
	})
	return ret, err
}

// BuildPolicies returns an object that can list and get BuildPolicies.
func (s *buildPolicyLister) BuildPolicies(namespace string) BuildPolicyNamespaceLister {
	return buildPolicyNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// BuildPolicyNamespaceLister helps list and get BuildPolicies.
type BuildPolicyNamespaceLister interface {
	// List lists all BuildPolicies in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1beta1.BuildPolicy, err error)
	// Get retrieves the BuildPolicy from the indexer for a given namespace and name.
	Get(name string) (*v1beta1.BuildPolicy, error)
	BuildPolicyNamespaceListerExpansion
}

// buildPolicyNamespaceLister implements the BuildPolicyNamespaceLister
// interface.
type buildPolicyNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all BuildPolicies in the indexer for a given namespace.
func (s buildPolicyNamespaceLister) List(selector labels.Selector) (ret []*v1beta1.BuildPolicy, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1beta1.BuildPolicy))
	})
	return ret, err
}

// Get retrieves the BuildPolicy from the indexer for a given namespace and name.
func (s buildPolicyNamespaceLister) Get(name string) (*v1beta1.BuildPolicy, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1beta1.Resource("buildpolicy"), name)
	}
	return obj.(*v1beta1.BuildPolicy), nil
}
