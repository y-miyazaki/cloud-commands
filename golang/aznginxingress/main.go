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
	app.Usage = `This command installs nginx ingress.
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
`
	app.Version = "v1.1.0"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "namespace, n",
			Usage: "The namespace of nginx ingress",
			Value: "default",
		},
		cli.StringFlag{
			Name:  "release-name, r",
			Usage: "The release name of nginx ingress",
			Value: "nginx-ingress",
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
			Name:  "other-options, o",
			Usage: "Specify the helm install option to set Nginx Ingress at the time of installation.",
		},
		cli.BoolFlag{
			Name:  "install, i",
			Usage: "install nginx ingress.",
		},
		cli.BoolFlag{
			Name:  "uninstall, u",
			Usage: "uninstall nginx ingress.",
		},
	}
	app.Action = func(c *cli.Context) error {
		installNginxIngressFlag := c.Bool("install")
		uninstallNginxIngress := c.Bool("uninstall")
		if installNginxIngressFlag {
			err := installNginxIngress(c)
			if err != nil {
				return err
			}
			nginxNamespace := c.String("namespace")
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
		} else if uninstallNginxIngress {
			nginxReleaseName := c.String("release-name")
			_, _ = command.CombinedOutputStr("helm delete --purge "+nginxReleaseName, true)
		}
		return nil
	}
	return app.Run(args)
}

func installNginxIngress(c *cli.Context) error {
	nginxNamespace := c.String("namespace")
	nginxReleaseName := c.String("release-name")
	nginxReplicaCount := c.String("replicacount")
	nginxLoadBalancerIP := c.String("load-balancer-ip")
	internal := c.Bool("internal")
	nginxOtherOptions := c.String("other-options")
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

func main() {
	_ = run(os.Args)
}
