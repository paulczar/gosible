package ansible

import (
)

// Options for installing ansible
type Options struct {
	Version  string
  Method   string
  Sudo     bool
}

// Options for installing ansible via pip
type Pip struct {
  RequirementsTXT string
  VirtualEnv      string
  Options
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
