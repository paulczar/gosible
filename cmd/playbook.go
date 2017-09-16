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
  "github.com/spf13/cobra"
  "github.com/paulczar/gosible/ansible"
  "github.com/paulczar/gosible/provisioner"
)

var playbookOptions = &ansible.Options{}
var po = &provisioner.Options{}

// playbookCmd represents the playbook command
var playbookCmd = &cobra.Command{
  Use:   "playbook",
  Short: "wrapper around ansible command",
  Long: `
Gosible playbook is a wrapper around ansible-playbook that adds some
additional useful features.
`,

}

func init() {
  RootCmd.AddCommand(playbookCmd)
}
