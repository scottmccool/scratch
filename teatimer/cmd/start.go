/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/kyokomi/emoji"

	"time"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a timer for <duration> (Xs)",
	Long: `Starts a timer for the specified duration
	Valid durations:
	h - hour
	m - minute
	s - second
	including "3m,30s"
	`,
	Run: func(cmd *cobra.Command, args []string) {
		sTime, err := time.ParseDuration(args[0])
		if err == nil {
			fmt.Println("Starting a timer for ", sTime)
			time.Sleep(sTime)
			emoji.Println(":heart: :coffee: Your tea is ready!  :coffee:  :heart: ")
			return
		} else {
			log.Println(err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
