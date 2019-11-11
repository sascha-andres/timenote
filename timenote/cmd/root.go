// Copyright © 2017 Sascha Andres <sascha.andres@outlook.com>
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
	"livingit.de/code/timenote/internal/cache"
	"os"
	"path"

	log "github.com/sirupsen/logrus"

	"github.com/google/gops/agent"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var caching *cache.Cache

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "timenote",
	Short: "Take notes with attached timestamps",
	Long: `A timestamp will be attached when you start a not and a second
one as soon as you stop working on that note

You can tag notes`,
	Run: func(cmd *cobra.Command, args []string) {

		if err := agent.Listen(agent.Options{}); err != nil {
			log.Fatal(err)
		}

		if err := run(); err != nil {
			log.Fatal(err)
		}
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	c, err := cache.NewCache(viper.GetInt("cache.max-age"), viper.GetString("cache.path"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	caching = c
	defer func() {
		err := caching.Close()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}()
	if err = RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(home)
		os.Exit(1)
	}

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.timenote.yaml)")
	RootCmd.PersistentFlags().StringP("dsn", "d", "toggl-token", "Token to access Toggl API")
	RootCmd.PersistentFlags().IntP("workspace", "w", 0, "Set to work within this workspace, leave to zero to have it guessed (first workspace)")
	RootCmd.PersistentFlags().StringP("output-format", "", "text", "test or json")
	RootCmd.PersistentFlags().StringP("separator", "", ";", "Separator for existing value and new value")
	RootCmd.PersistentFlags().IntP("cache-max-age", "", 360, "Maximum age of cache in minutes")
	RootCmd.PersistentFlags().StringP("cache-path", "", path.Join(home, ".timenote"), "Where to store cache")

	_ = viper.BindPFlag("separator", RootCmd.PersistentFlags().Lookup("separator"))
	_ = viper.BindPFlag("dsn", RootCmd.PersistentFlags().Lookup("dsn"))
	_ = viper.BindPFlag("output-format", RootCmd.PersistentFlags().Lookup("output-format"))

	_ = viper.BindPFlag("cache.max-age", RootCmd.PersistentFlags().Lookup("cache-max-age"))
	_ = viper.BindPFlag("cache.path", RootCmd.PersistentFlags().Lookup("cache-path"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(home)
			os.Exit(1)
		}

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(path.Join(home, ".config/timenote"))
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigName(".timenote")
	}

	viper.AutomaticEnv() // read in environment variables that match
	_ = viper.ReadInConfig()
}
