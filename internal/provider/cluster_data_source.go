package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &ClusterDataSource{}
	_ datasource.DataSourceWithConfigure = &ClusterDataSource{}
)

// NewClusterRegistrationResource returns a new resource.Resource.
func NewClusterDataSource() datasource.DataSource {
	return &ClusterDataSource{}
}

// ClusterRegistrationResource is the resource implementation.
type ClusterDataSource struct {
	client *KarporClient
}

// ClusterDataSourceModel is the datasource model.
type ClusterDataSourceModel struct {
	ClusterName types.String `tfsdk:"cluster_name"`
	DisplayName types.String `tfsdk:"display_name"`
	Description types.String `tfsdk:"description"`
	Id          types.String `tfsdk:"id"`
}

// Metadata returns the metadata for the datasource.
func (d *ClusterDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cluster"
}

// Schema returns the schema for the datasource.
func (d *ClusterDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Get cluster information",
		Attributes: map[string]schema.Attribute{
			"cluster_name": schema.StringAttribute{
				Required: true,
			},
			"display_name": schema.StringAttribute{
				Computed: true,
			},
			"description": schema.StringAttribute{
				Computed: true,
			},
			"id": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

// Read reads the datasource.
func (d *ClusterDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ClusterDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	cluster, err := d.client.GetCluster(ctx, data.ClusterName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to get cluster", err.Error())
		return
	}

	state := ClusterDataSourceModel{
		ClusterName: data.ClusterName,
		DisplayName: types.StringValue(cluster.DisplayName.ValueString()),
		Description: types.StringValue(cluster.Description.ValueString()),
		Id:          types.StringValue(cluster.Id.ValueString()),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure configures the datasource.
func (d *ClusterDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*KarporClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *KarporClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}
