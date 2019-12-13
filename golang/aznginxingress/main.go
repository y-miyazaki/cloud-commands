package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/urfave/cli"
)

func run(args []string) error {
	app := cli.NewApp()
	app.Name = "aznginxingress"
	app.Usage = "This command installs nginx ingress and cert-manager."
	app.Version = "v1.0.0"
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
			Usage: "The load balancer ip address",
		},
		cli.BoolFlag{
			Name:  "internal",
			Usage: "If set, The load balancer ip is internal",
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
		cli.BoolFlag{
			Name:  "install, i",
			Usage: "install nginx ingress and cert-manager.",
		},
		cli.BoolFlag{
			Name:  "uninstall",
			Usage: "uninstall nginx ingress and cert-manager.",
		},
	}
	app.Action = func(c *cli.Context) error {
		nginxNamespace := c.String("nginx-namespace")
		nginxReleaseName := c.String("nginx-release-name")
		nginxReplicaCount := c.String("replicacount")
		nginxLoadBalancerIP := c.String("load-balancer-ip")
		internal := c.Bool("internal")
		certManagerVersion := c.String("cert-manager-version")
		certManagerNamespace := c.String("cert-manager-namespace")
		certManagerReleaseName := c.String("cert-manager-release-name")
		install := c.Bool("install")
		uninstall := c.Bool("uninstall")
		if install {
			// ------------------------------------------------------------------------
			// nginx ingress check namespace & create namespace
			// ------------------------------------------------------------------------
			commandStr := "kubectl get ns | grep " + nginxNamespace + " | wc -l | tr -d \"\\n\""
			out, err := execCommandStr(commandStr)
			if err != nil {
				return err
			}
			if string(out) == "0" {
				_, err = execCommand("kubectl", "create", "namespace", nginxNamespace)
				if err != nil {
					return err
				}
			}
			// ------------------------------------------------------------------------
			// Use Helm to deploy an NGINX ingress controller
			// https://docs.microsoft.com/ja-jp/azure/aks/ingress-tls#create-an-ingress-controller
			// ------------------------------------------------------------------------
			commandStr = "helm list | grep " + nginxReleaseName + " | wc -l | tr -d \"\\n\""
			out, err = execCommandStr(commandStr)
			if err != nil {
				return err
			}
			if string(out) == "0" {
				commandStr = "helm install stable/nginx-ingress"
				commandStr += " --set controller.replicaCount=" + nginxReplicaCount
				commandStr += " --set controller.nodeSelector.\"beta\\.kubernetes\\.io/os\"=linux"
				commandStr += " --set defaultBackend.nodeSelector.\"beta\\.kubernetes\\.io/os\"=linux"
				commandStr = commandStr + " --namespace " + nginxNamespace
				commandStr = commandStr + " --name " + nginxReleaseName
				if internal {
					// commandStr += " --set controller.service.annotations.\"service.beta.kubernetes.io/azure-load-balancer-internal\"='\"true\"'"
					commandStr += " --set-string controller.service.annotations.\"service\\.beta\\.kubernetes\\.io/azure-load-balancer-internal\"=true"
				}
				if nginxLoadBalancerIP != "" {
					commandStr += " --set controller.service.loadBalancerIP=\"" + nginxLoadBalancerIP + "\""
				}
				_, err := execCommandStr(commandStr)
				if err != nil {
					return err
				}
			}
			// ------------------------------------------------------------------------
			// Use Helm to deploy an Cert Manager
			// ------------------------------------------------------------------------
			commandStr = "helm list | grep " + certManagerReleaseName + " | wc -l | tr -d \"\\n\""
			out, err = execCommandStr(commandStr)
			if err != nil {
				return err
			}
			if out == "0" {
				// Install the CustomResourceDefinition resources separately
				// see document.
				// https://cert-manager.io/docs/installation/kubernetes/
				commandStr = "kubectl apply --validate=false -f https://raw.githubusercontent.com/jetstack/cert-manager/release-" + certManagerVersion + "/deploy/manifests/00-crds.yaml"
				_, err = execCommandStr(commandStr)
				if err != nil {
					return err
				}

				// cert manager check namespace & create namespace
				commandStr := "kubectl get ns | grep " + certManagerNamespace + " | wc -l | tr -d \"\\n\""
				out, err = execCommandStr(commandStr)
				if err != nil {
					return err
				}
				if out == "0" {
					_, err := execCommand("kubectl", "create", "namespace", certManagerNamespace)
					if err != nil {
						return err
					}
				}
				// Label the cert-manager namespace to disable resource validation
				_, err = execCommand("kubectl", "label", "namespace", "--overwrite", certManagerNamespace, "certmanager.k8s.io/disable-validation=true")
				if err != nil {
					return err
				}
				// Add the Jetstack Helm repository
				_, err = execCommandStr("helm repo add jetstack https://charts.jetstack.io")
				if err != nil {
					return err
				}
				// Update your local Helm chart repository cache
				_, err = execCommandStr("helm repo update")
				if err != nil {
					return err
				}
				// Update your local Helm chart repository cache
				_, err = execCommandStr("helm install --name " + certManagerReleaseName + " --namespace " + certManagerNamespace + " --version v" + certManagerVersion + ".0 jetstack/cert-manager")
				if err != nil {
					return err
				}
			}
			fmt.Println("# ------------------------------------------------------------------------")
			fmt.Println("# Nginx Ingress associates ip address now.")
			fmt.Println("# Please check this following command.")
			fmt.Println("# ------------------------------------------------------------------------")
			fmt.Println("kubectl get service -l app=nginx-ingress --namespace " + nginxNamespace)

			fmt.Println("# ------------------------------------------------------------------------")
			fmt.Println("# describe Certificate.")
			fmt.Println("# Please check this following command.")
			fmt.Println("# ------------------------------------------------------------------------")
			fmt.Println("kubectl describe certificate {your secret name} --namespace " + nginxNamespace)
		}
		if uninstall {
			execCommandStr("helm delete --purge " + nginxReleaseName)
			execCommandStr("helm delete --purge " + certManagerReleaseName)
		}
		return nil
	}
	return app.Run(args)
}

// execCommand execute exec.Command and output command.
func execCommand(command string, options ...string) (string, error) {
	out, err := exec.Command(command, options...).CombinedOutput()
	outputCommand := command + " "
	for _, s := range options {
		outputCommand = s + " "
	}
	fmt.Println(outputCommand)
	fmt.Println(string(out))
	return string(out), err
}

// execCommandStr execute exec.Command and output command.
func execCommandStr(command string) (string, error) {
	out, err := exec.Command("sh", "-c", command).CombinedOutput()
	fmt.Println(command)
	fmt.Println(string(out))
	return string(out), err
}

func main() {
	run(os.Args)
}
