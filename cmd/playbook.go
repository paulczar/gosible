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
	//"fmt"
  //"errors"
	"github.com/spf13/cobra"
	"github.com/paulczar/gosible/ansible/playbook"
)

var ao = &playbook.Options{}

// playbookCmd represents the playbook command
var playbookCmd = &cobra.Command{
	Use:   "playbook [flags] file.yml [ansible-playbook arguments]",
	Short: "wrapper around ansible-playbook command",
	Long: `
Gosible playbook is a wrapper around ansible-playbook that adds some
additional useful features.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		playbook.Run(ao,args)
		return nil
	},
}

func init() {
	RootCmd.AddCommand(playbookCmd)

	// stops parsing flags after first unknown flag is found
	playbookCmd.Flags().SetInterspersed(false)
	playbookCmd.Flags().StringVarP(&ao.SSHConfigFile, "ssh-config-file",
		"s", "", "Path to ssh config file to use.")
	playbookCmd.Flags().BoolVarP(&ao.SSHForwardAgent, "ssh-forward-agent",
		"f", false, "path to ssh config file to use")
	playbookCmd.Flags().StringVarP(&ao.Environment, "environment",
		"e", "", "directory that contains ansible inventory")
	playbookCmd.Flags().StringVarP(&ao.KnownHostsFile, "known-hosts-file",
		"", "", "location of known hosts file")
}
