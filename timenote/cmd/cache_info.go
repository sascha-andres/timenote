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

import "fmt"

// cacheInfo prints when the cache was last updated and when it will be
// updated next.
func cacheInfo() {
	p := newPersistor()

	mdProjects, err := caching.ProjectMetaData(p.Workspace())
	if err != nil {
		fatal(err)
	}
	mdClients, err := caching.ClientMetaData(p.Workspace())
	if err != nil {
		fatal(err)
	}

	fmt.Println("Project cache")
	fmt.Println("-------------")
	fmt.Println()
	fmt.Printf("Updated:     %s\n", mdProjects.Updated.String())
	fmt.Printf("Next update: %s\n", mdProjects.NextUpdate.String())
	fmt.Println()
	fmt.Println("Client cache")
	fmt.Println("------------")
	fmt.Println()
	fmt.Printf("Updated:     %s\n", mdClients.Updated.String())
	fmt.Printf("Next update: %s\n", mdClients.NextUpdate.String())
}
