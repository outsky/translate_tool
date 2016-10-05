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

var update_dbname string
var update_data string

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update translation to dictionary",
	Long:  `Update translation to dictionary`,
	Run: func(cmd *cobra.Command, args []string) {
		analysis.GetInstance().Update(path.Clean(update_dbname), path.Clean(update_data))
	},
}

func init() {
	RootCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringVarP(&update_dbname, "db", "d", "dictionary.txt", "Translation data dictionary")
	updateCmd.Flags().StringVarP(&update_data, "update", "u", "chinese.txt", "The new translation data")
}
