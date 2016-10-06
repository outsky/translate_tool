// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
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
	"path"

	"github.com/spf13/cobra"

	"trans/analysis"
)

var dbname string
var updateData string
var srcPath string
var output string
var logpath string
var routine int

// cmd represents the translate command
var cmd = &cobra.Command{
	Use:   "translate",
	Short: "Translation file or directory",
	Long:  `Translation using dictionary file or directory. If the output does not exist will be created automatically`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(srcPath) == 0 || len(output) == 0 {
			cmd.Help()
			return
		}
		analysis.GetInstance().Translate(
			path.Clean(dbname),
			path.Clean(updateData),
			path.Clean(srcPath),
			path.Clean(output),
			routine,
			path.Clean(logpath))
	},
}

func init() {
	RootCmd.AddCommand(cmd)

	cmd.Flags().StringVarP(&logpath, "log", "l", "", "Log path")
	cmd.Flags().StringVarP(&dbname, "db", "d", "dictionary.txt", "Translation data dictionary")
	cmd.Flags().StringVarP(&updateData, "update", "u", "chinese.txt", "The new translation data")
	cmd.Flags().StringVarP(&srcPath, "src", "s", "", "Translated file or directory path")
	cmd.Flags().StringVarP(&output, "output", "o", "", "The output file or directory path translated")
	cmd.Flags().IntVarP(&routine, "routine", "r", 1, "Goroutine number. This is a test parameters")
}
