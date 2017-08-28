// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package playbook

import (
  	"fmt"
    "os"
  	"os/exec"
    "strings"
    "path/filepath"
)

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
func sshConfigFile(options *Options) {
  if options.SSHConfigFile == "" {
    sshConfigFile := filepath.Join(options.Environment, "ssh_config")
    if _, err := os.Stat(sshConfigFile); !os.IsNotExist(err) {
      options.SSHConfigFile = sshConfigFile
    }
  } else {
    if _, err := os.Stat(options.SSHConfigFile); os.IsNotExist(err) {
      fmt.Println("ssh_config file", options.SSHConfigFile, "does not exist")
      os.Exit(1)
    }
  }
  if options.SSHConfigFile != "" {
    appendEnvironmentVariable("ANSIBLE_SSH_ARGS",
        fmt.Sprintf("-F %s", options.SSHConfigFile))
  }
}

// if the user specifies an environment we will attempt to ensure that
// it exists, we will also set the inventory arg to be passed onto ansible.
func configureEnvironment(options *Options) {
  if options.Environment != "" {
    fi, err := os.Stat(options.Environment)
    switch {
      case err != nil:
        fmt.Printf("Environment %s does not exist", options.Environment)
        os.Exit(1)
      case fi.IsDir():
        options.Inventory = filepath.Join(options.Environment, "hosts")
      default:
        fmt.Println("Environment", options.Environment, "should be a directory")
        os.Exit(1)
    }
  }
  if options.Inventory != "" {
    if _, err := os.Stat(options.Inventory); os.IsNotExist(err) {
      fmt.Println("inventory file", options.Inventory, "does not exist")
      os.Exit(1)
    }
  }
}

// checks if the user specifies an alternative known hosts file, if not
// looks for one in the specified environment or just defaults to none.
func configureKnownHostsFile(options *Options) {
  if options.KnownHostsFile != "" {
    if _, err := os.Stat(options.KnownHostsFile); os.IsNotExist(err) {
      fmt.Println("known hosts file", options.KnownHostsFile, "does not exist")
      os.Exit(1)
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
}

// enables ssh agent forwarding
func configureSSHForwardAgent(options *Options) {
  if options.SSHForwardAgent {
    appendEnvironmentVariable("ANSIBLE_SSH_ARGS", "-o ForwardAgent=yes")
  }
}

func Run(options *Options, ansibleArgs []string) {
  var (
		cmdOut []byte
		err    error
	)

  configureEnvironment(options)
  sshConfigFile(options)
  configureKnownHostsFile(options)
  configureSSHForwardAgent(options)

  gosibleArgs := []string{"--inventory", options.Inventory}

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
	fmt.Println("Successfully ran ansible?")
}
