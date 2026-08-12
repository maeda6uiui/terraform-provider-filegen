package provider

import "github.com/hashicorp/terraform-plugin-framework/types"

type FilledFileResourceModel struct {
	Filename types.String `tfsdk:"filename"`
	Size     types.Int32  `tfsdk:"size"`
	Filler   types.Int32  `tfsdk:"filler"`
	Md5sum   types.String `tfsdk:"md5sum"`
}
