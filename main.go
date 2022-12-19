
package main

import (
	"context"
	"terraform-provider-sdwan/sdwan"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

// Provider documentation generation.
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-name hashicups
func main() {
	providerserver.Serve(context.Background(), sdwan.New, providerserver.ServeOpts{
		Address: "hashicorp.com/edu/sdwan",
	})
}
