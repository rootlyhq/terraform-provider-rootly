package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (r *AlertRouteResource) update(ctx context.Context, plan AlertRouteModel, data *AlertRouteModel, resp *resource.UpdateResponse) {

	// httpResp, err := r.client.UpdateAlertRouteWithApplicationVndAPIPlusJSONBodyWithResponse(ctx, data.Id.ValueString(), body)
	// if err != nil {
	// 	resp.Diagnostics.AddError("API Error", err.Error())
	// 	return
	// } else if httpResp.StatusCode() < 200 || httpResp.StatusCode() >= 300 {
	// 	resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to update, got status code: %d", httpResp.StatusCode()))
	// 	return
	// } else if httpResp.ApplicationvndApiJSON200 == nil {
	// 	resp.Diagnostics.AddError("API Error", "Unable to update, got empty response")
	// 	return
	// }

}
