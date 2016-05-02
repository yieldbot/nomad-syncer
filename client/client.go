/*
 * nomad-syncer
 * Copyright (c) 2016 Yieldbot, Inc.
 * For the full copyright and license information, please view the LICENSE.txt file.
 */

// Package client provides Nomad operations
package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// Client represents the client interface
type Client struct {
	URL      string
	ProxyURL *url.URL
}

// Jobs returns the jobs
func (cl Client) Jobs() ([]Jobs, error) {

	// Get jobs
	req, err := http.NewRequest("GET", cl.nomadURL()+"/v1/jobs", nil)
	res, err := cl.doRequest(req)
	if err != nil {
		return nil, errors.New("failed to fetch jobs due to " + err.Error())
	}

	// Parse jobs
	var jobs []Jobs
	if err = json.Unmarshal(res, &jobs); err != nil {
		return nil, errors.New("failed to unmarshal JSON data due to " + err.Error())
	}

	return jobs, nil
}

// PrintJobs prints the jobs
func (cl Client) PrintJobs(pretty bool) error {

	// Get jobs
	jobs, err := cl.Jobs()
	if err != nil {
		return err
	}

	// Parse jobs
	var buf []byte

	// If pretty is true then
	if pretty {
		buf, err = json.MarshalIndent(jobs, "", "  ")
	} else {
		// Otherwise just parse it
		buf, err = json.Marshal(jobs)
	}

	if err != nil {
		return err
	}

	os.Stdout.Write(buf)

	return nil
}

// AddJob adds a job
func (cl Client) AddJob(jsonContent string) error {

	// Check job
	buf := []byte(jsonContent)
	var job SyncJob
	if err := json.Unmarshal(buf, &job); err != nil {
		return errors.New("failed to unmarshal JSON data due to " + err.Error())
	}

	// Add job
	req, err := http.NewRequest("PUT", cl.nomadURL()+"/v1/job/", bytes.NewBuffer(buf))
	req.Header.Set("Content-Type", "application/json")
	_, err = cl.doRequest(req)
	if err != nil {
		return errors.New("failed to add job due to " + err.Error())
	}

	return nil
}

// GetJob returns the job information
func (cl Client) GetJob(jobID string) (Job, error) {

	var job Job

	// Check job
	if jobID == "" {
		return job, errors.New("invalid job Id")
	}

	// Get jobs
	req, err := http.NewRequest("GET", cl.nomadURL()+"/v1/job/"+jobID, nil)
	res, err := cl.doRequest(req)
	if err != nil {
		return job, errors.New("failed to fetch jobs due to " + err.Error())
	}

	// Parse job
	if err = json.Unmarshal(res, &job); err != nil {
		return job, errors.New("failed to unmarshal JSON data due to " + err.Error())
	}

	return job, nil
}

// PrintJob prints the job
func (cl Client) PrintJob(jobID string, pretty bool) error {

	// Get job
	job, err := cl.GetJob(jobID)
	if err != nil {
		return err
	}

	// Parse job
	var buf []byte

	// If pretty is true then
	if pretty {
		buf, err = json.MarshalIndent(job, "", "  ")
	} else {
		// Otherwise just parse it
		buf, err = json.Marshal(job)
	}

	if err != nil {
		return err
	}

	os.Stdout.Write(buf)

	return nil
}

// DeleteJob deletes a job
func (cl Client) DeleteJob(jobID string) error {

	// Check job
	if jobID == "" {
		return errors.New("invalid job Id")
	}

	// Delete job
	req, err := http.NewRequest("DELETE", cl.nomadURL()+"/v1/job/"+jobID, nil)
	res, err := cl.doRequest(req)
	if err != nil {
		return errors.New("failed to delete job due to " + err.Error())
	} else if res != nil {

	}

	return nil
}

// doRequest makes a request to the REST API
func (cl Client) doRequest(req *http.Request) ([]byte, error) {

	// Init a client
	var client *http.Client

	if cl.ProxyURL != nil {
		client = &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(cl.ProxyURL)}}
	} else {
		client = &http.Client{}
	}

	// Do request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read data
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return data, errors.New("bad response: " + fmt.Sprintf("%d", resp.StatusCode))
	}

	return data, nil
}

// nomadURL returns the url
func (cl Client) nomadURL() string {
	return strings.TrimSuffix(cl.URL, "/")
}
