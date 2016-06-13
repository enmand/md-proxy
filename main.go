package main

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/fasthttp"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
	"github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	"github.com/rackspace/gophercloud/pagination"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	port   = kingpin.Flag("port", "Port Number").Short('p').Default(":8080").String()
	region = kingpin.Flag("region", "OpenStack region name").Short('r').Default("regionOne").String()
)

func main() {
	kingpin.Version("0.1")
	kingpin.Parse()

	provider := initOpenStackProvider()
	fetchServers(provider)

	e := echo.New()
	osg := e.Group("/openstack")
	osg_2012 := osg.Group("/2012-08-10")
	osg_2013 := osg.Group("/2013-04-04")
	osg_latest := osg.Group("/latest")
	ec2g := e.Group("/ec2")
	ec2g_2009 := ec2g.Group("/2009-04-04")
	ec2g_latest := ec2g.Group("/latest")
	initRoutes(osg_2012, osg_2013, osg_latest, ec2g_2009, ec2g_latest)

	e.Run(fasthttp.New(fmt.Sprintf("%s", *port)))
}

func initOpenStackProvider() *gophercloud.ProviderClient {
	opts, err := openstack.AuthOptionsFromEnv()
	if err != nil {
		panic("Invalid authentication options: Please configure your environment with the proper OS_ variables")
	}

	provider, err := openstack.AuthenticatedClient(opts)
	if err != nil {
		panic("Could not create OpenStack Provider. Are your credentials correct?")
	}

	return provider
}

func initWebServer() *echo.Echo {
	e := echo.New()

	return e
}

func fetchServers(p *gophercloud.ProviderClient) {
	c, err := openstack.NewComputeV2(p, gophercloud.EndpointOpts{
		Region: *region,
	})
	if err != nil {
		panic(err)
	}
	pager := servers.List(c, servers.ListOpts{})
	_ = pager.EachPage(func(page pagination.Page) (bool, error) {
		list, _ := servers.ExtractServers(page)
		for _, s := range list {
			addServer(c, s)
		}

		return true, nil
	})
}
