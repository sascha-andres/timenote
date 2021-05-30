// Copyright © 2021 Sascha Andres <sascha.andres@outlook.com>
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
	"github.com/spf13/cobra"
	"os"
)

var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "shell completion",
	Long: `To load completions:

Bash:

  $ source <(yourprogram completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ go-logsink completion bash > /etc/bash_completion.d/go-logsink
  # macOS:
  $ go-logsink completion bash > /usr/local/etc/bash_completion.d/go-logsink

Zsh:

  # If shell completion is not already enabled in your environment,
  # you will need to enable it.  You can execute the following once:

  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  $ go-logsink completion zsh > "${fpath[1]}/_go-logsink"

  # You will need to start a new shell for this setup to take effect.

fish:

  $ go-logsink completion fish | source

  # To load completions for each session, execute once:
  $ go-logsink completion fish > ~/.config/fish/completions/go-logsink.fish

PowerShell:

  PS> go-logsink completion powershell | Out-String | Invoke-Expression

  # To load completions for every new session, run:
  PS> go-logsink completion powershell > go-logsink.ps1
  # and source this file from your PowerShell profile.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Fprintf(os.Stderr, "wrong number of arguments")
			os.Exit(1)
		}
		switch args[0] {
		case "bash":
			cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			cmd.Root().GenFishCompletion(os.Stdout, true)
		case "powershell":
			cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
		}
	},
}

func init() {
	RootCmd.AddCommand(completionCmd)
}
