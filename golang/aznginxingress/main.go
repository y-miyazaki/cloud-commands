package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
	"github.com/y-miyazaki/go-common/pkg/command"
)

func run(args []string) error {
	app := cli.NewApp()
	app.Name = "aznginxingress"
	app.Usage = `This command installs nginx ingress and cert-manager.
   This command requires "kubectl" and "helm" commands.
   Also, since you actually install to the Cluster, you need to have permission to deploy to the Cluster.

   1. This is the base installation of nginx-ingress executed by this command.
	  If you want to specify additional options, use the load-balancer-ip / internal / nginx-other-options option.

      helm install stable/nginx-ingress \
         --set controller.replicaCount=2 \
         --set controller.nodeSelector."beta\.kubernetes\.io/os"=linux \
         --set defaultBackend.nodeSelector."beta\.kubernetes\.io/os"=linux \
         --set controller.service.externalTrafficPolicy=Local \
         --namespace {nginxNamespace} \
         --name {nginxReleaseName}

      Reference document
         https://docs.microsoft.com/ja-jp/azure/aks/ingress-static-ip#create-an-ingress-controller
         https://github.com/helm/charts/tree/master/stable/nginx-ingress

   2. This is the base installation of cert-manager executed by this command.
      If you want to specify additional options, use the cert-manager-other-options option.

      helm install --name {certManagerReleaseName} \
         --namespace {certManagerNamespace} \
         --version v{certManagerVersion}.0 \
         jetstack/cert-manager

      Reference document
         https://docs.microsoft.com/ja-jp/azure/aks/ingress-static-ip#install-cert-manager
         https://github.com/jetstack/cert-manager
`
	app.Version = "v1.0.1"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "nginx-release-name, n-rn",
			Usage: "The release name of nginx ingress",
			Value: "nginx-ingress",
		},
		cli.StringFlag{
			Name:  "nginx-namespace, n-n",
			Usage: "The namespace of nginx ingress",
			Value: "default",
		},
		cli.IntFlag{
			Name:  "replicacount",
			Usage: "The replica count of nginx ingress",
			Value: 2,
		},
		cli.StringFlag{
			Name:  "load-balancer-ip, lbip",
			Usage: "If set, add the option when installing nginx ingress and allocate only fixed IP.",
		},
		cli.BoolFlag{
			Name:  "internal",
			Usage: "If set, add the option when installing nginx ingress and allocate only the internal IP.",
		},
		cli.StringFlag{
			Name:  "nginx-other-options",
			Usage: "Specify the helm install option to set Nginx Ingress at the time of installation.",
		},
		cli.StringFlag{
			Name:  "cert-manager-version, cm-v",
			Usage: "The version of cert-manager",
			Value: "0.11",
		},
		cli.StringFlag{
			Name:  "cert-manager-namespace, cm-namespace",
			Usage: "The namespace of cert-manager",
			Value: "cert-manager",
		},
		cli.StringFlag{
			Name:  "cert-manager-release-name, cm-release-name",
			Usage: "The release name of cert-manager",
			Value: "cert-manager",
		},
		cli.StringFlag{
			Name:  "cert-manager-other-options",
			Usage: "Specify the helm install option to set cert manager at the time of installation.",
		},
		cli.BoolFlag{
			Name:  "install-nginxingress, i-n",
			Usage: "install nginx ingress.",
		},
		cli.BoolFlag{
			Name:  "install-cert-manager, i-c",
			Usage: "install cert-manager.",
		},
		cli.BoolFlag{
			Name:  "uninstall-nginxingress, ui-n",
			Usage: "uninstall nginx ingress and cert-manager.",
		},
		cli.BoolFlag{
			Name:  "uninstall-cert-manager, ui-c",
			Usage: "uninstall nginx ingress and cert-manager.",
		},
	}
	app.Action = func(c *cli.Context) error {
		installNginxIngressFlag := c.Bool("install-nginxingress")
		installCertManagerFlag := c.Bool("install-cert-manager")
		uninstallNginxIngress := c.Bool("uninstall-nginxingress")
		uninstallCertManager := c.Bool("uninstall-cert-manager")
		if installNginxIngressFlag {
			err := installNginxIngress(c)
			if err != nil {
				return err
			}
		}
		if installCertManagerFlag {
			err := installCertManager(c)
			if err != nil {
				return err
			}
		}
		if installNginxIngressFlag {
			nginxNamespace := c.String("nginx-namespace")
			fmt.Println("#------------------------------------------------------------------------")
			fmt.Println("# Nginx Ingress associates ip address now.")
			fmt.Println("# Please check this following command.")
			fmt.Println("#------------------------------------------------------------------------")
			fmt.Println("kubectl get service -l app=nginx-ingress --namespace " + nginxNamespace)

			fmt.Println("#------------------------------------------------------------------------")
			fmt.Println("# describe Certificate.")
			fmt.Println("# Please check this following command.")
			fmt.Println("#------------------------------------------------------------------------")
			fmt.Println("kubectl describe certificate {your secret name} --namespace " + nginxNamespace)
		}
		if uninstallNginxIngress {
			nginxReleaseName := c.String("nginx-release-name")
			_, _ = command.CombinedOutputStr("helm delete --purge "+nginxReleaseName, true)
		}
		if uninstallCertManager {
			certManagerReleaseName := c.String("cert-manager-release-name")
			_, _ = command.CombinedOutputStr("helm delete --purge "+certManagerReleaseName, true)
		}
		return nil
	}
	return app.Run(args)
}

func installNginxIngress(c *cli.Context) error {
	nginxNamespace := c.String("nginx-namespace")
	nginxReleaseName := c.String("nginx-release-name")
	nginxReplicaCount := c.String("replicacount")
	nginxLoadBalancerIP := c.String("load-balancer-ip")
	internal := c.Bool("internal")
	nginxOtherOptions := c.String("nginx-other-options")
	// ------------------------------------------------------------------------
	// nginx ingress check namespace & create namespace
	// ------------------------------------------------------------------------
	commandStr := "kubectl get ns | grep " + nginxNamespace + " | wc -l | tr -d \"\\n\""
	out, err := command.CombinedOutputStr(commandStr, true)
	if err != nil {
		return err
	}
	if out == "0" {
		_, err = command.CombinedOutput("kubectl", true, "create", "namespace", nginxNamespace)
		if err != nil {
			return err
		}
	}
	// ------------------------------------------------------------------------
	// Use Helm to deploy an NGINX ingress controller
	// https://docs.microsoft.com/ja-jp/azure/aks/ingress-tls#create-an-ingress-controller
	// ------------------------------------------------------------------------
	commandStr = "helm list | grep " + nginxReleaseName + " | wc -l | tr -d \"\\n\""
	out, err = command.CombinedOutputStr(commandStr, true)
	if err != nil {
		return err
	}
	if out == "0" {
		commandStr = "helm install stable/nginx-ingress"
		commandStr += " --set controller.replicaCount=" + nginxReplicaCount
		commandStr += " --set controller.nodeSelector.\"beta\\.kubernetes\\.io/os\"=linux"
		commandStr += " --set defaultBackend.nodeSelector.\"beta\\.kubernetes\\.io/os\"=linux"
		commandStr += " --set controller.service.externalTrafficPolicy=Local"
		commandStr = commandStr + " --namespace " + nginxNamespace
		commandStr = commandStr + " --name " + nginxReleaseName
		if internal {
			commandStr += " --set-string controller.service.annotations.\"service\\.beta\\.kubernetes\\.io/azure-load-balancer-internal\"=true"
		}
		if nginxLoadBalancerIP != "" {
			commandStr += " --set controller.service.loadBalancerIP=\"" + nginxLoadBalancerIP + "\""
		}
		// other options
		commandStr = commandStr + " " + nginxOtherOptions

		_, err = command.CombinedOutputStr(commandStr, true)
		if err != nil {
			return err
		}
	}
	return nil
}

func installCertManager(c *cli.Context) error {
	certManagerVersion := c.String("cert-manager-version")
	certManagerNamespace := c.String("cert-manager-namespace")
	certManagerReleaseName := c.String("cert-manager-release-name")
	certManagerOtherOptions := c.String("cert-manager-other-options")
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
