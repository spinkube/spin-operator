package runtimeconfig

import (
	"testing"

	spinv1 "github.com/spinkube/spin-operator/api/v1"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func secretKeySelector(name string) *corev1.SecretKeySelector {
	return &corev1.SecretKeySelector{
		LocalObjectReference: corev1.LocalObjectReference{
			Name: name,
		},
	}
}

func configMapSelector(name string) *corev1.ConfigMapKeySelector {
	return &corev1.ConfigMapKeySelector{
		LocalObjectReference: corev1.LocalObjectReference{
			Name: name,
		},
	}
}

func Test_ExtractRuntimeConfigDependencies(t *testing.T) {
	t.Parallel()

	basicSpec := spinv1.SpinAppSpec{}

	table := []struct {
		name               string
		inputAppSpec       func() spinv1.SpinAppSpec
		expectedSecrets    []types.NamespacedName
		expectedConfigMaps []types.NamespacedName
	}{
		{
			name: "only_static_values",
			inputAppSpec: func() spinv1.SpinAppSpec {
				spec := basicSpec
				spec.RuntimeConfig.KeyValueStores = []spinv1.KeyValueStoreConfig{
					{
						Name: "my-kv-store",
						Type: "magical",
						Options: []spinv1.RuntimeConfigOption{
							{
								Name:  "secret",
								Value: "i-put-secrets-in-plain-text",
							},
						},
					},
				}

				return spec
			},
		},
		{
			name: "unique_configmaps",
			inputAppSpec: func() spinv1.SpinAppSpec {
				spec := basicSpec
				spec.RuntimeConfig.SqliteDatabases = []spinv1.SqliteDatabaseConfig{
					{
						Name: "my-turso-db",
						Type: "libsql",
						Options: []spinv1.RuntimeConfigOption{
							{
								Name: "url",
								ValueFrom: &spinv1.RuntimeConfigVarSource{
									ConfigMapKeyRef: configMapSelector("my-cm-b"),
								},
							},
						},
					},
				}
				spec.RuntimeConfig.KeyValueStores = []spinv1.KeyValueStoreConfig{
					{
						Name: "my-kv-store",
						Type: "magical",
						Options: []spinv1.RuntimeConfigOption{
							{
								Name: "secret",
								ValueFrom: &spinv1.RuntimeConfigVarSource{
									ConfigMapKeyRef: configMapSelector("my-cm-a"),
								},
							},
						},
					},
				}

				return spec
			},
			expectedConfigMaps: []types.NamespacedName{
				{
					Name:      "my-cm-a",
					Namespace: "test-ns",
				},
				{
					Name:      "my-cm-b",
					Namespace: "test-ns",
				},
			},
		},
		{
			name: "duplicate_config_maps",
			inputAppSpec: func() spinv1.SpinAppSpec {
				spec := basicSpec
				spec.RuntimeConfig.SqliteDatabases = []spinv1.SqliteDatabaseConfig{
					{
						Name: "my-turso-db",
						Type: "libsql",
						Options: []spinv1.RuntimeConfigOption{
							{
								Name: "url",
								ValueFrom: &spinv1.RuntimeConfigVarSource{
									ConfigMapKeyRef: configMapSelector("my-cm-a"),
								},
							},
						},
					},
				}
				spec.RuntimeConfig.KeyValueStores = []spinv1.KeyValueStoreConfig{
					{
						Name: "my-kv-store",
						Type: "magical",
						Options: []spinv1.RuntimeConfigOption{
							{
								Name: "secret",
								ValueFrom: &spinv1.RuntimeConfigVarSource{
									ConfigMapKeyRef: configMapSelector("my-cm-a"),
								},
							},
						},
					},
				}

				return spec
			},
			expectedConfigMaps: []types.NamespacedName{
				{
					Name:      "my-cm-a",
					Namespace: "test-ns",
				},
			},
		},
		{
			name: "unique_secrets",
			inputAppSpec: func() spinv1.SpinAppSpec {
				spec := basicSpec
				spec.RuntimeConfig.SqliteDatabases = []spinv1.SqliteDatabaseConfig{
					{
						Name: "my-turso-db",
						Type: "libsql",
						Options: []spinv1.RuntimeConfigOption{
							{
								Name: "url",
								ValueFrom: &spinv1.RuntimeConfigVarSource{
									SecretKeyRef: secretKeySelector("my-secret-b"),
								},
							},
						},
					},
				}
				spec.RuntimeConfig.KeyValueStores = []spinv1.KeyValueStoreConfig{
					{
						Name: "my-kv-store",
						Type: "magical",
						Options: []spinv1.RuntimeConfigOption{
							{
								Name: "secret",
								ValueFrom: &spinv1.RuntimeConfigVarSource{
									SecretKeyRef: secretKeySelector("my-secret-a"),
								},
							},
						},
					},
				}

				return spec
			},
			expectedSecrets: []types.NamespacedName{
				{
					Name:      "my-secret-a",
					Namespace: "test-ns",
				},
				{
					Name:      "my-secret-b",
					Namespace: "test-ns",
				},
			},
		},
		{
			name: "duplicate_secrets",
			inputAppSpec: func() spinv1.SpinAppSpec {
				spec := basicSpec
				spec.RuntimeConfig.SqliteDatabases = []spinv1.SqliteDatabaseConfig{
					{
						Name: "my-turso-db",
						Type: "libsql",
						Options: []spinv1.RuntimeConfigOption{
							{
								Name: "url",
								ValueFrom: &spinv1.RuntimeConfigVarSource{
									SecretKeyRef: secretKeySelector("my-secret-a"),
								},
							},
						},
					},
				}
				spec.RuntimeConfig.KeyValueStores = []spinv1.KeyValueStoreConfig{
					{
						Name: "my-kv-store",
						Type: "magical",
						Options: []spinv1.RuntimeConfigOption{
							{
								Name: "secret",
								ValueFrom: &spinv1.RuntimeConfigVarSource{
									SecretKeyRef: secretKeySelector("my-secret-a"),
								},
							},
						},
					},
				}

				return spec
			},
			expectedSecrets: []types.NamespacedName{
				{
					Name:      "my-secret-a",
					Namespace: "test-ns",
				},
			},
		},
	}

	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			app := &spinv1.SpinApp{
				ObjectMeta: metav1.ObjectMeta{
					Name:      test.name,
					Namespace: "test-ns",
				},
				Spec: test.inputAppSpec(),
			}

			deps := extractRuntimeConfigDependencies(app)

			secrets := mapKeys(deps.Secrets)
			require.ElementsMatch(t, test.expectedSecrets, secrets)

			configMaps := mapKeys(deps.ConfigMaps)
			require.ElementsMatch(t, test.expectedConfigMaps, configMaps)
		})
	}
}

func mapKeys[T comparable, V any, M ~map[T]V](input M) []T {
	result := make([]T, 0, len(input))
	for key := range input {
		result = append(result, key)
	}
	return result
}
