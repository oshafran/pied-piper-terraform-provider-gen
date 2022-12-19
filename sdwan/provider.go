package sdwan

import (
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	sdwanAPI "github.com/oshafran/pied-piper-openapi-client-go"
)

func auth(client *http.Client, host string, username string, password string) string {
	apiUrl := host 
	{
		resource := "/j_security_check"
		data := url.Values{}
		data.Set("j_username", username)
		data.Set("j_password", password)
		u, _ := url.ParseRequestURI(apiUrl)
		u.Path = resource
		urlStr := u.String()                                                               // "https://api.com/user/"
		r, _ := http.NewRequest(http.MethodPost, urlStr, strings.NewReader(data.Encode())) // URL-encoded payload
		r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		resp, _ := client.Do(r)
		fmt.Println(resp.Status)
	}
	{
		req, err := http.NewRequest("GET", fmt.Sprintf("%s/dataservice/client/token", host), nil)
		if err != nil {
			panic(err)
		}
		resp, err := client.Do(req)
		body, err := ioutil.ReadAll(resp.Body)
		fmt.Println(string(body), err)
    token = string(body)
    return string(body)
	}
}

// Ensure the implementation satisfies the expected interfaces
var (
	_ provider.Provider = &sdwanProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New() provider.Provider {
	return &sdwanProvider{}
}

// sdwanProvider is the provider implementation.
type sdwanProvider struct{}

// sdwanProviderModel maps provider schema data to a Go type.
type sdwanProviderModel struct {
	Host     types.String `tfsdk:"host"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

// Metadata returns the provider type name.
func (p *sdwanProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "sdwan"
}

// GetSchema defines the provider-level schema for configuration data.
func (p *sdwanProvider) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "Interact with sdwan.",
		Attributes: map[string]tfsdk.Attribute{
			"host": {
				Description: "URI for sdwan API. May also be provided via HASHICUPS_HOST environment variable.",
				Type:        types.StringType,
				Optional:    true,
			},
			"username": {
				Description: "Username for sdwan API. May also be provided via HASHICUPS_USERNAME environment variable.",
				Type:        types.StringType,
				Optional:    true,
			},
			"password": {
				Description: "Password for sdwan API. May also be provided via HASHICUPS_PASSWORD environment variable.",
				Type:        types.StringType,
				Optional:    true,
				Sensitive:   true,
			},
		},
	}, nil
}

// Configure prepares a sdwan API client for data sources and resources.
func (p *sdwanProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring sdwan client")

	// Retrieve provider data from configuration
	var config sdwanProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown sdwan API Host",
			"The provider cannot create the sdwan API client as there is an unknown configuration value for the HashiCups API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the sdwan_HOST environment variable.",
		)
	}

	if config.Username.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Unknown sdwan API Username",
			"The provider cannot create the sdwan API client as there is an unknown configuration value for the HashiCups API username. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the sdwan_USERNAME environment variable.",
		)
	}

	if config.Password.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Unknown sdwan API Password",
			"The provider cannot create the sdwan API client as there is an unknown configuration value for the HashiCups API password. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the sdwan_PASSWORD environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	host := os.Getenv("sdwan_HOST")
	username := os.Getenv("sdwan_USERNAME")
	password := os.Getenv("sdwan_PASSWORD")

	if !config.Host.IsNull() {
		host = config.Host.Value
	}

	if !config.Username.IsNull() {
		username = config.Username.Value
	}

	if !config.Password.IsNull() {
		password = config.Password.Value
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing sdwan API Host",
			"The provider cannot create the sdwan API client as there is a missing or empty value for the HashiCups API host. "+
				"Set the host value in the configuration or use the sdwan_HOST environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if username == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Missing sdwan API Username",
			"The provider cannot create the sdwan API client as there is a missing or empty value for the HashiCups API username. "+
				"Set the username value in the configuration or use the sdwan_USERNAME environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if password == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Missing sdwan API Password",
			"The provider cannot create the sdwan API client as there is a missing or empty value for the HashiCups API password. "+
				"Set the password value in the configuration or use the sdwan_PASSWORD environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "sdwan_host", host)
	ctx = tflog.SetField(ctx, "sdwan_username", username)
	ctx = tflog.SetField(ctx, "sdwan_password", password)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "sdwan_password")

	tflog.Debug(ctx, "Creating sdwan client")

	// Create a new sdwan client using the configuration values
	jar, err := cookiejar.New(nil)
	if err != nil {
		// error handling
	}
	httpClient := &http.Client{
		Jar: jar,
	}
	httpClient.Transport = http.DefaultTransport
	httpClient.Transport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	auth(httpClient, host, username, password)
	// return
	configuration := sdwanAPI.NewConfiguration()
	configuration.Servers = sdwanAPI.ServerConfigurations{
		{
			URL:         fmt.Sprintf("%s/dataservice", host),
			Description: "No description provided",
		},
	}
  configuration.DefaultHeader = map[string]string{
    "X-XSRF-TOKEN": token,
  }
	configuration.HTTPClient = httpClient
	client := sdwanAPI.NewAPIClient(configuration)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create sdwan API Client",
			"An unexpected error occurred when creating the sdwan API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"sdwan Client Error: "+err.Error(),
		)
		return
	}

	// Make the sdwan client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured sdwan client", map[string]any{"success": true})
}

// DataSources defines the data sources implemented in the provider.
func (p *sdwanProvider) DataSources(_ context.Context) []func() datasource.DataSource {
  return nil;
}

// Resources defines the resources implemented in the provider.
func (p *sdwanProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewVpnSiteListResource,
	}
}
