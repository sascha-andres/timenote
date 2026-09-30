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

	"github.com/jason0x43/go-toggl"
)

// dispatchProjects routes "projects" sub-verbs; with none given it lists
// projects.
func dispatchProjects(verbs []string) {
	if len(verbs) == 0 {
		projectsList()
		return
	}
	switch verbs[0] {
	case "create":
		projectsCreate()
	case "delete":
		projectsDelete()
	default:
		fatal(fmt.Errorf("unknown projects command %q", verbs[0]))
	}
}

func projectsList() {
	p := newPersistor()
	projects, err := p.Projects()
	if err != nil {
		fatal(err)
	}

	if *outputFormat != "json" {
		writeProjectsTable(filterProjects(projects))
	} else {
		writeProjectsJson(filterProjects(projects))
	}
}

func writeProjectsJson(projects []toggl.Project) {
	data, err := json.Marshal(projects)
	if err != nil {
		fatal(err)
	}
	fmt.Println(string(data))
}

func writeProjectsTable(projects []toggl.Project) {
	w := new(tabwriter.Writer)
	// Format in tab-separated columns with a tab stop of 8.
	w.Init(os.Stdout, 0, 8, 2, '\t', 0)
	fmt.Fprintln(w, "ID\tName\t")
	for _, prj := range projects {
		fmt.Fprintf(w, "%d\t%s\t\n", prj.ID, prj.Name)
	}
	_ = w.Flush()
}

// filterProjects removes projects named in the excluded-projects flag.
func filterProjects(projects []toggl.Project) []toggl.Project {
	excludeList := excludedProjects()
	filtered := make([]toggl.Project, 0)
	for _, prj := range projects {
		excluded := false
		for _, v := range excludeList {
			if v == prj.Name {
				excluded = true
				break
			}
		}
		if !excluded {
			filtered = append(filtered, prj)
		}
	}
	return filtered
}
