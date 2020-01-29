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
	"context"
	"fmt"
	"log"

	"github.com/google/go-github/github"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List Available Helm Versions to Download",
	Long:  "List Available Helm Versions to Download",
	Run: func(cmd *cobra.Command, args []string) {
		client := github.NewClient(nil)
		ctx := context.Background()
		opt := github.ListOptions{
			Page:    1,
			PerPage: 50,
		}
		tags, _, err := client.Repositories.ListReleases(ctx, "helm", "helm", &opt)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Last 50 Helm Releases:")
		for _, element := range tags {
			fmt.Println(element.GetTagName())
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
