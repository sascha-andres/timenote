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
	"livingit.de/code/timenote/persistence/factory"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// timestampDoneCmd represents the done command
var timestampDoneCmd = &cobra.Command{
	Use:   "done",
	Short: "done with current entry",
	Long:  `Stops the current entry`,
	Run: func(cmd *cobra.Command, args []string) {
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

		err = persistence.Done()
		if err != nil {
			log.Error(err)
			return
		}
	},
}

func init() {
	timestampCmd.AddCommand(timestampDoneCmd)
}
