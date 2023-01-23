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
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"livingit.de/code/timenote/internal/persistence"
)

var timestampProjectCmd = &cobra.Command{
	Use:   "project",
	Short: "current entry is part of project",
	Long:  `Depending on the persistor, this command sets the project`,
	Run: func(cmd *cobra.Command, args []string) {
		p, err := persistence.NewToggl(token, viper.GetInt("workspace"), caching)
		if err != nil {
			log.Fatal(err)
		}

		err = p.SetProjectForCurrentTimestamp(viper.GetString("project.name"), viper.GetBool("project.auto-create-project"))
		if err != nil {
			log.Fatalf("error setting project: %s", err)
		}
	},
}

func init() {
	timestampCmd.AddCommand(timestampProjectCmd)

	timestampProjectCmd.Flags().StringP("name", "", "", "Name for project")
	timestampProjectCmd.Flags().BoolP("auto-create-project", "", false, "automatically create project if it does not exist")
	_ = timestampProjectCmd.MarkFlagRequired("name")

	_ = viper.BindPFlag("project.name", timestampProjectCmd.Flags().Lookup("name"))
	_ = viper.BindPFlag("project.auth-create-project", timestampProjectCmd.Flags().Lookup("auto-create-project"))
}
