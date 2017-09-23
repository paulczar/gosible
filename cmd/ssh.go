package cmd

import (
  //"os/exec"
  "os"
  "fmt"
  "path/filepath"
  "github.com/spf13/cobra"
  "github.com/paulczar/gosible/ssh"
)

var (
	sshOptions = &ssh.Options{}
	environment string
)

// runCmd represents the run command
var sshCmd = &cobra.Command{
  Use:   "ssh [flags] hostname",
  Short: "ssh to host",
  Long: `
Gosible ssh will ssh to the host name provided
`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var (
			sshConfigFile string
			err error
		)
		if len(args) < 1 {
			return fmt.Errorf("must specify a host to ssh to")
	 	}
		sshOptions.Host = args[0]
		// TODO allow commands ?

		if sshOptions.SSHConfigFile == "" {
			// Check if there's an ssh config in provided environment
			if environment != "" {
				sshConfigFile = filepath.Join(environment, "ssh_config")
				if _, err = os.Stat(sshConfigFile); err == nil {
					sshOptions.SSHConfigFile = sshConfigFile
				}
			// Check if there's an ssh config in current working dir
			} else {
				cwd, _ := os.Getwd()
				sshConfigFile = filepath.Join(cwd, "ssh_config")
				if _, err = os.Stat(sshConfigFile); err == nil {
					sshOptions.SSHConfigFile = sshConfigFile
				}
			}
		}
		return nil
	},
  RunE: func(cmd *cobra.Command, args []string) error {
    err := ssh.SSH(sshOptions)
    if err != nil {
      return err
    }
    return nil
  },
}

func init() {
  RootCmd.AddCommand(sshCmd)
  sshCmd.Flags().StringVarP(&sshOptions.SSHConfigFile, "ssh-config-file",
    "s", "", "Path to ssh config file to use.")
	sshCmd.Flags().StringVarP(&environment, "environment",
    "e", "", "directory that contains ansible inventory")
}
