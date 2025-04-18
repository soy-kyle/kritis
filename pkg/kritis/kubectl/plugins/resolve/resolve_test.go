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
package resolve

import (
	"fmt"
	"sort"
	"testing"

	"github.com/soy-kyle/kritis/pkg/kritis/testutil"
	yaml "gopkg.in/yaml.v2"
)

var testYaml1 = `apiVersion: v1
kind: Pod
metadata:
  name: test
spec:
  containers:
  - name: docker
    image: golang:1.12
    args: ["--arg1=<first>",
           "--arg2=<second>",
           "--arg3=<third>"]
`

var testYaml2 = `apiVersion: v1
kind: Pod
metadata:
  name: test
spec:
  containers:
  - name: digest
    image: gcr.io/google-appengine/debian9@sha256:547f82a1a5a194b22d1178f4c6aae3de006152757c0da267fd3a68b03e8b6d85
    env:
    key: ENV
    value: ENV_VALUE
    moreImages:
        image: gcr.io/distroless/base:debug
  - name: no-tag
    image: gcr.io/distroless/base
  - name: docker
    image: busybox
`

func Test_recursiveGetTaggedImages(t *testing.T) {
	tests := []struct {
		name     string
		yaml     string
		expected []string
	}{
		{
			name: "test one tagged image",
			yaml: testYaml1,
			expected: []string{
				"golang:1.12",
			},
		},
		{
			name: "test multiple tagged images",
			yaml: testYaml2,
			expected: []string{
				"busybox",
				"gcr.io/distroless/base",
				"gcr.io/distroless/base:debug",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m := yaml.MapSlice{}
			if err := yaml.Unmarshal([]byte(test.yaml), &m); err != nil {
				t.Fatalf("couldn't unmarshal yaml: %v", err)
			}
			actual := recursiveGetTaggedImages(m)
			sort.Strings(actual)
			testutil.DeepEqual(t, test.expected, actual)
		})
	}
}

type fakeResolver struct {
	tagMap map[string]string
}

func newFakeResolver() fakeResolver {
	f := fakeResolver{
		tagMap: map[string]string{},
	}
	return f
}

func (f *fakeResolver) resolve(image string) (string, error) {
	digest, ok := f.tagMap[image]
	if !ok {
		return "", fmt.Errorf("image %s not found", image)
	}
	return digest, nil
}

func setResolver(f func(string) (string, error)) func() {
	oldResolver := resolver
	resolver = f
	return func() {
		resolver = oldResolver
	}
}

func Test_resolveTagsToDigests(t *testing.T) {
	r := newFakeResolver()
	r.tagMap["gcr.io/google-appengine/debian9:2017-09-07-161610"] = "gcr.io/google-appengine/debian9@sha256:foo"
	r.tagMap["golang:1.12"] = "index.docker.io/library/golang@sha256:bar"

	defer setResolver(r.resolve)()

	tests := []struct {
		name     string
		images   []string
		expected map[string]string
	}{
		{
			name: "gcr image",
			images: []string{
				"gcr.io/google-appengine/debian9:2017-09-07-161610",
			},
			expected: map[string]string{
				"gcr.io/google-appengine/debian9:2017-09-07-161610": "gcr.io/google-appengine/debian9@sha256:foo",
			},
		},
		{
			name: "docker registry image",
			images: []string{
				"golang:1.12",
			},
			expected: map[string]string{
				"golang:1.12": "index.docker.io/library/golang@sha256:bar",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := resolveTagsToDigests(test.images)
			testutil.CheckErrorAndDeepEqual(t, false, err, test.expected, actual)
		})
	}
}

func Test_recursiveReplaceImage(t *testing.T) {
	tests := []struct {
		name         string
		yaml         interface{}
		replacements map[string]string
		expected     interface{}
	}{
		{
			name: "replace one image",
			yaml: yaml.MapSlice{
				yaml.MapItem{
					Key:   "image",
					Value: "image:tag",
				},
			},
			replacements: map[string]string{
				"image:tag": "image:digest",
			},
			expected: yaml.MapSlice{
				yaml.MapItem{
					Key:   "image",
					Value: "image:digest",
				},
			},
		},
		{
			name: "yaml without image field",
			yaml: yaml.MapSlice{
				yaml.MapItem{
					Key:   "key",
					Value: "image:tag",
				},
			},
			replacements: map[string]string{
				"image:tag": "image:digest",
			},
			expected: yaml.MapSlice{
				yaml.MapItem{
					Key:   "key",
					Value: "image:tag",
				},
			},
		},
		{
			name: "replace some images",
			yaml: formatMapSlice(t, []string{"image:tag", "something", "image:tag2"}),
			replacements: map[string]string{
				"image:tag":  "image:digest",
				"image:tag2": "image:digest2",
			},
			expected: formatMapSlice(t, []string{"image:digest", "something", "image:digest2"}),
		},
		{
			name:         "replace no images",
			yaml:         formatTestYaml1(t),
			replacements: map[string]string{},
			expected:     formatTestYaml1(t),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := recursiveReplaceImage(test.yaml, test.replacements)
			testutil.DeepEqual(t, test.expected, actual)
		})
	}
}

func Test_MultiYaml(t *testing.T) {
	r := newFakeResolver()
	r.tagMap["image:tag"] = "image:digest"
	r.tagMap["image:tag2"] = "image:digest2"

	defer setResolver(r.resolve)()
	multiYaml := `apiVersion: v1
kind: Pod
metadata:
  name: test
spec:
  containers:
  - name: first
    image: %s
---
apiVersion: v1
kind: Pod
metadata:
  name: test
spec:
  containers:
  - name: second
    image: %s
`
	expected := fmt.Sprintf(multiYaml, "image:digest", "image:digest2")
	actual, err := executeSubstitution(fmt.Sprintf(multiYaml, "image:tag", "image:tag2"))
	testutil.CheckErrorAndDeepEqual(t, false, err, expected, actual)
}

func formatTestYaml1(t *testing.T) yaml.MapSlice {
	t.Helper()
	m := yaml.MapSlice{}
	err := yaml.Unmarshal([]byte(testYaml1), &m)
	if err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	return m
}

func formatMapSlice(t *testing.T, args []string) yaml.MapSlice {
	t.Helper()
	testYaml := `apiVersion: v1
kind: Pod
metadata:
    name: test
    label: test
spec:
    containers:
    - name: tag
      image: "%s"
    env:
        key: ENV
        value: ENV_VALUE
    containers:
    - name: key1
      image: "image:digest"
      ports:
      - containerPort: 80
    - name: key2
      value: value
    - name: key3
      value: 6
    moreImages:
        image: "%s"
    nest:
        value:  0
        value1: 1
        value2: 2
        nest:
        nest:
            nest:
            - name: digest
            image: "%s"
`
	y := fmt.Sprintf(testYaml, args[0], args[1], args[2])

	m := yaml.MapSlice{}
	err := yaml.Unmarshal([]byte(y), &m)
	if err != nil {
		t.Fatalf("unmarshal failed. YAML:\n%s\nerror: %v", y, err)
	}
	return m
}
