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
	"github.com/sascha-andres/go-toggl"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"livingit.de/code/timenote/internal/persistence"
	"os"
	"text/tabwriter"
)

// timestampCmd represents the timestamp command
var projectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "projects management",
	Long:  `List projects and manage projects using sub commands`,
	Run: func(cmd *cobra.Command, args []string) {
		p, err := persistence.NewToggl(viper.GetString("dsn"), viper.GetInt("workspace"), caching)
		if err != nil {
			log.Fatal(err)
		}

		projects, err := p.Projects()
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		if viper.GetString("output-format") != "json" {
			writeProjectsTable(filterProjects(projects))
		} else {
			writeProjectsJson(filterProjects(projects))
		}
	},
}

func writeProjectsJson(projects []toggl.Project) {
	data, err := json.Marshal(projects)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	_, _ = fmt.Println(string(data))
}

func writeProjectsTable(projects []toggl.Project) {
	w := new(tabwriter.Writer)
	// Format in tab-separated columns with a tab stop of 8.
	w.Init(os.Stdout, 0, 8, 2, '\t', 0)
	_, _ = fmt.Fprintln(w, "ID\tName\t")
	for _, prj := range projects {
		_, _ = fmt.Fprintln(w, fmt.Sprintf("%d\t%s\t", prj.ID, prj.Name))
	}
	_ = w.Flush()
}

// isProjectInExcludeList returns true if project is filtered by exclude
func filterProjects(projects []toggl.Project) []toggl.Project {
	excludeList := viper.GetStringSlice("excluded-projects")
	filteredProjectList := make([]toggl.Project, 0)
	for _, prj := range projects {
		c := false
		for _, v := range excludeList {
			if v == prj.Name {
				c = true
				break
			}
		}
		if c {
			continue
		}
		filteredProjectList = append(filteredProjectList, prj)
	}
	return filteredProjectList
}

func init() {
	RootCmd.AddCommand(projectsCmd)
}
