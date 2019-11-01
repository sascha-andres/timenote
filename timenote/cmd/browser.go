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
	"github.com/pkg/browser"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// browserCmd represents the browser command
var browserCmd = &cobra.Command{
	Use:   "browser",
	Short: "open a browser with your overview",
	Long: `Depending on your backend this may open a browser
	with your dashboard.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := browser.OpenURL("https://toggl.com/app/timer")
		if err != nil {
			log.Warnf("error executing browser: %s", err)
		}
	},
}

func init() {
	RootCmd.AddCommand(browserCmd)
}
