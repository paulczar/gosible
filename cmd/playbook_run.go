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

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/paulczar/gosible/ansible/playbook"
)

// runCmd represents the run command
var playbookRunCmd = &cobra.Command{
	Use:   "run [flags] file.yml [ansible-playbook arguments]",
	Short: "wrapper around ansible command",
	Long: `
Gosible playbook is a wrapper around ansible-playbook that adds some
additional useful features.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
 			return fmt.Errorf("must specify a playbook to run")
    }
		err := playbook.Run(ao,args)
		if err != nil {
			return err
		} else {
			fmt.Println("Successfully ran ansible?")
			return nil
		}
	},
}

func init() {
	playbookCmd.AddCommand(playbookRunCmd)
	// stops parsing flags after first unknown flag is found
	playbookRunCmd.Flags().SetInterspersed(false)
	playbookRunCmd.Flags().StringVarP(&ao.SSHConfigFile, "ssh-config-file",
		"s", "", "Path to ssh config file to use.")
	playbookRunCmd.Flags().BoolVarP(&ao.SSHForwardAgent, "ssh-forward-agent",
		"f", false, "path to ssh config file to use")
	playbookRunCmd.Flags().StringVarP(&ao.Environment, "environment",
		"e", "", "directory that contains ansible inventory")
	playbookRunCmd.Flags().StringVarP(&ao.KnownHostsFile, "known-hosts-file",
		"", "", "location of known hosts file")
}
