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
)

// appendCmd represents the append command
var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "current entry is part of project",
	Long:  `Depending on the persistor, this command sets the project`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("todo")
	},
}

func init() {
	RootCmd.AddCommand(projectCmd)

	projectCmd.Flags().StringP("name", "", "", "Name for project")
	_ = projectCmd.MarkFlagRequired("name")

	_ = viper.BindPFlag("project.name", projectCmd.Flags().Lookup("name"))
}
