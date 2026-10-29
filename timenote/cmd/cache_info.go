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
	"context"
	"fmt"
	"strings"
)

// cacheInfo prints when the cache was last updated and when it will be
// updated next.
func cacheInfo() {
	status, err := newCachingBackend().CacheStatus(context.Background())
	if err != nil {
		fatal(err)
	}

	for i, s := range status {
		if i > 0 {
			fmt.Println()
		}
		title := s.Name + " cache"
		fmt.Println(title)
		fmt.Println(strings.Repeat("-", len(title)))
		fmt.Println()
		fmt.Printf("Updated:     %s\n", s.Updated.String())
		fmt.Printf("Next update: %s\n", s.NextUpdate.String())
	}
}
