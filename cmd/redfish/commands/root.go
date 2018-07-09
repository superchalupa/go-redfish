// Copyright © 2018 Dell EMC
//
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

package commands

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cliCfgFile string
var cliCfg *viper.Viper

var rootCmd = &cobra.Command{
	Use:   "redfish",
	Short: "A redfish HTTP server",
	Long:  `Go-redfish is an HTTP server that implements the redfish protocol and is intended to be run in an embedded system.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cliCfgFile, "config", "", "config file (default is /etc/redfish.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	cliCfg = viper.New()

	if cliCfgFile != "" {
		// Use config file from the flag.
		cliCfg.SetConfigFile(cliCfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".redfish" (without extension).
		cliCfg.AddConfigPath("/etc/")
		cliCfg.AddConfigPath(home)
		cliCfg.SetConfigName("redfish")
	}

	cliCfg.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := cliCfg.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", cliCfg.ConfigFileUsed())
	}
}