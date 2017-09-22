package cmd

import (
  "os"
  "os/exec"
  "fmt"
  "github.com/spf13/cobra"
  "github.com/paulczar/gosible/ansible"
  "github.com/paulczar/gosible/provisioner"
  "path/filepath"
)

// runCmd represents the run command
var playbookRunCmd = &cobra.Command{
  Use:   "run [flags] file.yml [--] [ansible-playbook arguments]",
  Short: "wrapper around ansible command",
  Long: `
Gosible playbook run will run an ansible playbook over your provided environment.
You can set additional ansible-playbook flags by providing a double dash "--" and 
then the additional flags.

Examples:
Run playbook/yaml over the env/test environment passing flags to ansible to use 
the user "root" and run in verbose mode:
    $ gosible playbook run -e env/test playbook.yml -- --user=root -vvvv
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
    if playbookOptions.Environment != cwd {
      virtualEnv = filepath.Join(playbookOptions.Environment, "virtualenv/bin")
      if _, err = os.Stat(virtualEnv); err == nil {
        os.Setenv("PATH", fmt.Sprintf("%s:%s", virtualEnv, os.Getenv("PATH")))
      }
    }
    
    // check if ansible-playbook binary exists
    _, err = exec.LookPath("ansible-playbook")
    if err != nil {
      return err
    }
    // check if provisioning needs to happen
    if po.Provisioner != "" {
      po.Environment = playbookOptions.Environment
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
    err := ansible.PlaybookRun(playbookOptions, args)
    if err != nil {
      return err
    }
    return nil
  },
}

func init() {
  playbookCmd.AddCommand(playbookRunCmd)
  playbookRunCmd.Flags().StringVarP(&playbookOptions.SSHConfigFile, "ssh-config-file",
    "s", "", "Path to ssh config file to use.")
  playbookRunCmd.Flags().BoolVarP(&playbookOptions.SSHForwardAgent, "ssh-forward-agent",
    "f", false, "path to ssh config file to use")
  playbookRunCmd.Flags().StringVarP(&playbookOptions.Environment, "environment",
    "e", "", "directory that contains ansible inventory")
  playbookRunCmd.Flags().StringVarP(&playbookOptions.KnownHostsFile, "known-hosts-file",
    "", "", "location of known hosts file")
  playbookRunCmd.Flags().StringVarP(&po.Provisioner, "provisioner",
    "", "", "provisioner (vagrant)")
}
