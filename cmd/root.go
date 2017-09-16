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
  "os"

  homedir "github.com/mitchellh/go-homedir"
  "github.com/spf13/cobra"
  "github.com/spf13/viper"
)

type Options struct {
  cfgFile     string
}

type RootOptions struct {
  Options
}

var ro = &RootOptions{}

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
  Use:   "gosible",
  Short: "Gosible is a wrapper around Ansible",
  Long: `
Gosible is a CLI tool designed to implement stronger 
Infrastructure-as-Code abilities to Ansible.

https://github.com/paulczar/gosible

`,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
  if err := RootCmd.Execute(); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}

func init() {
  cobra.OnInitialize(initConfig)

  RootCmd.PersistentFlags().StringVar(&ro.cfgFile, "config", "",
    "config file (default is $HOME/.gosible.yaml)")
}


// initConfig reads in config file and ENV variables if set.
func initConfig() {
  if ro.cfgFile != "" {
    // Use config file from the flag.
    viper.SetConfigFile(ro.cfgFile)
  } else {
    // Find home directory.
    home, err := homedir.Dir()
    if err != nil {
      fmt.Println(err)
      os.Exit(1)
    }

    // Search config in home directory with name ".gosible" (without extension).
    viper.AddConfigPath(home)
    viper.SetConfigName(".gosible")
  }

  viper.AutomaticEnv() // read in environment variables that match

  // If a config file is found, read it in.
  if err := viper.ReadInConfig(); err == nil {
    fmt.Println("Using config file:", viper.ConfigFileUsed())
  }
}
