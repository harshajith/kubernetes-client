/**
 * Copyright (C) 2015 Red Hat, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *         http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package builder

import (
	"math/rand"
	"reflect"
	"strings"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	buildapiv1 "github.com/openshift/api/build/v1"
	"github.com/openshift/origin/pkg/git"
)

func TestBuildInfo(t *testing.T) {
	b := &buildapiv1.Build{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "sample-app",
			Namespace: "default",
		},
		Spec: buildapiv1.BuildSpec{
			CommonSpec: buildapiv1.CommonSpec{
				Source: buildapiv1.BuildSource{
					Git: &buildapiv1.GitBuildSource{
						URI: "github.com/openshift/sample-app",
						Ref: "master",
					},
				},
				Strategy: buildapiv1.BuildStrategy{
					SourceStrategy: &buildapiv1.SourceBuildStrategy{
						Env: []corev1.EnvVar{
							{Name: "RAILS_ENV", Value: "production"},
						},
					},
				},
			},
		},
	}
	sourceInfo := &git.SourceInfo{}
	sourceInfo.CommitID = "1575a90c569a7cc0eea84fbd3304d9df37c9f5ee"
	got := buildInfo(b, sourceInfo)
	want := []KeyValue{
		{"OPENSHIFT_BUILD_NAME", "sample-app"},
		{"OPENSHIFT_BUILD_NAMESPACE", "default"},
		{"OPENSHIFT_BUILD_SOURCE", "github.com/openshift/sample-app"},
		{"OPENSHIFT_BUILD_REFERENCE", "master"},
		{"OPENSHIFT_BUILD_COMMIT", "1575a90c569a7cc0eea84fbd3304d9df37c9f5ee"},
		{"RAILS_ENV", "production"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("buildInfo(%+v) = %+v; want %+v", b, got, want)
	}

	b.Spec.Revision = &buildapiv1.SourceRevision{
		Git: &buildapiv1.GitSourceRevision{
			Commit: "1575a90c569a7cc0eea84fbd3304d9df37c9f5ee",
		},
	}
	got = buildInfo(b, nil)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("buildInfo(%+v) = %+v; want %+v", b, got, want)
	}

}

func TestRandomBuildTag(t *testing.T) {
	tests := []struct {
		namespace, name string
		want            string
	}{
		{"test", "build-1", "docker.io/test/build-1:f1f85ff5"},
		// For long build namespace + build name, the returned random build tag
		// would be longer than the limit of reference.NameTotalLengthMax (255
		// chars). We do not truncate the repository name because it could create an
		// invalid repository name (e.g., namespace=abc, name=d, repo=abc/d,
		// trucated=abc/ -> invalid), so we simply take a SHA1 hash of the
		// repository name (which is guaranteed to be a valid repository name) and
		// preserve the random tag.
		{
			"namespace" + strings.Repeat(".namespace", 20),
			"name" + strings.Repeat(".name", 20),
			"3250b1b251c90df3963d8faf6525732f56c44f8e:f1f85ff5",
		},
	}
	for _, tt := range tests {
		rand.Seed(0)
		got := randomBuildTag(tt.namespace, tt.name)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("randomBuildTag(%q, %q) = %q, want %q", tt.namespace, tt.name, got, tt.want)
		}
	}
}

func TestRandomBuildTagNoDupes(t *testing.T) {
	rand.Seed(0)
	previous := make(map[string]struct{})
	for i := 0; i < 100; i++ {
		tag := randomBuildTag("test", "build-1")
		_, exists := previous[tag]
		if exists {
			t.Errorf("randomBuildTag returned a recently seen tag: %q", tag)
		}
		previous[tag] = struct{}{}
	}
}

func TestContainerName(t *testing.T) {
	rand.Seed(0)
	got := containerName("test-strategy", "my-build", "ns", "hook")
	want := "openshift_test-strategy-build_my-build_ns_hook_f1f85ff5"
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}
