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
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"livingit.de/code/timenote/persistence/factory"
)

// timestampNewCmd represents the new command
var timestampNewCmd = &cobra.Command{
	Use:   "new",
	Short: "add a new timestamp",
	Long:  `Starts a new timestamp`,
	Run: func(cmd *cobra.Command, args []string) {
		description := viper.GetString("timestamp.new.description")
		persistence, err := factory.CreatePersistence(viper.GetString("dsn"))
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			err := persistence.Close()
			if err != nil {
				log.Fatal(err)
			}
		}()

		if err := persistence.New(); err != nil {
			log.Fatal(err)
		} else {
			persistence.Append(description)
		}

	},
}

func init() {
	timestampCmd.AddCommand(timestampNewCmd)
	timestampNewCmd.Flags().StringP("description", "", "", "Description for timestamp")
	timestampNewCmd.MarkFlagRequired("description")
	viper.BindPFlag("timestamp.new.description", timestampNewCmd.Flags().Lookup("description"))
}
