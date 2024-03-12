package runtimeconfig

import (
	"fmt"

	spinv1alpha1 "github.com/spinkube/spin-operator/api/v1alpha1"
	"github.com/spinkube/spin-operator/pkg/secret"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

type EnvVariablesProviderOptions struct {
	Prefix     string `toml:"prefix,omitempty"`
	DotEnvPath string `toml:"dotenv_path,omitempty"`
}

type VariablesProvider struct {
	Type string `toml:"type,omitempty"`
	EnvVariablesProviderOptions
}

type KeyValueStoreOptions map[string]secret.String

type SQLiteDatabaseOptions map[string]secret.String

type Spin struct {
	Variables []VariablesProvider `toml:"config_provider,omitempty"`

	KeyValueStores map[string]KeyValueStoreOptions `toml:"key_value_store,omitempty"`

	SQLiteDatabases map[string]SQLiteDatabaseOptions `toml:"sqlite_database,omitempty"`

	LLMCompute map[string]secret.String `toml:"llm_compute,omitempty"`
}

func (s *Spin) AddKeyValueStore(
	name, storeType, namespace string,
	secrets map[types.NamespacedName]*corev1.Secret,
	configMaps map[types.NamespacedName]*corev1.ConfigMap,
	opts []spinv1alpha1.RuntimeConfigOption) error {
	if s.KeyValueStores == nil {
		s.KeyValueStores = make(map[string]KeyValueStoreOptions)
	}

	if _, ok := s.KeyValueStores[name]; ok {
		return fmt.Errorf("duplicate definition for key value store with name: %s", name)
	}

	options, err := renderOptionsIntoMap(storeType, namespace, opts, secrets, configMaps)
	if err != nil {
		return err
	}

	s.KeyValueStores[name] = options
	return nil
}

func (s *Spin) AddSQLiteDatabase(
	name, storeType, namespace string,
	secrets map[types.NamespacedName]*corev1.Secret,
	configMaps map[types.NamespacedName]*corev1.ConfigMap,
	opts []spinv1alpha1.RuntimeConfigOption) error {
	if s.SQLiteDatabases == nil {
		s.SQLiteDatabases = make(map[string]SQLiteDatabaseOptions)
	}
	if _, ok := s.SQLiteDatabases[name]; ok {
		return fmt.Errorf("duplicate definition for sqlite database with name: %s", name)
	}

	options, err := renderOptionsIntoMap(storeType, namespace, opts, secrets, configMaps)
	if err != nil {
		return err
	}

	s.SQLiteDatabases[name] = SQLiteDatabaseOptions(options)
	return nil
}

func (s *Spin) AddLLMCompute(computeType, namespace string,
	secrets map[types.NamespacedName]*corev1.Secret,
	configMaps map[types.NamespacedName]*corev1.ConfigMap,
	opts []spinv1alpha1.RuntimeConfigOption) error {
	computeOpts, err := renderOptionsIntoMap(computeType, namespace, opts, secrets, configMaps)
	if err != nil {
		return err
	}

	s.LLMCompute = computeOpts
	return nil
}

func renderOptionsIntoMap(typeOpt, namespace string,
	opts []spinv1alpha1.RuntimeConfigOption,
	secrets map[types.NamespacedName]*corev1.Secret, configMaps map[types.NamespacedName]*corev1.ConfigMap) (map[string]secret.String, error) {
	options := map[string]secret.String{
		"type": secret.String(typeOpt),
	}

	for _, opt := range opts {
		var value string
		if opt.Value != "" {
			value = opt.Value
		} else if valueFrom := opt.ValueFrom; valueFrom != nil {
			if cmKeyRef := valueFrom.ConfigMapKeyRef; cmKeyRef != nil {
				cm, ok := configMaps[types.NamespacedName{Name: cmKeyRef.Name, Namespace: namespace}]
				if !ok {
					// This error shouldn't happen - we validate dependencies ahead of time, add this as a fallback error
					return nil, fmt.Errorf("unmet dependency while building config: configmap (%s/%s) not found", namespace, cmKeyRef.Name)
				}

				value = cm.Data[cmKeyRef.Key]
			} else if secKeyRef := valueFrom.SecretKeyRef; secKeyRef != nil {
				sec, ok := secrets[types.NamespacedName{Name: secKeyRef.Name, Namespace: namespace}]
				if !ok {
					// This error shouldn't happen - we validate dependencies ahead of time, add this as a fallback error
					return nil, fmt.Errorf("unmet dependency while building config: secret (%s/%s) not found", namespace, secKeyRef.Name)
				}

				value = string(sec.Data[secKeyRef.Key])
			}
		}

		options[opt.Name] = secret.String(value)
	}

	return options, nil
}
