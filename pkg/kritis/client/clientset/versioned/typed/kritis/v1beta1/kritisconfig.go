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

// KritisConfigsGetter has a method to return a KritisConfigInterface.
// A group's client should implement this interface.
type KritisConfigsGetter interface {
	KritisConfigs() KritisConfigInterface
}

// KritisConfigInterface has methods to work with KritisConfig resources.
type KritisConfigInterface interface {
	Create(ctx context.Context, kritisConfig *v1beta1.KritisConfig, opts v1.CreateOptions) (*v1beta1.KritisConfig, error)
	Update(ctx context.Context, kritisConfig *v1beta1.KritisConfig, opts v1.UpdateOptions) (*v1beta1.KritisConfig, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1beta1.KritisConfig, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1beta1.KritisConfigList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1beta1.KritisConfig, err error)
	KritisConfigExpansion
}

// kritisConfigs implements KritisConfigInterface
type kritisConfigs struct {
	client rest.Interface
}

// newKritisConfigs returns a KritisConfigs
func newKritisConfigs(c *KritisV1beta1Client) *kritisConfigs {
	return &kritisConfigs{
		client: c.RESTClient(),
	}
}

// Get takes name of the kritisConfig, and returns the corresponding kritisConfig object, and an error if there is any.
func (c *kritisConfigs) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1beta1.KritisConfig, err error) {
	result = &v1beta1.KritisConfig{}
	err = c.client.Get().
		Resource("kritisconfigs").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of KritisConfigs that match those selectors.
func (c *kritisConfigs) List(ctx context.Context, opts v1.ListOptions) (result *v1beta1.KritisConfigList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1beta1.KritisConfigList{}
	err = c.client.Get().
		Resource("kritisconfigs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested kritisConfigs.
func (c *kritisConfigs) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("kritisconfigs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a kritisConfig and creates it.  Returns the server's representation of the kritisConfig, and an error, if there is any.
func (c *kritisConfigs) Create(ctx context.Context, kritisConfig *v1beta1.KritisConfig, opts v1.CreateOptions) (result *v1beta1.KritisConfig, err error) {
	result = &v1beta1.KritisConfig{}
	err = c.client.Post().
		Resource("kritisconfigs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(kritisConfig).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a kritisConfig and updates it. Returns the server's representation of the kritisConfig, and an error, if there is any.
func (c *kritisConfigs) Update(ctx context.Context, kritisConfig *v1beta1.KritisConfig, opts v1.UpdateOptions) (result *v1beta1.KritisConfig, err error) {
	result = &v1beta1.KritisConfig{}
	err = c.client.Put().
		Resource("kritisconfigs").
		Name(kritisConfig.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(kritisConfig).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the kritisConfig and deletes it. Returns an error if one occurs.
func (c *kritisConfigs) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Resource("kritisconfigs").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *kritisConfigs) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Resource("kritisconfigs").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched kritisConfig.
func (c *kritisConfigs) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1beta1.KritisConfig, err error) {
	result = &v1beta1.KritisConfig{}
	err = c.client.Patch(pt).
		Resource("kritisconfigs").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
