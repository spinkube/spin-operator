package e2e

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	nodev1 "k8s.io/api/node/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/e2e-framework/klient/wait"
	"sigs.k8s.io/e2e-framework/klient/wait/conditions"
	"sigs.k8s.io/e2e-framework/pkg/env"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/envfuncs"
	"sigs.k8s.io/e2e-framework/support/utils"

	spinapps_v1alpha1 "github.com/spinkube/spin-operator/api/v1alpha1"
)

const ErrFormat = "%v: %v\n"

var (
	testEnv                    env.Environment
	testNamespace              string
	testCACertSecret           = "test-spin-ca"
	spinOperatorDeploymentName = "spin-operator-controller-manager"
	spinOperatorNamespace      = "spin-operator"
	cluster                    = &Cluster{}
)

func TestMain(m *testing.M) {
	cfg, _ := envconf.NewFromFlags()
	testEnv = env.NewWithConfig(cfg)
	testNamespace = envconf.RandomName("my-ns", 10)
	cluster.name = envconf.RandomName("crdtest-", 16)

	testEnv.Setup(
		func(ctx context.Context, e *envconf.Config) (context.Context, error) {
			if _, err := cluster.Create(ctx, cluster.name); err != nil {
				return ctx, err
			}
			e.WithKubeconfigFile(cluster.kubecfgFile)

			return ctx, nil
		},

		envfuncs.CreateNamespace(testNamespace),

		// build and load spin operator image into cluster
		func(ctx context.Context, _ *envconf.Config) (context.Context, error) {
			if os.Getenv("E2E_SKIP_BUILD") == "" { // nolint:forbidigo
				if p := utils.RunCommand(`bash -c "cd .. && IMG=ghcr.io/spinkube/spin-operator:dev make docker-build"`); p.Err() != nil {
					return ctx, fmt.Errorf(ErrFormat, p.Err(), p.Out())
				}
			}

			if p := utils.RunCommand(("k3d image import -c " + cluster.name + " ghcr.io/spinkube/spin-operator:dev")); p.Err() != nil {
				return ctx, fmt.Errorf(ErrFormat, p.Err(), p.Out())
			}
			return ctx, nil
		},

		// install spin operator and pre-reqs
		func(ctx context.Context, _ *envconf.Config) (context.Context, error) {
			// install crds
			if p := utils.RunCommand(`bash -c "cd .. && make install"`); p.Err() != nil {

				return ctx, fmt.Errorf(ErrFormat, p.Err(), p.Out())

			}

			// install cert-manager
			if p := utils.RunCommand("kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.14.2/cert-manager.yaml"); p.Err() != nil {
				return ctx, fmt.Errorf(ErrFormat, p.Err(), p.Out())
			}
			// wait for cert-manager to be ready
			if p := utils.RunCommand("kubectl wait --for=condition=Available --timeout=300s deployment/cert-manager-webhook -n cert-manager"); p.Err() != nil {
				return ctx, fmt.Errorf(ErrFormat, p.Err(), p.Out())
			}

			if p := utils.RunCommand(`bash -c "cd .. && IMG=ghcr.io/spinkube/spin-operator:dev make deploy"`); p.Err() != nil {
				return ctx, fmt.Errorf(ErrFormat, p.Err(), p.Out())
			}

			// wait for the controller deployment to be ready
			client := cfg.Client()

			if err := spinapps_v1alpha1.AddToScheme(client.Resources().GetScheme()); err != nil {
				return ctx, fmt.Errorf("failed to register the spinapps_v1alpha1 types with Kubernetes scheme: %w", err)
			}

			if err := wait.For(
				conditions.New(client.Resources()).
					DeploymentAvailable(spinOperatorDeploymentName, spinOperatorNamespace),
				wait.WithTimeout(3*time.Minute),
				wait.WithInterval(10*time.Second),
			); err != nil {
				return ctx, err
			}

			return ctx, nil
		},

		// create runtime class
		func(ctx context.Context, c *envconf.Config) (context.Context, error) {
			client := cfg.Client()
			runtimeClass := &nodev1.RuntimeClass{
				ObjectMeta: metav1.ObjectMeta{
					Name: runtimeClassName,
				},
				Handler: "spin",
			}

			err := client.Resources().Create(ctx, runtimeClass)
			return ctx, err
		},

		// create executor
		func(ctx context.Context, c *envconf.Config) (context.Context, error) {
			client := cfg.Client()

			err := client.Resources().Create(ctx, newContainerdShimExecutor(testNamespace))
			return ctx, err
		},
	)

	testEnv.Finish(
		func(ctx context.Context, _ *envconf.Config) (context.Context, error) {
			return ctx, cluster.Destroy()
		},
	)

	os.Exit(testEnv.Run(m))
}

func newContainerdShimExecutor(namespace string) *spinapps_v1alpha1.SpinAppExecutor {
	return &spinapps_v1alpha1.SpinAppExecutor{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "containerd-shim-spin",
			Namespace: namespace,
		},
		Spec: spinapps_v1alpha1.SpinAppExecutorSpec{
			CreateDeployment: true,
			DeploymentConfig: &spinapps_v1alpha1.ExecutorDeploymentConfig{
				RuntimeClassName:      runtimeClassName,
				InstallDefaultCACerts: true,
				CACertSecret:          testCACertSecret,
			},
		},
	}
}
