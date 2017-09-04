package provisioner

import (
  "fmt"
  "os"
  "github.com/paulczar/vagrantutil"
  "path/filepath"
  "io/ioutil"
)

var vagrant *vagrantutil.Vagrant

func findVagrantfile(paths []string) (string) {
  var (
    vagrantFile string
    path        string
  )
  for _, path = range paths {
    vagrantFile = filepath.Join(path, "Vagrantfile")
    if _, err := os.Stat(vagrantFile); !os.IsNotExist(err) {
      return path
    }
  }
  return ""
}

func runVagrantUp() {
  output, _ := vagrant.Up()
  // print the output
  for line := range output {
    fmt.Println(line)
  }
}

// VagrantUp provisioner
func VagrantUp(po *Options) error {
  var (
    err error
  )
  cwd, _ := os.Getwd()
  paths := []string{po.Environment, cwd} 
  vagrantPath := findVagrantfile(paths)
  if vagrantPath == "" {
    return fmt.Errorf("cannot find vagrantfile in %s", paths)
  }
  vagrant, _ = vagrantutil.NewVagrant(vagrantPath)
  vagrant.Status()
	switch state := vagrant.State; state {
    case "NotCreated":
      runVagrantUp()
    case "PowerOff":
      runVagrantUp()
    case "Running":
      // do nothing
    default:
      return fmt.Errorf("cannot handle vagrant state %s.", state)
    }
      
  if vagrant.State == "NotCreated" {
    output, _ := vagrant.Up()
    // print the output
    for line := range output {
      fmt.Println(line)
    }    
  }
  sshConfig, err := (vagrant.SSHConfig())
  if err != nil {
    return err
  }
  sshConfigFile := filepath.Join(vagrantPath, "ssh_config")
  err = ioutil.WriteFile(sshConfigFile, []byte(sshConfig), 0644)
  if err != nil {
    return err
  }

  return nil
}