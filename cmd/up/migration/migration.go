package migration

import (
	"github.com/alecthomas/kong"
	"github.com/upbound/up/internal/kube"
	"github.com/upbound/up/internal/migration"
	"github.com/upbound/up/internal/upbound"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

// AfterApply constructs and binds Upbound-specific context to any subcommands
// that have Run() methods that receive it.
func (c *Cmd) AfterApply(kongCtx *kong.Context) error {
	upCtx, err := upbound.NewFromFlags(c.Flags)
	if err != nil {
		return err
	}
	cfg, err := kube.GetKubeConfig(c.Kubeconfig)
	if err != nil {
		return err
	}
	if upCtx.WrapTransport != nil {
		cfg.Wrap(upCtx.WrapTransport)
	}
	client, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return err
	}
	d, err := dynamic.NewForConfig(cfg)
	if err != nil {
		return err
	}

	kongCtx.Bind(&migration.Context{
		KubeCore:    client.CoreV1(),
		KubeDynamic: d,
		Namespace:   c.Namespace,
	})
	return nil
}

type Cmd struct {
	Export exportCmd `cmd:"" help:"Export a control plane."`
	Import importCmd `cmd:"" help:"Import a control plane."`

	Kubeconfig string `type:"existingfile" help:"Override default kubeconfig path."`
	Namespace  string `short:"n" env:"UXP_NAMESPACE" default:"upbound-system" help:"Kubernetes namespace for UXP."`

	// Common Upbound API configuration
	Flags upbound.Flags `embed:""`
}
