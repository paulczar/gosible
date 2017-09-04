package ansible

import (
)

// Options for installing ansible
type InstallOptions struct {
	Version  string
  Method   string
  Sudo     bool
}

// Options for installing ansible via pip
type InstallViaPip struct {
  RequirementsTXT string
  VirtualEnv      string
  InstallOptions
}

func InstallAnsible() error {
  // TODO
  return nil
}

func installAnsibleFromPip() error {
 // TODO
  return nil 
}

func installAnsibleFromApt() error {
  // TODO
  return nil
}

func installAnsibleFromYum() error {
  // TODO
  return nil
}
