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

// timestampAppendCmd represents the append command
var cacheInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "see information about cache",
	Long: `See some information about the cache. When it was updated
and when it will be updated`,
	Run: func(cmd *cobra.Command, args []string) {
		p, err := persistence.NewToggl(token, viper.GetInt("workspace"), caching)
		if err != nil {
			log.Fatal(err)
		}

		mdProjects, err := caching.ProjectMetaData(p.Workspace())
		if err != nil {
			log.Fatal(err)
		}
		mdClients, err := caching.ClientMetaData(p.Workspace())
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Project cache")
		fmt.Println("-------------")
		fmt.Println()
		fmt.Printf("Updated:     %s\n", mdProjects.Updated.String())
		fmt.Printf("Next update: %s\n", mdProjects.NextUpdate.String())
		fmt.Println()
		fmt.Println("Client cache")
		fmt.Println("------------")
		fmt.Println()
		fmt.Printf("Updated:     %s\n", mdClients.Updated.String())
		fmt.Printf("Next update: %s\n", mdClients.NextUpdate.String())
	},
}

func init() {
	cacheCmd.AddCommand(cacheInfoCmd)
}
