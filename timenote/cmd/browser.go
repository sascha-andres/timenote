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
	"github.com/pkg/browser"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"livingit.de/code/timenote/persistence"
)

// browserCmd represents the browser command
var browserCmd = &cobra.Command{
	Use:   "browser",
	Short: "open a browser with your overview",
	Long: `Depending on your backend this may open a browser
	with your dashboard.`,
	Run: func(cmd *cobra.Command, args []string) {
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

		has, url, err := p.GetWebsite()
		if err != nil {
			log.Error(err)
			return
		}
		if has {
			err = browser.OpenURL(url)
			if err != nil {
				log.Warnf("error executing browser: %s", err)
			}
		} else {
			fmt.Println("no url for backend")
		}
	},
}

func init() {
	RootCmd.AddCommand(browserCmd)
}
