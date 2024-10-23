/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup a new react-native git project",
	Run:   setup,
}

func setup(cmd *cobra.Command, args []string) {
	fmt.Println("setup called")
	createFolderIfNotExist("~/.config/celeste")
	repoUrl := getRepoUrl()
	fmt.Println("repo: %s", repoUrl)
	// clone repo on .config/celeste using bare mode
}

func createFolderIfNotExist(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Create the directory
		err := os.Mkdir(path, 0755) // 0755 is the permission
		if err != nil {
			fmt.Println("Error creating directory:", err)
			return
		}
		fmt.Println("Directory created:", path)
	} else {
		fmt.Println("Directory already exists:", path)
	}
}

func getRepoUrl() string {
	// git remote -v | awk '/fetch/ {print $2}'
	gitCmd := exec.Command("git", "remote", "-v")

	pipe, err := gitCmd.StdoutPipe()
	if err != nil {
		log.Fatal("Error creating pipe:", err)
	}

	if err := gitCmd.Start(); err != nil {
		log.Fatal("Error starting git command:", err)
		return ""
	}

	awkCmd := exec.Command("awk", "/fetch/ {print $2}")
	awkCmd.Stdin = pipe

	awkOutput, err := awkCmd.Output()
	if err != nil {
		log.Fatal("Error running awk command:", err)
		return ""
	}

	if err := gitCmd.Wait(); err != nil {
		log.Fatal("Error waiting for git command:", err)
		return ""
	}

	if err := awkCmd.Wait(); err != nil {
		log.Fatal("Error waiting for awk command:", err)
		return ""
	}

	return strings.TrimSpace(string(awkOutput))
}

func init() {
	rootCmd.AddCommand(setupCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
