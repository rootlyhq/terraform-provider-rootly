package provider

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/DataDog/jsonapi"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/rootlyhq/terraform-provider-rootly/v2/internal/apiclient"
)

func (r *AlertRouteResource) create(ctx context.Context, data *AlertRouteModel, resp *resource.CreateResponse) {
	modelIn, err := data.ToClientModel(ctx)
	if err != nil {
		resp.Diagnostics.AddError("ToClientModel Error", err.Error())
		return
	}

	b, err := jsonapi.Marshal(&modelIn, jsonapi.MarshalClientMode())
	if err != nil {
		resp.Diagnostics.AddError("JSONAPI Marshal Error", err.Error())
		return
	}

	httpResp, err := r.client.CreateAlertRouteWithBodyWithResponse(ctx, "application/vnd.api+json", bytes.NewReader(b))
	if err != nil {
		resp.Diagnostics.AddError("API Error", err.Error())
		return
	} else if httpResp.StatusCode() < 200 || httpResp.StatusCode() >= 300 {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to create, got status code: %d", httpResp.StatusCode()))
		return
	} else if httpResp.Body == nil {
		resp.Diagnostics.AddError("API Error", "Unable to create, got empty response")
		return
	}

	var modelOut apiclient.AlertRouteModel
	if err := jsonapi.Unmarshal(httpResp.Body, &modelOut); err != nil {
		resp.Diagnostics.AddError("JSONAPI Unmarshal Error", err.Error())
		return
	}

	if err := FillAlertRouteModel(ctx, modelOut, data); err != nil {
		resp.Diagnostics.AddError("FillModel Error", err.Error())
		return
	}
}

func (r *AlertRouteResource) read(ctx context.Context, data *AlertRouteModel, resp *resource.ReadResponse) {
	httpResp, err := r.client.GetAlertRouteWithResponse(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error", err.Error())
		return
	} else if httpResp.StatusCode() < 200 || httpResp.StatusCode() >= 300 {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to read, got status code: %d", httpResp.StatusCode()))
		return
	} else if httpResp.Body == nil {
		resp.Diagnostics.AddError("API Error", "Unable to read, got empty response")
		return
	}

	var modelOut apiclient.AlertRouteModel
	if err := jsonapi.Unmarshal(httpResp.Body, &modelOut); err != nil {
		resp.Diagnostics.AddError("JSONAPI Unmarshal Error", err.Error())
		return
	}

	if err := FillAlertRouteModel(ctx, modelOut, data); err != nil {
		resp.Diagnostics.AddError("FillModel Error", err.Error())
		return
	}
}

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

func (r *AlertRouteResource) delete(ctx context.Context, data *AlertRouteModel, resp *resource.DeleteResponse) {
	httpResp, err := r.client.DeleteAlertRouteWithResponse(
		ctx,
		data.Id.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError("API Error", err.Error())
		return
	} else if httpResp.StatusCode() == http.StatusNotFound {
		return
	} else if httpResp.StatusCode() < 200 || httpResp.StatusCode() >= 300 {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to delete, got status code: %d", httpResp.StatusCode()))
		return
	}
}
