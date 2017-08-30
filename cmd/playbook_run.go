package cmd

import (
	//"os"
	"os/exec"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/paulczar/gosible/ansible/playbook"
	"github.com/paulczar/gosible/provisioner"
)

// runCmd represents the run command
var playbookRunCmd = &cobra.Command{
	Use:   "run [flags] file.yml [ansible-playbook arguments]",
	Short: "wrapper around ansible command",
	Long: `
Gosible playbook is a wrapper around ansible-playbook that adds some
additional useful features.
`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var (
			err error
		)
		// check if ansible-playbook binary exists
		_, err = exec.LookPath("ansible-playbook")
		if err != nil {
			return err
		}
		// check if provisioning needs to happen
		if po.Provisioner != "" {
			po.Environment = ao.Environment
			err = provisioner.Up(po)
			if err != nil {
				return err
			}
		}


		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
 			return fmt.Errorf("must specify a playbook to run")
    }
		err := playbook.Run(ao, args)
		if err != nil {
			return err
		}
		fmt.Println("Successfully ran ansible?")
		return nil
	},
}

func init() {
	playbookCmd.AddCommand(playbookRunCmd)
	// stops parsing flags after first unknown flag is found
	playbookRunCmd.Flags().SetInterspersed(false)
	playbookRunCmd.Flags().StringVarP(&ao.SSHConfigFile, "ssh-config-file",
		"s", "", "Path to ssh config file to use.")
	playbookRunCmd.Flags().BoolVarP(&ao.SSHForwardAgent, "ssh-forward-agent",
		"f", false, "path to ssh config file to use")
	playbookRunCmd.Flags().StringVarP(&ao.Environment, "environment",
		"e", "", "directory that contains ansible inventory")
	playbookRunCmd.Flags().StringVarP(&ao.KnownHostsFile, "known-hosts-file",
		"", "", "location of known hosts file")
	playbookRunCmd.Flags().StringVarP(&po.Provisioner, "provisioner",
		"", "", "provisioner (vagrant)")
}
