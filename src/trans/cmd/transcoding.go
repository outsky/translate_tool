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
	"trans/filetool"

	"github.com/spf13/cobra"
)

var tc_decoding string
var tc_input string
var tc_encoding string
var tc_output string

// transcodingCmd represents the transcoding command
var transcodingCmd = &cobra.Command{
	Use:   "transcoding",
	Short: "file transcoding",
	Long: `file transcoding Support utf8, gbk, hz-gb2312, gb18030, big5, euc-jp, iso-2022-jp, shift_jis, euc-kr.
Notice: This tool can only transcoding, and can not be translated.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		if len(tc_decoding) <= 0 || len(tc_encoding) <= 0 || len(tc_input) <= 0 || len(tc_output) <= 0 {
			cmd.Help()
			return
		}
		filetool.GetInstance().Transcoding(tc_input, tc_decoding, tc_output, tc_encoding)
	},
}

func init() {
	RootCmd.AddCommand(transcodingCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// transcodingCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// transcodingCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	transcodingCmd.Flags().StringVarP(&tc_decoding, "decoding", "d", "", "encoding of input file")
	transcodingCmd.Flags().StringVarP(&tc_encoding, "encoding", "e", "", "encoding of output file")
	transcodingCmd.Flags().StringVarP(&tc_input, "input", "i", "", "Input file or directory")
	transcodingCmd.Flags().StringVarP(&tc_output, "output", "o", "", "Output file or directory")
}
