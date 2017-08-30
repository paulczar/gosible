package provisioner

import (
	"fmt"
)

// Provisioner settings
type Options struct {
	Provisioner  string
	Environment	 string
}

// Up instantiates an environment
func Up(po *Options ) error {
	if po.Provisioner == "vagrant" {
		return VagrantUp(po)
	}
	return fmt.Errorf("%s is not a valid provisioner choice", po.Provisioner)
}