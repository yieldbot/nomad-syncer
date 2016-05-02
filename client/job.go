/*
 * nomad-syncer
 * Copyright (c) 2016 Yieldbot, Inc.
 * For the full copyright and license information, please view the LICENSE.txt file.
 */

package client

import (
	"time"
)

// Jobs represents jobs structure
type Jobs struct {
	ID                string
	Name              string
	Type              string
	Priority          int
	Status            string
	StatusDescription string
	CreateIndex       uint64
	ModifyIndex       uint64
}

// Job represents job structure
type Job struct {
	Region            string
	ID                string
	Name              string
	Type              string
	Priority          int
	AllAtOnce         bool
	Datacenters       []string
	Constraints       []*Constraint
	TaskGroups        []*TaskGroup
	Update            *UpdateStrategy
	Periodic          *PeriodicConfig
	Meta              map[string]string
	Status            string
	StatusDescription string
	CreateIndex       uint64
	ModifyIndex       uint64
}

// Constraint represents constraint structure
type Constraint struct {
	LTarget string
	RTarget string
	Operand string
}

// TaskGroup represents the task group field
type TaskGroup struct {
	Name          string
	Count         int
	Constraints   []*Constraint
	RestartPolicy *RestartPolicy
	Tasks         []*Task
	Meta          map[string]string
}

// RestartPolicy represents the restart policy field
type RestartPolicy struct {
	Interval         time.Duration
	Attempts         int
	Delay            time.Duration
	RestartOnSuccess bool
	Mode             string
}

// Task represents the task field
type Task struct {
	Name        string
	Driver      string
	Config      map[string]interface{}
	Env         map[string]string
	Services    []Service
	Constraints []*Constraint
	Resources   *Resources
	Meta        map[string]string
	KillTimeout time.Duration
}

// Service represents the service field
type Service struct {
	ID        string
	Name      string
	Tags      []string
	PortLabel string `mapstructure:"port"`
	Checks    []ServiceCheck
}

// ServiceCheck represents the service check field
type ServiceCheck struct {
	ID       string
	Name     string
	Type     string
	Script   string
	Path     string
	Protocol string
	Interval time.Duration
	Timeout  time.Duration
}

// Resources represents the resources field
type Resources struct {
	CPU      int
	MemoryMB int
	DiskMB   int
	IOPS     int
	Networks []*NetworkResource
}

// NetworkResource represents the network resource field
type NetworkResource struct {
	Device        string
	CIDR          string
	IP            string
	MBits         int
	ReservedPorts []Port
	DynamicPorts  []Port
	Public        bool
}

// Port represents port structure
type Port struct {
	Label string
	Value int
}

// UpdateStrategy represents the update strategy field
type UpdateStrategy struct {
	Stagger     time.Duration
	MaxParallel int
}

// PeriodicConfig represents the periodic config field
type PeriodicConfig struct {
	Enabled         bool
	Spec            string
	SpecType        string
	ProhibitOverlap bool
}

// SyncJob represents sync job structure
type SyncJob struct {
	Job *Job
}
