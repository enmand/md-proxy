package main

import (
	"fmt"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	"github.com/rackspace/gophercloud/pagination"
)

type MetaData struct {
	AdminPass        string            `json:"admin_pass"`
	AvailabilityZone string            `json:"availability_zone"`
	Hostname         string            `json:"hostname"`
	LaunchIndex      int               `json:"launch_index"`
	Name             string            `json:"name"`
	PublicKeys       map[string]string `json:"public_keys"`
	UUID             string            `json:"uuid"`
}

func FindMetadata(s servers.Server) MetaData {
	md := MetaData{
		AdminPass:        s.AdminPass,
		AvailabilityZone: "nova",
		Hostname:         fmt.Sprintf("%s.novalocal", s.Name),
		LaunchIndex:      0,
		Name:             s.Name,
		PublicKeys:       nil,
		UUID:             s.ID,
	}

	return md
}

var serverMetadata = make(map[string]MetaData)

func addServer(c *gophercloud.ServiceClient, s servers.Server) {
	md := FindMetadata(s)

	addrsPager := servers.ListAddresses(c, s.ID)

	_ = addrsPager.EachPage(func(page pagination.Page) (bool, error) {
		addrs, _ := servers.ExtractAddresses(page)

		for _, a := range addrs {
			for i := 0; i < len(a); i++ {
				serverMetadata[a[i].Address] = md
			}
		}

		return true, nil
	})

}
