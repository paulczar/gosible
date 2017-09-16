package ansible

import (
    "fmt"
)

// ModuleOptions for running ansible
type ModuleOptions struct {
  ModuleHosts     string
  Module          string
  ModulePath      string
  ModuleArgs      string
}

// Module ansible playbook
func Module(options *Options, ansibleArgs []string) error {
  var (
    err         error
    gosibleArgs []string
  )

  err = configureEnvironment(options)
  if err != nil {
    return err
  }
  err = sshConfigFile(options)
  if err != nil {
    return err
  }
  err = configureKnownHostsFile(options)
  if err != nil {
    return err
  }
  configureSSHForwardAgent(options)

  gosibleArgs = []string{options.ModuleHosts, "--module-name", options.Module}

  if options.ModulePath != "" {
    gosibleArgs = append(gosibleArgs,
    []string{"--module-path", options.ModulePath}...)
  }

  if options.ModuleArgs != "" {
    gosibleArgs = append(gosibleArgs,
    []string{"--args", fmt.Sprintf("'%s'", options.ModuleArgs)}...)
  }

  if options.Inventory != "" {
    gosibleArgs = append(gosibleArgs,
      []string{"--inventory", options.Inventory}...)
  }

  cmdName := "ansible"
  cmdArgs := append(gosibleArgs, ansibleArgs...)

  return runCmd(cmdName, cmdArgs)

}
