// Copyright © 2021 Sascha Andres <sascha.andres@outlook.com>
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
)

// timestampCmd represents the timestamp command
var interactiveCmd = &cobra.Command{
	Use:   "i",
	Short: "interactive shell",
	Long:  `Interact with toggl interactively with a "shell"`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := run(); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(interactiveCmd)
}
