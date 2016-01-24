/*
 * nomad-syncer
 * Copyright (c) 2016 Yieldbot, Inc. (http://github.com/yieldbot/nomad-syncer)
 * For the full copyright and license information, please view the LICENSE.txt file.
 */

// Package app provides the app information
package app

import (
	"flag"
	"io/ioutil"
	"net/url"
	"os"
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
}

// Run runs the app
func Run() {

	// Init cli
	cli = gocli.Cli{
		AppName:    "nomad-syncer",
		AppVersion: "1.0.0",
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
func syncFile(path string) {
	// Read file
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		cli.LogErr.Fatal(err)
	}

	// Add the job
	if err := nomadClient.AddJob(string(buf)); err != nil {
		cli.LogErr.Fatal(err)
	} else {
		cli.LogOut.Printf("%s is synced\n", path)
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
