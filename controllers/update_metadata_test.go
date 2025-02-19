/*
 * update_labels_test.go
 *
 * This source file is part of the FoundationDB open source project
 *
 * Copyright 2021 Apple Inc. and the FoundationDB project authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package controllers

import (
	fdbv1beta2 "github.com/FoundationDB/fdb-kubernetes-operator/api/v1beta2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = FDescribe("Update metadata", func() {
	type testCase struct {
		pod          *corev1.Pod
		metadata     metav1.ObjectMeta
		expected     bool
		expectedMeta metav1.ObjectMeta
	}

	DescribeTable("Test metadata correctness",
		func(tc testCase) {
			result := metadataCorrect(tc.metadata, &tc.pod.ObjectMeta)
			Expect(result).To(Equal(tc.expected))
			Expect(tc.pod.ObjectMeta.Labels).To(Equal(tc.expectedMeta.Labels))
			Expect(tc.pod.ObjectMeta.Annotations).To(Equal(tc.expectedMeta.Annotations))
		},
		Entry("Metadata matches with Pod metadata",
			testCase{
				pod: &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Annotations: map[string]string{
							fdbv1beta2.LastSpecKey: "1",
						},
					},
				},
				metadata: metav1.ObjectMeta{
					Annotations: map[string]string{
						fdbv1beta2.LastSpecKey: "1",
					},
				},
				expected: true,
				expectedMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						fdbv1beta2.LastSpecKey: "1",
					},
				},
			},
		),
		Entry("Metadata last spec is not matching",
			testCase{
				pod: &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Annotations: map[string]string{
							fdbv1beta2.LastSpecKey: "1",
						},
					},
				},
				metadata: metav1.ObjectMeta{
					Annotations: map[string]string{
						fdbv1beta2.LastSpecKey: "2",
					},
				},
				expected: true,
				expectedMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						fdbv1beta2.LastSpecKey: "1",
					},
				},
			},
		),
		Entry("Metadata Annotation is not matching",
			testCase{
				pod: &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Annotations: map[string]string{
							fdbv1beta2.LastSpecKey: "1",
							"special":              "43",
						},
					},
				},
				metadata: metav1.ObjectMeta{
					Annotations: map[string]string{
						fdbv1beta2.LastSpecKey: "1",
						"special":              "42",
					},
				},
				expected: false,
				expectedMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						fdbv1beta2.LastSpecKey: "1",
						"special":              "42",
					},
				},
			},
		),
		Entry("Missing annotation on metadata",
			testCase{
				pod: &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Annotations: map[string]string{
							fdbv1beta2.LastSpecKey: "1",
						},
					},
				},
				metadata: metav1.ObjectMeta{
					Annotations: map[string]string{
						fdbv1beta2.LastSpecKey: "1",
						"controller/X":         "wrong",
					},
				},
				expected: false,
				expectedMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						fdbv1beta2.LastSpecKey: "1",
						"controller/X":         "wrong",
					},
				},
			},
		),
		Entry("Ignore additional annotation",
			testCase{
				pod: &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Annotations: map[string]string{
							fdbv1beta2.LastSpecKey: "1",
							"controller/X":         "wrong",
						},
					},
				},
				metadata: metav1.ObjectMeta{
					Annotations: map[string]string{
						fdbv1beta2.LastSpecKey: "1",
					},
				},
				expected: true,
				expectedMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						fdbv1beta2.LastSpecKey: "1",
						"controller/X":         "wrong",
					},
				},
			},
		),
		Entry("Annotation has wrong value",
			testCase{
				pod: &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Annotations: map[string]string{
							fdbv1beta2.LastSpecKey: "1",
							"controller/X":         "true",
						},
					},
				},
				metadata: metav1.ObjectMeta{
					Annotations: map[string]string{
						fdbv1beta2.LastSpecKey: "1",
						"controller/X":         "wrong",
					},
				},
				expected: false,
				expectedMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						fdbv1beta2.LastSpecKey: "1",
						"controller/X":         "wrong",
					},
				},
			},
		),
		Entry("Ignore additional label",
			testCase{
				pod: &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Annotations: map[string]string{
							fdbv1beta2.LastSpecKey: "1",
						},
						Labels: map[string]string{
							"test": "test",
						},
					},
				},
				metadata: metav1.ObjectMeta{
					Annotations: map[string]string{
						fdbv1beta2.LastSpecKey: "1",
					},
				},
				expected: true,
				expectedMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						fdbv1beta2.LastSpecKey: "1",
					},
					Labels: map[string]string{
						"test": "test",
					},
				},
			},
		),
		Entry("Label has wrong value",
			testCase{
				pod: &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Annotations: map[string]string{
							fdbv1beta2.LastSpecKey: "1",
						},
						Labels: map[string]string{
							fdbv1beta2.FDBProcessClassLabel: "storage",
						},
					},
				},
				metadata: metav1.ObjectMeta{
					Annotations: map[string]string{
						fdbv1beta2.LastSpecKey: "1",
					},
					Labels: map[string]string{
						fdbv1beta2.FDBProcessClassLabel: "globalControllerLogger",
					},
				},
				expected: false,
				expectedMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						fdbv1beta2.LastSpecKey: "1",
					},

					Labels: map[string]string{
						fdbv1beta2.FDBProcessClassLabel: "globalControllerLogger",
					},
				},
			},
		),
		Entry("Metadata for a Pod running on a node",
			testCase{
				pod: &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Annotations: map[string]string{
							fdbv1beta2.LastSpecKey: "1",
						},
					},
					Spec: corev1.PodSpec{
						NodeName: "testing-node",
					},
				},
				metadata: metav1.ObjectMeta{
					Annotations: map[string]string{
						fdbv1beta2.LastSpecKey:    "1",
						fdbv1beta2.NodeAnnotation: "testing-node",
					},
				},
				expected: false,
				expectedMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						fdbv1beta2.LastSpecKey:    "1",
						fdbv1beta2.NodeAnnotation: "testing-node",
					},
				},
			},
		),
	)
})
