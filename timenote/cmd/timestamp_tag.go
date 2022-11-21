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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"livingit.de/code/timenote/internal/persistence"
	"log"
)

// timestampAppendCmd represents the append command
var timestampTagCmd = &cobra.Command{
	Use:   "tag",
	Short: "tag current entry",
	Long:  `Depending on the persistor, this command adds or overwrites a tag`,
	Run: func(cmd *cobra.Command, args []string) {
		name := viper.GetString("tag.name")
		p, err := persistence.NewToggl(token, viper.GetInt("workspace"), caching)
		if err != nil {
			log.Fatal(err)
		}

		err = p.Tag(name)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	timestampCmd.AddCommand(timestampTagCmd)

	timestampTagCmd.Flags().StringP("name", "", "", "Name for tag")
	_ = timestampTagCmd.MarkFlagRequired("name")

	_ = viper.BindPFlag("tag.name", timestampTagCmd.Flags().Lookup("name"))
}
