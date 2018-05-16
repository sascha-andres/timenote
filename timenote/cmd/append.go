// Copyright © 2018 Sascha Andres <sascha.andres@outlook.com>
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

package cmd

import (
	"github.com/sascha-andres/timenote/persistence/factory"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// appendCmd represents the append command
var appendCmd = &cobra.Command{
	Use:   "append",
	Short: "append/overwrite description for running timeentry",
	Long: `Depending on the persistor, this command appends
to the description or sets the description`,
	Run: func(cmd *cobra.Command, args []string) {
		description := viper.GetString("append.description")
		persistence, err := factory.CreatePersistence(viper.GetString("persistor"), viper.GetString("dsn"))
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			err := persistence.Close()
			if err != nil {
				log.Fatal(err)
			}
		}()

		err = persistence.Append(description)
		if err != nil {
			log.Error(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(appendCmd)

	appendCmd.Flags().StringP("description", "", "", "Description for timestamp")
	appendCmd.MarkFlagRequired("description")

	viper.BindPFlag("append.description", appendCmd.Flags().Lookup("description"))
}