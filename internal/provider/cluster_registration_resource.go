package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &ClusterRegistrationResource{}
	_ resource.ResourceWithConfigure   = &ClusterRegistrationResource{}
	_ resource.ResourceWithImportState = &ClusterRegistrationResource{}
)

// NewClusterRegistrationResource returns a new resource.Resource.
func NewClusterRegistrationResource() resource.Resource {
	return &ClusterRegistrationResource{}
}

// ClusterRegistrationResource is the resource implementation.
type ClusterRegistrationResource struct {
	client *KarporClient
}

// ClusterRegistrationResourceModel is the resource model.
type ClusterRegistrationResourceModel struct {
	ClusterName types.String `tfsdk:"cluster_name"`
	DisplayName types.String `tfsdk:"display_name"`
	Credentials types.String `tfsdk:"credentials"`
	Description types.String `tfsdk:"description"`
	Id          types.String `tfsdk:"id"`
	LastUpdated types.String `tfsdk:"last_updated"`
}

// Metadata returns the resource type name.
func (r *ClusterRegistrationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cluster_registration"
}

// Schema returns the resource schema.
func (r *ClusterRegistrationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage cluster registration",
		Attributes: map[string]schema.Attribute{
			"cluster_name": schema.StringAttribute{
				Required:    true,
				Description: "Unique name for the cluster",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"display_name": schema.StringAttribute{
				Optional:    true,
				Description: "Human-readable display name",
			},
			"credentials": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "Path to kubeconfig file",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "Human-readable description",
			},
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "Unique identifier",
			},
			"last_updated": schema.StringAttribute{
				Computed:    true,
				Description: "Last updated timestamp",
			},
		},
	}
}

// Create creates the resource.
func (c *ClusterRegistrationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ClusterRegistrationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Validate the kubeconfig file
	success, err := c.client.ValidateClusterConfig(ctx, &plan)
	if !success || err != nil {
		resp.Diagnostics.AddError("Invalid kubeconfig file", err.Error())
		return
	}
	tflog.Info(ctx, "Valid kubeconfig file")

	// Register the cluster
	uid, err := c.client.RegisterCluster(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to register cluster", err.Error())
		return
	}

	// Set the resource ID (uid)
	plan.Id = types.StringValue(uid)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Save the resource state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read reads the resource.
func (c *ClusterRegistrationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state ClusterRegistrationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed order value from Karpor
	remoteState, err := c.client.GetCluster(ctx, state.ClusterName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Karpor Cluster",
			"Could not read Karpor Mananged Cluster Name "+state.ClusterName.String()+": "+err.Error(),
		)
		return
	}

	// Overwrite items with refreshed state
	state.ClusterName = remoteState.ClusterName
	state.DisplayName = remoteState.DisplayName
	state.Description = remoteState.Description
	state.Id = remoteState.Id

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource.
func (c *ClusterRegistrationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan ClusterRegistrationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the cluster
	success, err := c.client.UpdateCluster(ctx, &plan)
	if !success || err != nil {
		resp.Diagnostics.AddError("Failed to update cluster", err.Error())
		return
	}

	// Fetch updated items from GetCluster as UpdateCluster items are not populated.
	remoteState, err := c.client.GetCluster(ctx, plan.ClusterName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Karpor Cluster",
			"Could not read Karpor Cluster "+plan.ClusterName.ValueString()+": "+err.Error(),
		)
		return
	}

	// Update resource state with updated items and timestamp
	plan.DisplayName = remoteState.DisplayName
	plan.Description = remoteState.Description
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource.
func (c *ClusterRegistrationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state ClusterRegistrationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing order
	success, err := c.client.DeleteCluster(ctx, &state)
	if !success || err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Karpor Cluster",
			"Could not delete cluster, unexpected error: "+err.Error(),
		)
		return
	}
}

// ImportState imports the resource.
func (c *ClusterRegistrationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Terraform will automatically call the resource's Read method to import the rest of the Terraform state
	resource.ImportStatePassthroughID(ctx, path.Root("cluster_name"), req, resp)
}

// Configure configures the resource.
func (c *ClusterRegistrationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*KarporClient)
	if !ok {
		resp.Diagnostics.AddError("Unexpected data type",
			fmt.Sprintf("Expected *KarporClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	c.client = client
}
