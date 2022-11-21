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
	"encoding/json"
	"fmt"
	"livingit.de/code/timenote"
	"livingit.de/code/timenote/internal/persistence"
	"os"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

// timestampCurrentCmd represents the current command
var timestampTodayCmd = &cobra.Command{
	Use:   "today",
	Short: "Print timestamps from today",
	Long:  `Print all timestamps with a date from today or being active`,
	Run: func(cmd *cobra.Command, args []string) {
		p, err := persistence.NewToggl(token, viper.GetInt("workspace"), caching)
		if err != nil {
			log.Fatal(err)
		}

		ts, err := p.ListForDay()
		if err != nil {
			log.Fatal(err)
			return
		}
		if !viper.GetBool("timestamp.today.sum-only") {
			if viper.GetString("output-format") == "json" {
				writeTimeEntriesJson(ts)
				return
			}
			if viper.GetBool("timestamp.today.group") {
				ts2 := make(map[string]timenote.TimeEntry, 0)
				ts3 := make([]timenote.TimeEntry, 0)
				for _, e := range ts {
					if e.Duration <= 0 {
						e.Note = "[running] " + e.Note
						ts3 = append(ts3, e)
						continue
					}
					if val, ok := ts2[e.Note]; ok {
						v := ts2[val.Note]
						v.Duration += val.Duration
						ts2[e.Note] = v
					} else {
						ts2[e.Note] = e
					}
				}
				for _, e := range ts2 {
					ts3 = append(ts3, e)
				}
				ts = ts3
			}
			writeTimeEntriesTable(ts)
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
			if !viper.GetBool("timestamp.today.include-seconds") {
				td.OmitSeconds()
			}
			fmt.Println(td.String())
		}
	},
}

func writeTimeEntriesJson(ts []timenote.TimeEntry) {
	data, err := json.Marshal(ts)
	if err != nil {
		log.Fatal(err)
	}
	_, _ = fmt.Println(string(data))
}

func writeTimeEntriesTable(ts []timenote.TimeEntry) {
	w := new(tabwriter.Writer)
	// Format in tab-separated columns with a tab stop of 8.
	w.Init(os.Stdout, 0, 8, 2, '\t', 0)
	_, _ = fmt.Fprintln(w, "ID\tTime\tClient\tProject\tNote\t")
	for _, e := range ts {
		humanTime := ""
		if e.Duration >= 0 {
			td, _ := timenote.NewTogglDuration(e.Duration)
			if !viper.GetBool("timestamp.today.include-seconds") {
				td.OmitSeconds()
			}
			humanTime = td.String()
		} else {
			t := time.Now().UTC().Add(time.Duration(e.Duration) * time.Second)
			td2, _ := timenote.TogglDurationFromTime(t)
			if !viper.GetBool("timestamp.today.include-seconds") {
				td2.OmitSeconds()
			}
			humanTime = td2.String()
		}
		_, _ = fmt.Fprintln(w, fmt.Sprintf("%d\t%s\t%s\t%s\t%s\t", e.ID, humanTime, e.Client, e.Project, e.Note))
	}
	_, _ = fmt.Fprintln(w)
	_ = w.Flush()
}

func init() {
	RootCmd.AddCommand(timestampTodayCmd)

	timestampTodayCmd.Flags().BoolP("sum-only", "", false, "Just print sum of timestamps")
	timestampTodayCmd.Flags().BoolP("group", "", false, "Print grouped by name")
	timestampTodayCmd.Flags().BoolP("include-seconds", "", true, "Include seconds when writing out time entry")

	_ = viper.BindPFlag("timestamp.today.include-seconds", timestampTodayCmd.Flags().Lookup("include-seconds"))
	_ = viper.BindPFlag("timestamp.today.sum-only", timestampTodayCmd.Flags().Lookup("sum-only"))
	_ = viper.BindPFlag("timestamp.today.group", timestampTodayCmd.Flags().Lookup("group"))
}
