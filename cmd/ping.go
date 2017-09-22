package cmd

import (
  "os/exec"
  "os"
  "fmt"
  "path/filepath"
  "github.com/spf13/cobra"
  "github.com/paulczar/gosible/ansible"
)

var pingOptions = &ansible.Options{}

// runCmd represents the run command
var pingCmd = &cobra.Command{
  Use:   "ping [flags] [--] [ansible arguments]",
  Short: "check if all hosts are available",
  Long: `
Gosible ping uses the ansible ping module to check if all hosts are available.
This is useful for checking if your machine can successfully SSH/Ansible to
each host.
`,
  PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
    var (
      err error
      virtualEnv string
    )
    // check if there's a virtualenv we should use in your cwd
    cwd, _ := os.Getwd()
    virtualEnv = filepath.Join(cwd, "virtualenv/bin")
    if _, err = os.Stat(virtualEnv); err == nil {
      os.Setenv("PATH", fmt.Sprintf("%s:%s", virtualEnv, os.Getenv("PATH")))
    }

    // check if there's a virtualenv we should use in your environment
    if pingOptions.Environment != cwd {
      virtualEnv = filepath.Join(pingOptions.Environment, "virtualenv/bin")
      if _, err = os.Stat(virtualEnv); err == nil {
        os.Setenv("PATH", fmt.Sprintf("%s:%s", virtualEnv, os.Getenv("PATH")))
      }
    }
    // check if ansible-playbook binary exists
    _, err = exec.LookPath("ansible")
    if err != nil {
      return err
    }
    return nil
  },
  RunE: func(cmd *cobra.Command, args []string) error {
    pingOptions.Module = "ping"
    err := ansible.Module(pingOptions, args)
    if err != nil {
      return err
    }
    return nil
  },
}

func init() {
  RootCmd.AddCommand(pingCmd)
  // stops parsing flags after first unknown flag is found
  pingCmd.Flags().SetInterspersed(false)
  pingCmd.Flags().StringVarP(&pingOptions.SSHConfigFile, "ssh-config-file",
    "s", "", "Path to ssh config file to use.")
  pingCmd.Flags().StringVarP(&pingOptions.Environment, "environment",
    "e", "", "directory that contains ansible inventory")
  pingCmd.Flags().StringVarP(&pingOptions.KnownHostsFile, "known-hosts-file",
    "", "", "location of known hosts file")
  pingCmd.Flags().StringVarP(&pingOptions.ModuleHosts, "hosts",
  "", "all", "host or host pattern to run")
}
