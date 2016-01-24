/*
 * nomad-syncer
 * Copyright (c) 2016 Yieldbot, Inc. (http://github.com/yieldbot/nomad-syncer)
 * For the full copyright and license information, please view the LICENSE.txt file.
 */

// Package app provides the app information
package app

import (
	"encoding/json"
	"errors"
	"flag"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/yieldbot/gocli"
	"github.com/yieldbot/nomad-syncer/client"
)

var (
	cli             gocli.Cli
	nomadURL        string
	proxyURL        string
	nomadClient     client.Client
	usageFlag       bool
	versionFlag     bool
	versionExtFlag  bool
	prettyPrintFlag bool
	nomadFlag       string
	proxyFlag       string
	dockerPullFlag  bool
)

func init() {
	flag.BoolVar(&usageFlag, "h", false, "Display usage")
	flag.BoolVar(&usageFlag, "help", false, "Display usage")
	flag.BoolVar(&versionFlag, "version", false, "Display version information")
	flag.BoolVar(&versionFlag, "v", false, "Display version information")
	flag.BoolVar(&versionExtFlag, "vv", false, "Display extended version information")
	flag.BoolVar(&prettyPrintFlag, "pp", false, "Pretty print for JSON output")
	flag.StringVar(&nomadFlag, "nomad", "", "Nomad url (default \"http://localhost:4646\")")
	flag.StringVar(&proxyFlag, "proxy", "", "Proxy url")
	flag.BoolVar(&dockerPullFlag, "docker-pull", false, "Pull Docker images before sync")
}

// Run runs the app
func Run() {

	// Init cli
	cli = gocli.Cli{
		AppName:    "nomad-syncer",
		AppVersion: "1.1.0",
		AppDesc:    "An opinionated CLI for Nomad",
		CommandList: map[string]string{
			"jobs": "Retrieve jobs",
			"add":  "Add a job",
			"get":  "Get a job information",
			"del":  "Delete a job",
			"sync": "Sync jobs via a file or directory",
		},
	}
	cli.Init()

	// Run the app

	// Command
	if cli.Command != "" {

		// Init the Nomad client
		if nomadFlag != "" {
			nomadURL = nomadFlag
		} else if os.Getenv("NOMAD_URL") != "" {
			nomadURL = os.Getenv("NOMAD_URL")
		} else {
			nomadURL = "http://localhost:4646"
		}

		if proxyFlag != "" {
			proxyURL = proxyFlag
		} else if os.Getenv("NOMAD_SYNCER_PROXY_URL") != "" {
			proxyURL = os.Getenv("NOMAD_SYNCER_PROXY_URL")
		}

		if proxyURL != "" {
			p, err := url.Parse(proxyURL)
			if err != nil {
				cli.LogErr.Fatal("invalid proxy value due to " + err.Error())
			}
			if p != nil {

			}
			nomadClient = client.Client{URL: nomadURL, ProxyURL: p}
		} else {
			nomadClient = client.Client{URL: nomadURL}
		}

		// Run the command
		if cli.Command == "jobs" {
			// Get the jobs
			runJobsCmd()
		} else if cli.Command == "add" {
			// Add a job
			runAddCmd()
		} else if cli.Command == "get" {
			// Get a job information
			runGetCmd()
		} else if cli.Command == "del" {
			// Delete a job
			runDelCmd()
		} else if cli.Command == "sync" {
			// Sync jobs
			runSyncCmd()
		}
	} else if versionFlag || versionExtFlag {
		// Version
		cli.PrintVersion(versionExtFlag)
	} else {
		// Default
		cli.PrintUsage()
	}
}

// runAddCmd runs the add command
func runAddCmd() {
	// Get the job name
	var jobj string
	if len(cli.CommandArgs) > 0 {
		jobj = cli.CommandArgs[0]
	}

	// Add the job
	if err := nomadClient.AddJob(jobj); err != nil {
		cli.LogErr.Fatal(err)
	} else {
		cli.LogOut.Printf("job is added\n")
	}
}

// runJobsCmd runs the jobs command
func runJobsCmd() {
	if err := nomadClient.PrintJobs(prettyPrintFlag); err != nil {
		cli.LogErr.Fatal(err)
	}
}

// runGetCmd runs the get command
func runGetCmd() {
	// Get the job id
	var job string
	if len(cli.CommandArgs) > 0 {
		job = cli.CommandArgs[0]
	}

	if err := nomadClient.PrintJob(job, prettyPrintFlag); err != nil {
		cli.LogErr.Fatal(err)
	}
}

// runDelCmd runs the remove command
func runDelCmd() {
	// Get the job id
	var job string
	if len(cli.CommandArgs) > 0 {
		job = cli.CommandArgs[0]
	}

	// Delete the job
	if err := nomadClient.DeleteJob(job); err != nil {
		cli.LogErr.Fatal(err)
	} else {
		cli.LogOut.Printf("%s job is removed\n", job)
	}
}

// syncFile syncs the given file
func syncFile(file string) {
	// Read file
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		cli.LogErr.Fatal(err)
	}

	syncIt := true

	// Pull docker images
	if dockerPullFlag {
		cli.LogOut.Printf("syncing %s\n", path.Base(file))
		if err := findAndPullDockerImages(string(buf)); err != nil {
			cli.LogErr.Println("failed to sync " + path.Base(file) + " due to " + err.Error())
			syncIt = false
		}
	}

	// Add the job
	if syncIt {
		if err := nomadClient.AddJob(string(buf)); err != nil {
			cli.LogErr.Fatal(err)
		} else {
			cli.LogOut.Printf("%s is synced\n", path.Base(file))
		}
	}
}

// walkFn called for each directory during walk function execution
func walkFn(path string, info os.FileInfo, err error) error {
	// If it is not a directory then
	if !info.IsDir() {
		// Sync the file
		syncFile(path)
	}
	return nil
}

// runSyncCmd runs the sync command
func runSyncCmd() {
	// Get the file or directory path
	var path string
	if len(cli.CommandArgs) > 0 {
		path = cli.CommandArgs[0]
	}

	// Check file
	var fi os.FileInfo
	fi, err := os.Stat(path)
	if os.IsNotExist(err) {
		cli.LogErr.Fatal("no such file or directory: " + path) // fatal error
	}

	// If it is a file than
	if !fi.IsDir() {
		// Sync the file
		syncFile(path)
	} else {
		// Otherwise recursively sync files
		if err := filepath.Walk(path, walkFn); err != nil {
			cli.LogErr.Fatal(err)
		}
	}
}

// findDockerImages finds docker images from the given content
func findDockerImages(jobJSON string) ([]string, error) {

	// Check job
	buf := []byte(jobJSON)
	var sj client.SyncJob
	if err := json.Unmarshal(buf, &sj); err != nil {
		return nil, errors.New("failed to unmarshal JSON data due to " + err.Error())
	}

	// Get images
	var images []string
	if sj.Job != nil && sj.Job.TaskGroups != nil {
		for _, tg := range sj.Job.TaskGroups {
			if tg.Tasks != nil {
				for _, t := range tg.Tasks {
					if t.Config != nil && t.Driver == "docker" {
						for ck, cv := range t.Config {
							if ck == "image" {
								images = append(images, cv.(string))
							}
						}
					}
				}
			}
		}
	}

	return images, nil
}

// pullDockerImage pulls a docker image by the given name
func pullDockerImage(image string) error {
	cmd := exec.Command("docker", "pull", image)
	err := cmd.Start()
	if err != nil {
		return err
	}
	cli.LogOut.Println("pulling " + image)
	err = cmd.Wait()
	if err != nil {
		return errors.New("fail to pull " + image)
	}
	return nil
}

// findAndPullDockerImages finds and pulls docker images from the given content
func findAndPullDockerImages(jobJSON string) error {

	// Find images
	var images []string
	images, err := findDockerImages(jobJSON)
	if err != nil {
		return errors.New("failed to find docker images due to " + err.Error())
	}

	// Pull images
	for _, i := range images {
		if err := pullDockerImage(i); err != nil {
			return err
		}
	}

	return nil
}
