// Copyright 2014 ALTOROS
// Licensed under the AGPLv3, see LICENSE file for details.

package data

import "io"

// Meta describes properties of dataset
type Meta struct {
	Limit      int `json:"limit"`
	Offset     int `json:"offset"`
	TotalCount int `json:"total_count"`
}

// Resource describes properties of linked resource
type Resource struct {
	URI  string `json:"resource_uri"`
	UUID string `json:"uuid"`
}

// NIC describes properties of network interface card
type NIC struct {
	IPv4 struct {
		Conf string   `json:"conf"`
		IP   Resource `json:"ip"`
	} `json:"ip_v4_conf"`
	Model string   `json:"model"`
	VLAN  Resource `json:"vlan"`
}

// Drive describe properties of disk drive
type Drive struct {
	BootOrder int      `json:"boot_order"`
	Channel   string   `json:"dev_channel"`
	Device    string   `json:"device"`
	Drive     Resource `json:"drive"`
}

// ServerRecord contains main properties of cloud server instance
type ServerRecord struct {
	Name   string `json:"name"`
	URI    string `json:"resource_uri"`
	Status string `json:"status"`
	UUID   string `json:"uuid"`
}

// ServerRecords holds collection of Server objects
type ServerRecords struct {
	Meta    Meta           `json:"meta"`
	Objects []ServerRecord `json:"objects"`
}

// Server contains detail properties of cloud server instance
type Server struct {
	ServerRecord
	Cpu    int64             `json:"cpu"`
	Mem    int64             `json:"mem"`
	Meta   map[string]string `json:"meta"`
	NICs   []NIC             `json:"nics"`
	Drives []Drive           `json:"drives"`
}

// ServersInfo holds collection of ServerInfo objects
type Servers struct {
	Meta    Meta     `json:"meta"`
	Objects []Server `json:"objects"`
}

// ReadServers reads and unmarshalls information about cloud server instances from JSON stream
func ReadServers(r io.Reader) ([]Server, error) {
	var servers Servers
	if err := ReadJson(r, &servers); err != nil {
		return nil, err
	}
	return servers.Objects, nil
}

// ReadServer reads and unmarshalls information about single cloud server instance from JSON stream
func ReadServer(r io.Reader) (*Server, error) {
	var server Server
	if err := ReadJson(r, &server); err != nil {
		return nil, err
	}
	return &server, nil
}
