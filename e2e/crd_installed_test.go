package e2e

import (
	"context"
	"testing"

	apiextensionsV1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"
)

func TestCRDInstalled(t *testing.T) {
	crdInstalledFeature := features.New("crd installed").
		Assess("spinapp crd installed", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			client := cfg.Client()
			if err := apiextensionsV1.AddToScheme(client.Resources().GetScheme()); err != nil {
				t.Fatalf("failed to register the v1 API extension types with Kubernetes scheme: %s", err)
			}
			name := "spinapps.core.spinkube.dev"
			var crd apiextensionsV1.CustomResourceDefinition
			if err := client.Resources().Get(ctx, name, "", &crd); err != nil {
				t.Fatalf("SpinApp CRD not found: %s", err)
			}

			if crd.Spec.Group != "core.spinkube.dev" {
				t.Fatalf("SpinApp CRD has unexpected group: %s", crd.Spec.Group)
			}
			return ctx

		}).
		Assess("spinappexecutor crd installed", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			client := cfg.Client()
			if err := apiextensionsV1.AddToScheme(client.Resources().GetScheme()); err != nil {
				t.Fatalf("failed to register the v1 API extension types with Kubernetes scheme: %s", err)
			}

			name := "spinappexecutors.core.spinkube.dev"
			var crd apiextensionsV1.CustomResourceDefinition
			if err := client.Resources().Get(ctx, name, "", &crd); err != nil {
				t.Fatalf("SpinApp CRD not found: %s", err)
			}

			if crd.Spec.Group != "core.spinkube.dev" {
				t.Fatalf("SpinAppExecutor CRD has unexpected group: %s", crd.Spec.Group)
			}
			return ctx
		}).Feature()
	testEnv.Test(t, crdInstalledFeature)
}
