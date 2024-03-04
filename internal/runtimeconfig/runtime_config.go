// package runtimeconfig provides implementations and helpers for managing Spin
// runtime config.
//
// It is mostly a terrifying pile of secrets.
//
// The core flow of rendering a Spin Runtime Config from the CRD is as follows:
//  1. We extract all of the config options from the spin app
//  2. For all of the config options we build a map of secrets and config maps
//     this has the advantage of de-duping whole-secrets into a single reference
//     when different keys may be re-used.
//  3. We fetch all of those secrets and config maps, returning an error if any
//     are not found. (We do not currently support "optional" secrets, as there
//     are no runtimeConfig options that would make sense to be optional).
//  4. We then iterate over the RuntimeConfig CRD again, with the augmented data,
//     to finally build a `runtimeconfig.Spin` that has populated config options.
//
// This mostly results in a lot of incredibly ugly code because we need to translate
// between very disparate schemas.
//
// Spin Runtime Config is modelled using Rust `enums` for its schema, with schemas
// varying based on the `type` option - this isn't something that can be cleanly
// modelled in a Kubernetes CRD, which is where our "type plus list of options"
// schema comes from.
// i.e:
//
//	keyValueStores:
//	  - name: "mystore"
//	    type: "redis"
//	    options:
//	      - name: url
//	        value: "redis://localghost:9000"
//	  - name: "myotherstore"
//	    type: "sqlite"
//	    options:
//	      - name: path
//	        value: "/mnt/store/redis.db"
//
// Or when sourcing a value from a secret:
//
//	keyValueStores:
//	  - name: "mystore"
//	    type: "redis"
//	    options:
//	      - name: url
//	        valueFrom:
//	         secretKeyRef:
//	           name: "my-secret"
//	           key: "redis-url"
//
// Will render into toml that looks something like:
//
//	[key_value_store.mystore]
//	type = "redis"
//	url = "redis://localghost:9000"
//
//	[key_value_store.myotherstore]
//	type = "sqlite"
//	path = "/mnt/store/redis.db"
//
// To maximize compatibility with different spin options + custom builds, we do
// very little validation of runtime config options in the operator.
package runtimeconfig

import (
	"context"
	"slices"

	spinv1 "github.com/spinkube/spin-operator/api/v1"
	"github.com/spinkube/spin-operator/internal/generics"
	"github.com/spinkube/spin-operator/internal/logging"
	"golang.org/x/sync/errgroup"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Builder interface {
	// Build takes a spin app and attempts to fetch any dependent secrets and build
	// a Spin-compatible representation of the configuration that can be rendered into
	// a new Secret.
	Build(ctx context.Context, app *spinv1.SpinApp) (*Spin, error)
}

func NewBuilder(client client.Client) *K8sBuilder {
	return &K8sBuilder{
		client: client,
	}
}

type K8sBuilder struct {
	client client.Client
}

func (k *K8sBuilder) Build(ctx context.Context, app *spinv1.SpinApp) (rc *Spin, err error) {
	logger := logging.FromContext(ctx).WithValues("component", "runtime_config_builder")
	defer func() {
		if err != nil {
			logger.Error(err, "failed to build runtime config")
			return
		}

		logger.Debug("built runtime config")
	}()

	deps := extractExternalDependencies(app)
	err = deps.fetch(ctx, k.client)
	if err != nil {
		return nil, err
	}

	runtimeConfig := app.Spec.RuntimeConfig
	if runtimeConfig.LoadFromSecret != "" {
		return nil, nil
	}

	rc = &Spin{}

	// We don't currently expose a VariablesProvider in the CRD as we expose bindings
	// to Kubernetes ConfigMaps and Secrets directly.
	// We should consider adding support for configuring alternative Variables
	// providers in the future.
	rc.Variables = []VariablesProvider{
		{
			Type: "env",
			EnvVariablesProviderOptions: EnvVariablesProviderOptions{
				// SPIN_VARIABLE_ is the default prefix, but by specifying it here we lightly
				// protect against changes to Spin Defaults.
				Prefix: "SPIN_VARIABLE_",
			},
		},
	}

	for _, kvStore := range runtimeConfig.KeyValueStores {
		err := rc.AddKeyValueStore(kvStore.Name, kvStore.Type, app.ObjectMeta.Namespace,
			deps.Secrets, deps.ConfigMaps, kvStore.Options)
		if err != nil {
			return nil, err
		}
	}

	for _, database := range runtimeConfig.SqliteDatabases {
		err := rc.AddSQLiteDatabase(database.Name, database.Type, app.ObjectMeta.Namespace,
			deps.Secrets, deps.ConfigMaps, database.Options)
		if err != nil {
			return nil, err
		}
	}
	if llm := runtimeConfig.LLMCompute; llm != nil {
		err := rc.AddLLMCompute(llm.Type, app.ObjectMeta.Namespace, deps.Secrets, deps.ConfigMaps, llm.Options)
		if err != nil {
			return nil, err
		}
	}

	return rc, nil
}

// externalDependencies is a horrible hack around extracting all of the required
// secrets and config maps for a given spin application's runtimeConfig and
// realizing them all with limited concurrency. This should _probably_ be
// refactored or replaced in the future to be a little less brittle and easier to manage.
type externalDependencies struct {
	Secrets    map[types.NamespacedName]*corev1.Secret
	ConfigMaps map[types.NamespacedName]*corev1.ConfigMap
}

func (e *externalDependencies) fetch(ctx context.Context, client client.Client) error {
	logger := logging.FromContext(ctx).WithValues("component", "runtime_config_builder")

	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(10)

	for name, secret := range e.Secrets {
		name := name
		secret := secret

		g.Go(func() error {
			logger.Debug("fetching secret", "secret", name)
			return client.Get(ctx, name, secret)
		})
	}

	for name, cm := range e.ConfigMaps {
		name := name
		cm := cm

		g.Go(func() error {
			logger.Debug("fetching config map", "config_map", name)
			return client.Get(ctx, name, cm)
		})
	}

	return g.Wait()
}

func extractExternalDependencies(app *spinv1.SpinApp) *externalDependencies {
	result := &externalDependencies{
		Secrets:    make(map[types.NamespacedName]*corev1.Secret),
		ConfigMaps: make(map[types.NamespacedName]*corev1.ConfigMap),
	}

	runtimeConfig := app.Spec.RuntimeConfig
	if runtimeConfig.LoadFromSecret != "" {
		// TODO: Should we block on the runtime config secret for consistency?
		return result
	}

	var configOptions []spinv1.RuntimeConfigOption

	if runtimeConfig.LLMCompute != nil {
		configOptions = append(configOptions, runtimeConfig.LLMCompute.Options...)
	}
	for _, kvStore := range runtimeConfig.KeyValueStores {
		configOptions = append(configOptions, kvStore.Options...)
	}
	for _, sqlDB := range runtimeConfig.SqliteDatabases {
		configOptions = append(configOptions, sqlDB.Options...)
	}

	secretMapper := func(configOption spinv1.RuntimeConfigOption) *corev1.Secret {
		if configOption.ValueFrom == nil {
			return nil
		}

		if configOption.ValueFrom.SecretKeyRef != nil {
			return &corev1.Secret{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Secret",
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Namespace: app.ObjectMeta.Namespace,
					Name:      configOption.ValueFrom.SecretKeyRef.Name,
				},
			}
		}

		return nil
	}

	configMapMapper := func(configOption spinv1.RuntimeConfigOption) *corev1.ConfigMap {
		if configOption.ValueFrom == nil {
			return nil
		}

		if configOption.ValueFrom.ConfigMapKeyRef != nil {
			return &corev1.ConfigMap{
				TypeMeta: metav1.TypeMeta{
					Kind:       "ConfigMap",
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Namespace: app.ObjectMeta.Namespace,
					Name:      configOption.ValueFrom.ConfigMapKeyRef.Name,
				},
			}
		}

		return nil
	}

	// 1. Map lists of config options to secrets or config maps
	// 2. remove any `nil` entries because our mapper doesn't do that by default
	// 3. Turn the list into a map of namespaced name: secret/config map

	result.Secrets = generics.AssociateBy(slices.DeleteFunc(generics.MapList(configOptions, secretMapper), func(sec *corev1.Secret) bool {
		return sec == nil
	}), func(sec *corev1.Secret) types.NamespacedName {
		return types.NamespacedName{Name: sec.ObjectMeta.Name, Namespace: sec.ObjectMeta.Namespace}
	})

	result.ConfigMaps = generics.AssociateBy(slices.DeleteFunc(generics.MapList(configOptions, configMapMapper), func(cm *corev1.ConfigMap) bool {
		return cm == nil
	}), func(cm *corev1.ConfigMap) types.NamespacedName {
		return types.NamespacedName{Name: cm.ObjectMeta.Name, Namespace: cm.ObjectMeta.Namespace}
	})

	return result
}
