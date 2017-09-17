package cmd

import (
  //"os/exec"
  "github.com/spf13/cobra"
  "github.com/paulczar/gosible/ansible"
)

var installOptions = &ansible.InstallOptions{}

// runCmd represents the run command
var installCmd = &cobra.Command{
  Use:   "install",
  Short: "install ansible via pip",
  Long: `
Gosible install will install ansible and/or any dependencies provided.
If you do not specify a virtualenv or requirements file it will attempt
to find them in your local path, or it will try to find them in your
environment if you provide one.
`,
  RunE: func(cmd *cobra.Command, args []string) error {
    err := ansible.InstallViaPip(installOptions)
    if err != nil {
      return err
    }
    return nil
  },
}

func init() {
  RootCmd.AddCommand(installCmd)
  // stops parsing flags after first unknown flag is found
  installCmd.Flags().SetInterspersed(false)
  installCmd.Flags().StringVarP(&installOptions.VirtualEnv, "virtualenv",
    "v", "", "Path to VirtualEnv to use")
  installCmd.Flags().StringVarP(&installOptions.RequirementsTXT, "requirements",
    "r", "", "path to requirements.txt")
		installCmd.Flags().StringVarP(&installOptions.Path, "environment",
			"e", "", "path to your gosible environment")
	}
