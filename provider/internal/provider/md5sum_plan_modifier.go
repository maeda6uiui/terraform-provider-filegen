package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/maeda6uiui/terraform-provider-filegen/internal/client"
)

// requiresReplaceOnMd5sumMismatch plans md5sum as the hash the file is expected
// to have for the planned size and filler, and requests replacement when the
// hash recorded in state (refreshed from disk by Read) does not match it.
// Without this, md5sum would simply be planned as its prior state value, and any
// change made to the file outside of Terraform would go unnoticed.
type requiresReplaceOnMd5sumMismatch struct{}

func (m requiresReplaceOnMd5sumMismatch) Description(ctx context.Context) string {
	return "Recreates the file if its hash differs from the expected one"
}

func (m requiresReplaceOnMd5sumMismatch) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m requiresReplaceOnMd5sumMismatch) PlanModifyString(
	ctx context.Context,
	req planmodifier.StringRequest,
	resp *planmodifier.StringResponse) {
	// Nothing to compare on creation and destruction
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}

	var data FilledFileResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.Size.IsUnknown() || data.Filler.IsUnknown() {
		resp.PlanValue = types.StringUnknown()
		return
	}

	expected_md5sum := client.NewClient().CalcFilledFileMd5sum(
		int(data.Size.ValueInt32()),
		byte(data.Filler.ValueInt32()),
	)
	resp.PlanValue = types.StringValue(expected_md5sum)

	if req.StateValue.ValueString() != expected_md5sum {
		resp.RequiresReplace = true
	}
}
