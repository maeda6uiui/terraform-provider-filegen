package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/maeda6uiui/terraform-provider-filegen/internal/client"
)

type FilledFileResource struct {
	client *client.Client
}

func NewFilledFileResource() resource.Resource {
	return &FilledFileResource{
		client: nil,
	}
}

func (r *FilledFileResource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"filename": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Filename of the output file",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"size": schema.Int32Attribute{
				Required:            true,
				MarkdownDescription: "Size of the output file",
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.RequiresReplace(),
				},
			},
			"filler": schema.Int32Attribute{
				Optional:            true,
				Computed:            true,
				Default:             int32default.StaticInt32(0),
				MarkdownDescription: "Fills the file content with this value",
				Validators: []validator.Int32{
					int32validator.Between(0, 255),
				},
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.RequiresReplace(),
				},
			},
			"md5sum": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "File hash (MD5) of the output file",
				PlanModifiers: []planmodifier.String{
					NewRequiresReplaceOnMd5sumMismatch(),
				},
			},
		},
	}
}

func (r *FilledFileResource) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client := client.NewClient()
	r.client = client
}

func (r *FilledFileResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse) {
	var data FilledFileResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	creation_resp, err := r.client.CreateFilledFile(
		data.Filename.ValueString(),
		int(data.Size.ValueInt32()),
		byte(data.Filler.ValueInt32()),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create a file, got error: %s", err),
		)
		return
	}
	data.Md5sum = types.StringValue(creation_resp.Md5sum)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FilledFileResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse) {
	var data FilledFileResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !r.client.FileExists(data.Filename.ValueString()) {
		resp.State.RemoveResource(ctx)
		return
	}

	file_hash, err := r.client.GetFileMd5sum(data.Filename.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to read the file, got error: %s", err),
		)
		return
	}

	data.Md5sum = types.StringValue(file_hash)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FilledFileResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse) {
	// Every attribute requires replacement, so this should never be called
	resp.Diagnostics.AddError(
		"Update Not Supported",
		"A filled file cannot be updated in place, it must be recreated instead",
	)
}

func (r *FilledFileResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse) {
	var data FilledFileResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteFile(data.Filename.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to delete the file, got error: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FilledFileResource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_filled_file"
}
