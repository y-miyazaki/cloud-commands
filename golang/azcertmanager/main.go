package main

import (
	"os"

	"github.com/urfave/cli"
	"github.com/y-miyazaki/go-common/pkg/command"
)

func run(args []string) error {
	app := cli.NewApp()
	app.Name = "azcertmanager"
	app.Usage = `This command installs cert-manager.
   This command requires "kubectl" and "helm" commands.
   Also, since you actually install to the Cluster, you need to have permission to deploy to the Cluster.

   1. This is the base installation of cert-manager executed by this command.
      If you want to specify additional options, use the other-options option.

      helm install --name {certManagerReleaseName} \
         --namespace {certManagerNamespace} \
         --version v{certManagerVersion}.0 \
         jetstack/cert-manager

      Reference document
         https://docs.microsoft.com/ja-jp/azure/aks/ingress-static-ip#install-cert-manager
         https://github.com/jetstack/cert-manager
`
	app.Version = "v1.1.0"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "version, v",
			Usage: "The version of cert-manager",
			Value: "0.11",
		},
		cli.StringFlag{
			Name:  "namespace, n",
			Usage: "The namespace of cert-manager",
			Value: "cert-manager",
		},
		cli.StringFlag{
			Name:  "release-name, r",
			Usage: "The release name of cert-manager",
			Value: "cert-manager",
		},
		cli.StringFlag{
			Name:  "other-options, o",
			Usage: "Specify the helm install option to set cert manager at the time of installation.",
		},
		cli.BoolFlag{
			Name:  "install, i",
			Usage: "install cert-manager.",
		},
		cli.BoolFlag{
			Name:  "uninstall, u",
			Usage: "uninstall nginx ingress and cert-manager.",
		},
	}
	app.Action = func(c *cli.Context) error {
		installCertManagerFlag := c.Bool("install")
		uninstallCertManager := c.Bool("uninstall")
		if installCertManagerFlag {
			err := installCertManager(c)
			if err != nil {
				return err
			}
		} else if uninstallCertManager {
			certManagerReleaseName := c.String("release-name")
			_, _ = command.CombinedOutputStr("helm delete --purge "+certManagerReleaseName, true)
		}
		return nil
	}
	return app.Run(args)
}

func installCertManager(c *cli.Context) error {
	certManagerVersion := c.String("version")
	certManagerNamespace := c.String("namespace")
	certManagerReleaseName := c.String("release-name")
	certManagerOtherOptions := c.String("other-options")
	// ------------------------------------------------------------------------
	// Use Helm to deploy an Cert Manager
	// ------------------------------------------------------------------------
	commandStr := "helm list | grep " + certManagerReleaseName + " | wc -l | tr -d \"\\n\""
	out, err := command.CombinedOutputStr(commandStr, true)
	if err != nil {
		return err
	}
	if out == "0" {
		// Install the CustomResourceDefinition resources separately
		// see document.
		// https://cert-manager.io/docs/installation/kubernetes/
		commandStr = "kubectl apply --validate=false -f https://raw.githubusercontent.com/jetstack/cert-manager/release-" + certManagerVersion + "/deploy/manifests/00-crds.yaml"
		_, err = command.CombinedOutputStr(commandStr, true)
		if err != nil {
			return err
		}

		// cert manager check namespace & create namespace
		commandStr = "kubectl get ns | grep " + certManagerNamespace + " | wc -l | tr -d \"\\n\""
		out, err = command.CombinedOutputStr(commandStr, true)
		if err != nil {
			return err
		}
		if out == "0" {
			_, err = command.CombinedOutput("kubectl", true, "create", "namespace", certManagerNamespace)
			if err != nil {
				return err
			}
		}
		// Label the cert-manager namespace to disable resource validation
		_, err = command.CombinedOutput("kubectl", true, "label", "namespace", "--overwrite", certManagerNamespace, "cert-manager.io/disable-validation=true")
		if err != nil {
			return err
		}
		// Add the Jetstack Helm repository
		_, err = command.CombinedOutputStr("helm repo add jetstack https://charts.jetstack.io", true)
		if err != nil {
			return err
		}
		// Update your local Helm chart repository cache
		_, err = command.CombinedOutputStr("helm repo update", true)
		if err != nil {
			return err
		}
		// Update your local Helm chart repository cache
		_, err = command.CombinedOutputStr("helm install --name "+certManagerReleaseName+" --namespace "+certManagerNamespace+" --version v"+certManagerVersion+".0 "+certManagerOtherOptions+" jetstack/cert-manager", true)
		if err != nil {
			return err
		}
	}
	return nil
}
func main() {
	_ = run(os.Args)
}
