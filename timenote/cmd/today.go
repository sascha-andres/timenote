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
	"os"
	"text/tabwriter"

	"go.livingit.de/timenote"
)

// timestampToday prints all timestamps from today, or those still active.
func timestampToday() {
	p := newPersistor()

	ts, err := p.ListForDay()
	if err != nil {
		fatal(err)
	}

	if *sumOnly {
		printTodaySum(ts)
		return
	}

	if *outputFormat == "json" {
		writeTimeEntriesJson(ts)
		return
	}

	if *group {
		ts = groupTimeEntries(ts)
	}
	writeTimeEntriesTable(ts)
}

func printTodaySum(ts []timenote.TimeEntry) {
	var sum int64
	for _, e := range ts {
		sum += e.Seconds()
	}
	td, err := timenote.NewTogglDuration(sum)
	if err != nil {
		panic(err)
	}
	if !*includeSeconds {
		td.OmitSeconds()
	}
	fmt.Println(td.String())
}

func groupTimeEntries(ts []timenote.TimeEntry) []timenote.TimeEntry {
	grouped := make(map[string]timenote.TimeEntry)
	result := make([]timenote.TimeEntry, 0)
	for _, e := range ts {
		if e.Duration <= 0 {
			e.Note = "[running] " + e.Note
			result = append(result, e)
			continue
		}
		if val, ok := grouped[e.Note]; ok {
			v := grouped[val.Note]
			v.Duration += val.Duration
			grouped[e.Note] = v
		} else {
			grouped[e.Note] = e
		}
	}
	for _, e := range grouped {
		result = append(result, e)
	}
	return result
}

func writeTimeEntriesJson(ts []timenote.TimeEntry) {
	data, err := json.Marshal(ts)
	if err != nil {
		fatal(err)
	}
	fmt.Println(string(data))
}

func writeTimeEntriesTable(ts []timenote.TimeEntry) {
	w := new(tabwriter.Writer)
	// Format in tab-separated columns with a tab stop of 8.
	w.Init(os.Stdout, 0, 8, 2, '\t', 0)
	fmt.Fprintln(w, "ID\tTime\tClient\tProject\tNote\t")
	var totalDuration int64
	for _, e := range ts {
		totalDuration += e.Seconds()
		td, _ := timenote.NewTogglDuration(e.Seconds())
		if !*includeSeconds {
			td.OmitSeconds()
		}
		humanTime := td.String()
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\t\n", e.ID, humanTime, e.Client, e.Project, e.Note)
	}
	td, _ := timenote.NewTogglDuration(totalDuration)
	fmt.Fprintln(w, "------------\t------------\t------------\t------------\t------------\t")
	fmt.Fprintf(w, "%s\t%s\t\t\t\t\n", "total", td.String())
	fmt.Fprintln(w)
	_ = w.Flush()
}
