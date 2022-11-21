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
	"livingit.de/code/timenote/internal/persistence"
)

// timestampAppendCmd represents the append command
var timestampAppendCmd = &cobra.Command{
	Use:   "append",
	Short: "append/overwrite description for running timeentry",
	Long: `Depending on the persistor, this command appends
to the description or sets the description`,
	Run: func(cmd *cobra.Command, args []string) {
		description := viper.GetString("append.description")
		separator := viper.GetString("separator")
		p, err := persistence.NewToggl(token, viper.GetInt("workspace"), caching)
		if err != nil {
			log.Fatal(err)
		}

		err = p.Append(description, separator)
		if err != nil {
			log.Error(err)
		}
	},
}

func init() {
	timestampCmd.AddCommand(timestampAppendCmd)

	timestampAppendCmd.Flags().StringP("description", "", "", "Description for timestamp")
	_ = timestampAppendCmd.MarkFlagRequired("description")

	_ = viper.BindPFlag("append.description", timestampAppendCmd.Flags().Lookup("description"))
}
