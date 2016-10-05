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
	"trans/analysis"

	"github.com/spf13/cobra"
)

var translate_dbname string
var translate_update_data string
var translate_srcpath string
var translate_output string
var translate_routine int

// translateCmd represents the translate command
var translateCmd = &cobra.Command{
	Use:   "translate",
	Short: "Translation file or directory",
	Long:  `Translation using dictionary file or directory. If the output does not exist will be created automatically`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(translate_srcpath) == 0 || len(translate_output) == 0 {
			cmd.Help()
			return
		}
		analysis.GetInstance().Translate(
			path.Clean(translate_dbname),
			path.Clean(translate_update_data),
			path.Clean(translate_srcpath),
			path.Clean(translate_output),
			translate_routine)
	},
}

func init() {
	RootCmd.AddCommand(translateCmd)

	translateCmd.Flags().StringVarP(&translate_dbname, "db", "d", "dictionary.txt", "Translation data dictionary")
	translateCmd.Flags().StringVarP(&translate_update_data, "update", "u", "chinese.txt", "The new translation data")
	translateCmd.Flags().StringVarP(&translate_srcpath, "src", "s", "", "Translated file or directory path")
	translateCmd.Flags().StringVarP(&translate_output, "output", "o", "", "The output file or directory path translated")
	translateCmd.Flags().IntVarP(&translate_routine, "routine", "r", 1, "Goroutine number. This is a test parameters")
}
