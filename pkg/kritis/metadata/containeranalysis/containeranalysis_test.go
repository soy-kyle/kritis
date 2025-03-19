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

package containeranalysis

import (
	"testing"

	"github.com/soy-kyle/kritis/pkg/kritis/metadata"
	"github.com/soy-kyle/kritis/pkg/kritis/testutil"
)

func Test_isRegistryGCR(t *testing.T) {
	tests := []struct {
		name     string
		registry string
		expected bool
	}{
		{
			name:     "gcr image",
			registry: "gcr.io",
			expected: true,
		},
		{
			name:     "eu gcr image",
			registry: "eu.gcr.io",
			expected: true,
		},
		{
			name:     "invalid gcr image",
			registry: "foogcr.io",
			expected: false,
		},
		{
			name:     "non gcr image",
			registry: "index.docker.io",
			expected: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := isRegistryGCR(test.registry)
			testutil.DeepEqual(t, test.expected, actual)
		})
	}
}

func Test_isRegistryGAR(t *testing.T) {
	tests := []struct {
		name     string
		registry string
		expected bool
	}{
		{
			name:     "europe gar image",
			registry: "europe-docker.pkg.dev",
			expected: true,
		},
		{
			name:     "us-central1 gar image",
			registry: "us-central1-docker.pkg.dev",
			expected: true,
		},
		{
			name:     "invalid gar image",
			registry: "pkg.dev",
			expected: false,
		},
		{
			name:     "non gar image",
			registry: "index.docker.io",
			expected: false,
		},
		{
			name:     "invalid gar image",
			registry: "europe-python.pkg.dev",
			expected: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := isRegistryGAR(test.registry)
			testutil.DeepEqual(t, test.expected, actual)
		})
	}
}

func Test_getProjectFromContainerImage(t *testing.T) {
	tests := []struct {
		image   string
		project string
	}{
		{"gcr.io/project/1", "project"},
		{"gcr.io/project", "project"},
		{"gcr.io", ""},
	}
	for _, tc := range tests {
		t.Run(tc.image, func(t *testing.T) {
			testutil.DeepEqual(t, tc.project, getProjectFromContainerImage(tc.image))
		})
	}
}

func TestParseNoteReference(t *testing.T) {
	type ProjAndNote struct {
		projId string
		noteId string
	}
	tests := []struct {
		name   string
		input  string
		shdErr bool
		output ProjAndNote
	}{
		{"good", "projects/name/notes/noteName", false, ProjAndNote{"name", "noteName"}},
		{"bad1", "some", true, ProjAndNote{"", ""}},
		{"bad2", "v1aplha1/projects/name", true, ProjAndNote{"", ""}},
		{"bad3", "projects/name", true, ProjAndNote{"", ""}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actualProj, actualNote, err := metadata.ParseNoteReference(tc.input)
			testutil.CheckErrorAndDeepEqual(t, tc.shdErr, err, tc.output, ProjAndNote{actualProj, actualNote})
		})
	}
}
