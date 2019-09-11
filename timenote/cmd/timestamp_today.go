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
	"livingit.de/code/timenote"
	"os"
	"text/tabwriter"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"livingit.de/code/timenote/persistence/factory"
)

// timestampCurrentCmd represents the current command
var timestampTodayCmd = &cobra.Command{
	Use:   "today",
	Short: "Print timestamps from today",
	Long:  `Print all timestamps with a date from today or being active`,
	Run: func(cmd *cobra.Command, args []string) {
		persistence, err := factory.CreatePersistence(viper.GetString("persistor"), viper.GetString("dsn"))
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			err := persistence.Close()
			if err != nil {
				log.Fatal(err)
			}
		}()

		ts, err := persistence.ListForDay()
		if err != nil {
			log.Error(err)
			return
		}
		if !viper.GetBool("timestamp.today.sum-only") {
			w := new(tabwriter.Writer)
			// Format in tab-separated columns with a tab stop of 8.
			w.Init(os.Stdout, 0, 8, 2, '\t', 0)
			_, _ = fmt.Fprintln(w, "ID\tTime\tNote\t")
			for _, e := range ts {
				humanTime := ""
				if e.Duration >= 0 {
					td, _ := timenote.NewTogglDuration(e.Duration)
					humanTime = td.String()
				} else {
					t := time.Now().UTC().Add(time.Duration(e.Duration) * time.Second)
					td2, _ := timenote.TogglDurationFromTime(t)
					humanTime = td2.String()
				}
				_, _ = fmt.Fprintln(w, fmt.Sprintf("%d\t%s\t%s\t", e.ID, humanTime, e.Note))
			}
			_, _ = fmt.Fprintln(w)
			_ = w.Flush()
		} else {
			var sum int64
			for _, e := range ts {
				if e.Duration >= 0 {
					sum += e.Duration
				} else {
					t := time.Now().UTC().Add(time.Duration(e.Duration) * time.Second)
					td2, _ := timenote.TogglDurationFromTime(t)
					sum += td2.GetDuration()
				}
			}
			td, err := timenote.NewTogglDuration(sum)
			if err != nil {
				panic(err)
			}
			fmt.Println(td.String())
		}
	},
}

func init() {
	timestampCmd.AddCommand(timestampTodayCmd)

	timestampTodayCmd.Flags().BoolP("sum-only", "", false, "Just print sum of timestamps")

	_ = viper.BindPFlag("timestamp.today.sum-only", timestampTodayCmd.Flags().Lookup("sum-only"))
}
