module github.com/y-miyazaki/vscode-remote-development-terraform-azure/cmd/aznginxingress

go 1.13

require (
	github.com/urfave/cli v1.22.2
	github.com/y-miyazaki/cloud-commands/pkg/command v0.0.0
)

replace github.com/y-miyazaki/cloud-commands/pkg/command => ../../pkg/command
