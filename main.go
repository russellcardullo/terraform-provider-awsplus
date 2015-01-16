package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/russellcardullo/terraform-provider-awsplus/awsplus"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: awsplus.Provider,
	})
}
