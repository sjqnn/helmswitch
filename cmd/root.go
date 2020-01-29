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
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/spf13/cobra"

	"runtime"

	"context"

	"github.com/google/go-github/github"
	"github.com/manifoldco/promptui"
	"github.com/mholt/archiver/v3"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

// Check if Helm is installed. Display Current Version.
var installedHelmLocationString = "/usr/local/bin/helm"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "helmswitch",
	Short: "Switch your version of Helm",
	Long:  "Switch your version of Helm",
	// Check that you provided 1 argument that is a version.
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		selectedVersion := ""
		if len(args) < 1 {
			// Get Release Versions
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

			var helmVersions []string
			for _, element := range tags {
				helmVersions = append(helmVersions, element.GetTagName())
			}

			prompt := promptui.Select{
				Label: "Select Version",
				Items: helmVersions,
			}

			_, result, err := prompt.Run()

			if err != nil {
				fmt.Printf("Prompt failed %v\n", err)
				return
			}

			selectedVersion = result
			fmt.Println(selectedVersion)

			fmt.Printf("You choose %q\n", result)
		} else {
			selectedVersion = args[0]
		}
		installedHelmClientVersion, err := exec.Command("helm", "version", "--client", "--short").Output()
		if err != nil {
			fmt.Println("Failed getting Helm Client Version.")
			fmt.Println(err)
		}
		if len(installedHelmClientVersion) > 0 {
			fmt.Printf("Local Helm Client Version:\n%s", installedHelmClientVersion)
		}

		installedHelmServerVersion, err := exec.Command("helm", "version", "--server", "--short").Output()
		if err != nil {
			fmt.Println("Failed getting Helm Server Version.")
		}

		if len(installedHelmServerVersion) > 0 {
			fmt.Println(installedHelmServerVersion)
		}

		// Check if desired version of Helm exists.
		resp, err := http.Get("https://github.com/helm/helm/releases/tag/" + selectedVersion)

		if resp.StatusCode != 200 {
			fmt.Println("Could Not find Helm Release " + selectedVersion + " !")
			os.Exit(1)
		}
		if resp.StatusCode == 200 {
			fmt.Println("Found Helm Release " + selectedVersion + " !")
		}

		// Download and unzip desired version of Helm.
		tmpPath := "/tmp/helmswitch/"
		if _, err := os.Stat(tmpPath); os.IsNotExist(err) {
			os.Mkdir(tmpPath, 0777)
			os.Mkdir(tmpPath+"/"+selectedVersion+"/", 0777)
		}
		// Check OS
		if runtime.GOOS == "windows" {
			fmt.Println("Windows is not supported.")
			os.Exit(1)
		}
		if runtime.GOOS == "darwin" {
			fmt.Println("Mac OS detected")
			downloadLink := "https://get.helm.sh/helm-" + selectedVersion + "-darwin-amd64.tar.gz"
			if err := DownloadFile(tmpPath+selectedVersion+"-darwin-amd64.tar.gz", downloadLink); err != nil {
				panic(err)
			}
			err = archiver.Unarchive(tmpPath+selectedVersion+"-darwin-amd64.tar.gz", tmpPath+"/"+selectedVersion+"/")
			_, err := copy(tmpPath+"/"+selectedVersion+"/darwin-amd64/helm", installedHelmLocationString)
			err = os.Chmod(installedHelmLocationString, 0777)
			if err != nil {
				log.Println(err)
				os.Exit(1)
			}
			fmt.Println("Helm Release " + selectedVersion + " Installed!")
		}
		if runtime.GOOS == "linux" {
			fmt.Println("Linux OS detected")
			downloadLink := "https://get.helm.sh/helm-" + selectedVersion + "-linux-amd64.tar.gz"
			if err := DownloadFile(tmpPath+selectedVersion+"-linux-amd64.tar.gz", downloadLink); err != nil {
				panic(err)
			}
			err = archiver.Unarchive(tmpPath+selectedVersion+"-linux-amd64.tar.gz", tmpPath+"/"+selectedVersion+"/")
			_, err := copy(tmpPath+"/"+selectedVersion+"/linux-amd64/helm", installedHelmLocationString)
			err = os.Chmod(installedHelmLocationString, 0777)
			if err != nil {
				log.Println(err)
				os.Exit(1)
			}
			fmt.Println("Helm Release " + selectedVersion + " Installed!")
		}

		// Move Downloaded version of Helm to proper path.
		// Display new Helm version.
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.helmswitch.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".helmswitch" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".helmswitch")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

// DownloadFile downloads file from url to local.
func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}
