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

// dispatchCache routes "cache" sub-verbs: info, update.
func dispatchCache(verbs []string) {
	if len(verbs) == 0 {
		fatal(fmt.Errorf("cache requires a command: info, update"))
	}
	switch verbs[0] {
	case "info":
		cacheInfo()
	case "update":
		cacheUpdate()
	default:
		fatal(fmt.Errorf("unknown cache command %q", verbs[0]))
	}
}
