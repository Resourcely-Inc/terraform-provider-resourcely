package provider

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	//	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Resourcely-Inc/terraform-provider-resourcely/internal/client"
)

// Ensure ResourcelyProvider satisfies various provider interfaces.
var (
	_ provider.Provider = &ResourcelyProvider{}
)

// ResourcelyProvider defines the provider implementation.
type ResourcelyProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// ResourcelyProviderModel describes the provider data model.
type ResourcelyProviderModel struct {
	Host           types.String `tfsdk:"host"`
	AuthToken      types.String `tfsdk:"auth_token"`
	AllowedTenants types.List   `tfsdk:"allowed_tenants"`
}

func (p *ResourcelyProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "resourcely"
	resp.Version = p.version
}

func (p *ResourcelyProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Configure Resourcely resources",
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				MarkdownDescription: "URI for Resourcely API.",
				Optional:            true,
			},
			"auth_token": schema.StringAttribute{
				MarkdownDescription: "Authorization token for Resourcely API.",
				Optional:            true,
				Sensitive:           true,
			},
			"allowed_tenants": schema.ListAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "List of allowed tenant names (case-insensitive) to prevent accidently applying a configuration to the wrong one.",
				Optional:            true,
			},
		},
	}
}

func (p *ResourcelyProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Resourcely client")

	// Retrieve provider data from configuration
	var config ResourcelyProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	host := os.Getenv("RESOURCELY_HOST")
	authToken := os.Getenv("RESOURCELY_AUTH_TOKEN")
	allowedTenants := make([]string, 0)

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if !config.AuthToken.IsNull() {
		authToken = config.AuthToken.ValueString()
	}

	if !config.AllowedTenants.IsNull() {
		config.AllowedTenants.ElementsAs(ctx, &allowedTenants, false)
	}

	ctx = tflog.SetField(ctx, "resourcely_host", host)
	ctx = tflog.SetField(ctx, "resourcely_auth_token", authToken)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "resourcely_auth_token")
	ctx = tflog.SetField(ctx, "allowed_tenants", allowedTenants)

	tflog.Debug(ctx, "Creating Resourcely client")

	client, err := client.NewClient(nil, host, authToken)
	if err != nil {
		resp.Diagnostics.AddError(
			"Creating Resourcley Client Failed",
			err.Error(),
		)
	}

	err = client.Check()
	if err != nil {
		resp.Diagnostics.AddError(
			"Checking API Status Failed",
			err.Error(),
		)
	}

	if len(allowedTenants) > 0 {
		tenant, err := client.Tenant()
		if err != nil {
			resp.Diagnostics.AddError(
				"Getting Tenant Failed", err.Error(),
			)
		}

		found := false
		for _, allowedTenant := range allowedTenants {
			if tenant == allowedTenant {
				found = true
				break
			}
		}
		if !found {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Resourcely tenant not allowed: %s", tenant),
				fmt.Sprintf("Allowed tenants are %v", allowedTenants),
			)
			return
		}
	}

	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured Resourcely client", map[string]any{"success": true})
}

func (p *ResourcelyProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewBlueprintResource,
		NewContextQuestionResource,
	}
}

func (p *ResourcelyProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewBlueprintDataSource,
		NewContextQuestionDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &ResourcelyProvider{
			version: version,
		}
	}
}
