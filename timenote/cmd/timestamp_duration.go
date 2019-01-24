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
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"livingit.de/code/timenote/persistence/factory"
	"time"
)

// timestampCurrentCmd represents the current command
var timestampDurationCmd = &cobra.Command{
	Use:   "duration",
	Short: "Print current timestamp duration",
	Long: `Prints the current timestamp's duration in
hh:mm:ss'`,
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

		ts, err := persistence.Current()
		if err != nil {
			log.Error(err)
			return
		}
		if ts.Duration != 0 {
			t := time.Now().UTC().Add(time.Duration(ts.Duration) * time.Second)
			fmt.Printf("%s %s - %s", ts.Tag, t.Format("15:04:05"), ts.Note)
		} else {
			fmt.Println("not supported for storage type")
		}
	},
}

func init() {
	timestampCmd.AddCommand(timestampDurationCmd)
}