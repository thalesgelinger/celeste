/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
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

	usr, err := user.Current()
	if err != nil {
		fmt.Println("Error getting current user:", err)
		return
	}

	userPath := filepath.Join(usr.HomeDir, ".local", "share", "celeste")
	fmt.Println("setup called")
	err = createFolderIfNotExist(userPath)
	if err != nil {
		log.Fatal("Error creating celeste folder", err.Error())
	}
	repo, err := getRepoInfo()
	if err != nil {
		log.Fatal("Error getting repo info", err.Error())
	}
	err = cloneBareRepo(repo.Url, repo.Name, userPath)
	if err != nil {
		log.Fatal("Error cloning repo", err.Error())
	}
}

func cloneBareRepo(repoUrl, repoName, targetDir string) error {
	cmd := exec.Command("git", "clone", "--bare", repoUrl, fmt.Sprintf("%s/%s", targetDir, repoName))

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error cloning repository:", err.Error())
		fmt.Println("Output:", string(output))
		return err
	}
	fmt.Println("Repository cloned successfully into:", targetDir)
	return nil
}

func createFolderIfNotExist(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, 0755)
		if err != nil {
			return err
		}
		fmt.Println("Directory created:", path)
	}
	return nil
}

type Repo struct {
	Url  string
	Name string
}

func getRepoInfo() (Repo, error) {
	gitCmd := exec.Command("git", "remote", "-v")

	pipe, err := gitCmd.StdoutPipe()
	if err != nil {
		log.Fatal("Error creating pipe:", err)
	}

	if err := gitCmd.Start(); err != nil {
		log.Fatal("Error starting git command:", err)
		return Repo{}, err
	}

	awkCmd := exec.Command("awk", "/fetch/ {print $2}")
	awkCmd.Stdin = pipe

	awkOutput, err := awkCmd.Output()
	if err != nil {
		log.Fatal("Error running awk command:", err)
		return Repo{}, err
	}

	if err := gitCmd.Wait(); err != nil {
		log.Fatal("Error waiting for git command:", err)
		return Repo{}, err
	}

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return Repo{}, err
	}

	// Extract the last element of the path
	dirName := filepath.Base(cwd)

	return Repo{
		Url:  strings.TrimSpace(string(awkOutput)),
		Name: dirName,
	}, nil

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
