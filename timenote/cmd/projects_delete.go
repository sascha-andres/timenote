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
	"livingit.de/code/timenote/persistence"
)

// timestampAppendCmd represents the append command
var projectsDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Add a new project to the project list",
	Long: `This will add a new project to the time backend

If project already exists it will not do anything`,
	Run: func(cmd *cobra.Command, args []string) {
		name := viper.GetString("projects.delete.name")
		p, err := persistence.NewToggl(viper.GetString("dsn"), viper.GetInt("workspace"))
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			err := p.Close()
			if err != nil {
				log.Fatal(err)
			}
		}()

		err = p.DeleteProject(name)
		if err != nil {
			log.Error(err)
		}
	},
}

func init() {
	projectsCmd.AddCommand(projectsDeleteCmd)

	projectsDeleteCmd.Flags().StringP("name", "", "", "name for project. id can also be used")
	_ = projectsDeleteCmd.MarkFlagRequired("name")

	_ = viper.BindPFlag("projects.delete.name", projectsDeleteCmd.Flags().Lookup("name"))
}
