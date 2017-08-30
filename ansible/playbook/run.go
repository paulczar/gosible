package playbook

import (
  	"fmt"
    "os"
  	"os/exec"
    "strings"
    "path/filepath"
)

// Options for running ansible playbook
type Options struct {
	SSHConfigFile   string
  SSHForwardAgent bool
  Provisioner     string
  Environment     string
  Inventory       string
  KnownHostsFile  string
}

// sets an environment variable
func setEnvironmentVariables(variable, value string) {
  os.Setenv(variable, value)
}

// appends a string to an existing environment variable
func appendEnvironmentVariable(variable, value string) {
  old := os.Getenv(variable)
  new := []string{old,value}
  os.Setenv(variable, strings.Join(new, " "))
}

// Checks if specified ssh config file exists, if not it
// checks if there is an ssh_config in the specified environment
// otherwise defaults to no ssh config.
// appends result to env var ANSIBLE_SSH_ARGS
func sshConfigFile(options *Options) error {
  if options.SSHConfigFile == "" {
    sshConfigFile := filepath.Join(options.Environment, "ssh_config")
    if _, err := os.Stat(sshConfigFile); !os.IsNotExist(err) {
      options.SSHConfigFile = sshConfigFile
    }
  } else {
    if _, err := os.Stat(options.SSHConfigFile); os.IsNotExist(err) {
      return fmt.Errorf("ssh_config file %s does not exist", options.SSHConfigFile)
    }
  }
  if options.SSHConfigFile != "" {
    appendEnvironmentVariable("ANSIBLE_SSH_ARGS",
        fmt.Sprintf("-F %s", options.SSHConfigFile))
  }
  return nil
}

// if the user specifies an environment we will attempt to ensure that
// it exists, we will also set the inventory arg to be passed onto ansible.
func configureEnvironment(options *Options) error {
  if options.Environment != "" {
    fi, err := os.Stat(options.Environment)
    switch {
      case err != nil:
        return fmt.Errorf("Environment %s does not exist", options.Environment)
      case fi.IsDir():
        options.Inventory = filepath.Join(options.Environment, "hosts")
      default:
        return fmt.Errorf("Environment %s should be a directory", options.Environment)
    }
  }
  if options.Inventory != "" {
    if _, err := os.Stat(options.Inventory); os.IsNotExist(err) {
      return fmt.Errorf("inventory file %s does not exist", options.Inventory)
    }
  }
  return nil
}

// checks if the user specifies an alternative known hosts file, if not
// looks for one in the specified environment or just defaults to none.
func configureKnownHostsFile(options *Options) error {
  if options.KnownHostsFile != "" {
    if _, err := os.Stat(options.KnownHostsFile); os.IsNotExist(err) {
      return fmt.Errorf("known hosts file %s does not exist", options.KnownHostsFile)
    }
  } else {
    maybeKnownHostsFile := filepath.Join(options.Environment, "ssh_known_hosts")
    if _, err := os.Stat(maybeKnownHostsFile);! os.IsNotExist(err) {
      options.KnownHostsFile = maybeKnownHostsFile
    }
  }
  if options.KnownHostsFile != "" {
    appendEnvironmentVariable("ANSIBLE_SSH_ARGS",
        fmt.Sprintf("-o UserKnownHostsFile=%s", options.KnownHostsFile))
  }
  return nil
}

// enables ssh agent forwarding
func configureSSHForwardAgent(options *Options) {
  if options.SSHForwardAgent {
    appendEnvironmentVariable("ANSIBLE_SSH_ARGS", "-o ForwardAgent=yes")
  }
}

// Run ansible playbook
func Run(options *Options, ansibleArgs []string) error {
  var (
		cmdOut      []byte
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

  if options.Inventory != "" {
    gosibleArgs = append(gosibleArgs,
      []string{"--inventory", options.Inventory}...)
  }

  cmdName := "ansible-playbook"
	cmdArgs := append(gosibleArgs, ansibleArgs...)
  fmt.Println("running: ansible_playbook", strings.Join(cmdArgs, " "))

  // TODO to switch to streaming output
  cmdOut, err = exec.Command(cmdName, cmdArgs...).CombinedOutput()
  fmt.Println(string(cmdOut))
  if err != nil {
		fmt.Fprintln(os.Stderr, "There was an error running Ansible: ", err)
		os.Exit(1)
	}
  return nil
}
