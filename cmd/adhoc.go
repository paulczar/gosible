package cmd

import (
  "os"
  "fmt"
  "path/filepath"
  "os/exec"
  "github.com/spf13/cobra"
  "github.com/paulczar/gosible/ansible"
)

var adhocOptions = &ansible.Options{}

// runCmd represents the run command
var adhocCmd = &cobra.Command{
  Use:   "adhoc [flags] command [ansible arguments]",
  Short: "wrapper around ansible command",
  Long: `
Gosible adhoc is a wrapper around ansible --module shell that adds some
additional useful features.
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
    adhocOptions.Module = "raw"
    if len(args) < 1 {
      return fmt.Errorf("must specify an adhoc command to run")
    }
    adhocOptions.ModuleArgs = args[0]
    args = args[1:]
    err := ansible.Module(adhocOptions, args)
    if err != nil {
      return err
    }
    return nil
  },
}

func init() {
  RootCmd.AddCommand(adhocCmd)
  // stops parsing flags after first unknown flag is found
  adhocCmd.Flags().SetInterspersed(false)
  adhocCmd.Flags().StringVarP(&adhocOptions.SSHConfigFile, "ssh-config-file",
    "s", "", "Path to ssh config file to use.")
  adhocCmd.Flags().BoolVarP(&adhocOptions.SSHForwardAgent, "ssh-forward-agent",
    "f", false, "path to ssh config file to use")
  adhocCmd.Flags().StringVarP(&adhocOptions.Environment, "environment",
    "e", "", "directory that contains ansible inventory")
  adhocCmd.Flags().StringVarP(&adhocOptions.KnownHostsFile, "known-hosts-file",
    "", "", "location of known hosts file")
  adhocCmd.Flags().StringVarP(&adhocOptions.ModuleHosts, "hosts",
  "", "all", "host or host pattern to run")
}
