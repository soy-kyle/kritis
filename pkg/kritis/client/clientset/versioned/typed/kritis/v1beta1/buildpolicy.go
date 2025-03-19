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

// Code generated by client-gen. DO NOT EDIT.

package v1beta1

import (
	"context"
	"time"

	v1beta1 "github.com/soy-kyle/kritis/pkg/kritis/apis/kritis/v1beta1"
	scheme "github.com/soy-kyle/kritis/pkg/kritis/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// BuildPoliciesGetter has a method to return a BuildPolicyInterface.
// A group's client should implement this interface.
type BuildPoliciesGetter interface {
	BuildPolicies(namespace string) BuildPolicyInterface
}

// BuildPolicyInterface has methods to work with BuildPolicy resources.
type BuildPolicyInterface interface {
	Create(ctx context.Context, buildPolicy *v1beta1.BuildPolicy, opts v1.CreateOptions) (*v1beta1.BuildPolicy, error)
	Update(ctx context.Context, buildPolicy *v1beta1.BuildPolicy, opts v1.UpdateOptions) (*v1beta1.BuildPolicy, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1beta1.BuildPolicy, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1beta1.BuildPolicyList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1beta1.BuildPolicy, err error)
	BuildPolicyExpansion
}

// buildPolicies implements BuildPolicyInterface
type buildPolicies struct {
	client rest.Interface
	ns     string
}

// newBuildPolicies returns a BuildPolicies
func newBuildPolicies(c *KritisV1beta1Client, namespace string) *buildPolicies {
	return &buildPolicies{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the buildPolicy, and returns the corresponding buildPolicy object, and an error if there is any.
func (c *buildPolicies) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1beta1.BuildPolicy, err error) {
	result = &v1beta1.BuildPolicy{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("buildpolicies").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of BuildPolicies that match those selectors.
func (c *buildPolicies) List(ctx context.Context, opts v1.ListOptions) (result *v1beta1.BuildPolicyList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1beta1.BuildPolicyList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("buildpolicies").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested buildPolicies.
func (c *buildPolicies) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("buildpolicies").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a buildPolicy and creates it.  Returns the server's representation of the buildPolicy, and an error, if there is any.
func (c *buildPolicies) Create(ctx context.Context, buildPolicy *v1beta1.BuildPolicy, opts v1.CreateOptions) (result *v1beta1.BuildPolicy, err error) {
	result = &v1beta1.BuildPolicy{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("buildpolicies").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(buildPolicy).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a buildPolicy and updates it. Returns the server's representation of the buildPolicy, and an error, if there is any.
func (c *buildPolicies) Update(ctx context.Context, buildPolicy *v1beta1.BuildPolicy, opts v1.UpdateOptions) (result *v1beta1.BuildPolicy, err error) {
	result = &v1beta1.BuildPolicy{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("buildpolicies").
		Name(buildPolicy.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(buildPolicy).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the buildPolicy and deletes it. Returns an error if one occurs.
func (c *buildPolicies) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("buildpolicies").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *buildPolicies) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("buildpolicies").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched buildPolicy.
func (c *buildPolicies) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1beta1.BuildPolicy, err error) {
	result = &v1beta1.BuildPolicy{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("buildpolicies").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
