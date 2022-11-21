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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"livingit.de/code/timenote/internal/persistence"
	"log"
)

// timestampCmd represents the timestamp command
var timestampCmd = &cobra.Command{
	Use:   "timestamp",
	Short: "timestamp management",
	Long:  `Manage timestamps by adding`,
	Run: func(cmd *cobra.Command, args []string) {
		p, err := persistence.NewToggl(token, viper.GetInt("workspace"), caching)
		if err != nil {
			log.Fatal(err)
		}

		ts, err := p.Current()
		if err != nil {
			log.Printf("Error reading timestamp: %s", err)
			return
		}
		fmt.Println(ts)
	},
}

func init() {
	RootCmd.AddCommand(timestampCmd)
}
