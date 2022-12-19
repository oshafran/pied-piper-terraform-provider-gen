
package sdwan

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdwanAPI "github.com/oshafran/pied-piper-openapi-client-go"
)


func CustomControl(body map[string]interface{}, plan vpnSiteListResourceModel) map[string]interface{} {
	body["defaultAction"] = map[string]interface{}{
		"type": "reject",
	}
	body["sequences"].([]map[string]interface{})[0]["actions"] = []string{}
	entries := []map[string]interface{}{
		{
			"field": "vpnList",
			"ref":   plan.Vpn.ListId.Value,
		},
	}
	body["sequences"].([]map[string]interface{})[0]["match"] = map[string]interface{}{
		"entries": entries,
	}
	return body
}


var token string;
var (
	_ resource.Resource                = &vpnSiteListResource{}
	_ resource.ResourceWithConfigure   = &vpnSiteListResource{}
	_ resource.ResourceWithImportState = &vpnSiteListResource{}
)

func NewVpnSiteListResource() resource.Resource {
	return &vpnSiteListResource{}
}

func (d *vpnSiteListResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sdwanAPI.APIClient)
}

// vpnSiteListsResource is the data source implementation.
type vpnSiteListResource struct {
	client *sdwanAPI.APIClient
}
type vpnSiteListResourceModel struct{
  Vpn              vpnResource           `tfsdk:"vpn"`
  Site              siteResource           `tfsdk:"site"`
  Prefix              prefixResource           `tfsdk:"prefix"`
  Control              controlResource           `tfsdk:"control"`
}

type vpnResource struct{
  ListId              types.String           `tfsdk:"list_id"`
  Name              types.String           `tfsdk:"name"`
  Description              types.String           `tfsdk:"description"`
  Type              types.String           `tfsdk:"type"`
  Owner              types.String           `tfsdk:"owner"`
  LastUpdated              types.Float64           `tfsdk:"last_updated"`
  ReferenceCount              types.Float64           `tfsdk:"reference_count"`
  ReadOnly              types.Bool           `tfsdk:"read_only"`
  Version              types.String           `tfsdk:"version"`
  InfoTag              types.String           `tfsdk:"info_tag"`
  IsActivatedByVsmart              types.Bool           `tfsdk:"is_activated_by_vsmart"`
  Entries              []vpnEntriesResource           `tfsdk:"entries"`
}

type vpnEntriesResource struct{
  Vpn              types.String           `tfsdk:"vpn"`
}



type siteResource struct{
  ListId              types.String           `tfsdk:"list_id"`
  Name              types.String           `tfsdk:"name"`
  Description              types.String           `tfsdk:"description"`
  Type              types.String           `tfsdk:"type"`
  Owner              types.String           `tfsdk:"owner"`
  LastUpdated              types.Float64           `tfsdk:"last_updated"`
  ReferenceCount              types.Float64           `tfsdk:"reference_count"`
  ReadOnly              types.Bool           `tfsdk:"read_only"`
  Version              types.String           `tfsdk:"version"`
  InfoTag              types.String           `tfsdk:"info_tag"`
  IsActivatedByVsmart              types.Bool           `tfsdk:"is_activated_by_vsmart"`
  Entries              []siteEntriesResource           `tfsdk:"entries"`
}

type siteEntriesResource struct{
  SiteId              types.String           `tfsdk:"site_id"`
}



type prefixResource struct{
  ListId              types.String           `tfsdk:"list_id"`
  Name              types.String           `tfsdk:"name"`
  Description              types.String           `tfsdk:"description"`
  Type              types.String           `tfsdk:"type"`
  Owner              types.String           `tfsdk:"owner"`
  LastUpdated              types.Float64           `tfsdk:"last_updated"`
  ReferenceCount              types.Float64           `tfsdk:"reference_count"`
  ReadOnly              types.Bool           `tfsdk:"read_only"`
  Version              types.String           `tfsdk:"version"`
  InfoTag              types.String           `tfsdk:"info_tag"`
  IsActivatedByVsmart              types.Bool           `tfsdk:"is_activated_by_vsmart"`
  Entries              []prefixEntriesResource           `tfsdk:"entries"`
}

type prefixEntriesResource struct{
  IpPrefix              types.String           `tfsdk:"ip_prefix"`
}



type controlResource struct{
  DefinitionId              types.String           `tfsdk:"definition_id"`
  Name              types.String           `tfsdk:"name"`
  Type              types.String           `tfsdk:"type"`
  Description              types.String           `tfsdk:"description"`
  DefaultAction              controlDefaultActionResource           `tfsdk:"default_action"`
  Sequences              []controlSequencesResource           `tfsdk:"sequences"`
}

type controlDefaultActionResource struct{
  Type              types.String           `tfsdk:"type"`
}


type controlSequencesResource struct{
  SequenceId              types.Int64           `tfsdk:"sequence_id"`
  SequenceName              types.String           `tfsdk:"sequence_name"`
  BaseAction              types.String           `tfsdk:"base_action"`
  SequenceType              types.String           `tfsdk:"sequence_type"`
  SequenceIpType              types.String           `tfsdk:"sequence_ip_type"`
  Match              controlSequencesMatchResource           `tfsdk:"match"`
}

type controlSequencesMatchResource struct{
  Entries              []controlSequencesMatchEntriesResource           `tfsdk:"entries"`
}

type controlSequencesMatchEntriesResource struct{
  Field              types.String           `tfsdk:"field"`
}





func (r *vpnSiteListResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Metadata returns the data source type name.
func (d *vpnSiteListResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vpn_site_list"
}

func (d *vpnSiteListResource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "Fetches the list of coffees.",
		Attributes: map[string]tfsdk.Attribute{
      "vpn": {
        Description:"",
         Optional: true,
       Computed: false,
        Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{

          "list_id": {
        Description: "",
        Computed: true,
      Required: false,
        Type: types.StringType,
      },
      "name": {
        Description: "",
        Computed: false,
      Required: true,
        Type: types.StringType,
      },
      "description": {
        Description: "",
        Computed: false,
      Required: true,
        Type: types.StringType,
      },
      "type": {
        Description: "",
        Computed: false,
      Required: true,
        Type: types.StringType,
      },
      "owner": {
        Description: "",
        Computed: true,
      Required: false,
        Type: types.StringType,
      },
      "last_updated": {
        Description: "",
        Computed: true,
      Required: false,
        Type: types.Float64Type,
      },
      "reference_count": {
        Description: "",
        Computed: true,
      Required: false,
        Type: types.Float64Type,
      },
      "read_only": {
        Description: "",
        Computed: true,
      Required: false,
        Type: types.BoolType,
      },
      "version": {
        Description: "",
        Computed: true,
      Required: false,
        Type: types.StringType,
      },
      "info_tag": {
        Description: "",
        Computed: true,
      Required: false,
        Type: types.StringType,
      },
      "is_activated_by_vsmart": {
        Description: "",
        Computed: true,
      Required: false,
        Type: types.BoolType,
      },
      "entries": {
        Description:"",
        Computed: false,
         Optional: true,
       Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{

          "vpn": {
        Description: "",
        Computed: false,
      Required: true,
        Type: types.StringType,
      },
}),},
}),},
      "site": {
        Description:"",
         Optional: true,
       Computed: false,
        Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{

          "list_id": {
        Description: "",
        Computed: true,
      Required: false,
        Type: types.StringType,
      },
      "name": {
        Description: "",
        Computed: false,
      Required: true,
        Type: types.StringType,
      },
      "description": {
        Description: "",
        Computed: false,
      Required: true,
        Type: types.StringType,
      },
      "type": {
        Description: "",
        Computed: false,
      Required: true,
        Type: types.StringType,
      },
      "owner": {
        Description: "",
        Computed: true,
      Required: false,
        Type: types.StringType,
      },
      "last_updated": {
        Description: "",
        Computed: true,
      Required: false,
        Type: types.Float64Type,
      },
      "reference_count": {
        Description: "",
        Computed: true,
      Required: false,
        Type: types.Float64Type,
      },
      "read_only": {
        Description: "",
        Computed: true,
      Required: false,
        Type: types.BoolType,
      },
      "version": {
        Description: "",
        Computed: true,
      Required: false,
        Type: types.StringType,
      },
      "info_tag": {
        Description: "",
        Computed: true,
      Required: false,
        Type: types.StringType,
      },
      "is_activated_by_vsmart": {
        Description: "",
        Computed: true,
      Required: false,
        Type: types.BoolType,
      },
      "entries": {
        Description:"",
        Computed: false,
         Optional: true,
       Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{

          "site_id": {
        Description: "",
        Computed: false,
      Required: true,
        Type: types.StringType,
      },
}),},
}),},
      "prefix": {
        Description:"",
         Optional: true,
       Computed: false,
        Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{

          "list_id": {
        Description: "",
        Computed: true,
      Required: false,
        Type: types.StringType,
      },
      "name": {
        Description: "",
        Computed: false,
      Required: true,
        Type: types.StringType,
      },
      "description": {
        Description: "",
        Computed: false,
      Required: true,
        Type: types.StringType,
      },
      "type": {
        Description: "",
        Computed: false,
      Required: true,
        Type: types.StringType,
      },
      "owner": {
        Description: "",
        Computed: true,
      Required: false,
        Type: types.StringType,
      },
      "last_updated": {
        Description: "",
        Computed: true,
      Required: false,
        Type: types.Float64Type,
      },
      "reference_count": {
        Description: "",
        Computed: true,
      Required: false,
        Type: types.Float64Type,
      },
      "read_only": {
        Description: "",
        Computed: true,
      Required: false,
        Type: types.BoolType,
      },
      "version": {
        Description: "",
        Computed: true,
      Required: false,
        Type: types.StringType,
      },
      "info_tag": {
        Description: "",
        Computed: true,
      Required: false,
        Type: types.StringType,
      },
      "is_activated_by_vsmart": {
        Description: "",
        Computed: true,
      Required: false,
        Type: types.BoolType,
      },
      "entries": {
        Description:"",
        Computed: false,
         Optional: true,
       Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{

          "ip_prefix": {
        Description: "",
        Computed: false,
      Required: true,
        Type: types.StringType,
      },
}),},
}),},
      "control": {
        Description:"",
         Optional: true,
       Computed: false,
        Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{

          "definition_id": {
        Description: "",
        Computed: true,
      Required: false,
        Type: types.StringType,
      },
      "name": {
        Description: "",
        Computed: false,
      Required: true,
        Type: types.StringType,
      },
      "type": {
        Description: "",
        Computed: false,
      Required: true,
        Type: types.StringType,
      },
      "description": {
        Description: "",
        Computed: false,
      Required: true,
        Type: types.StringType,
      },
      "default_action": {
        Description:"",
         Optional: true,
       Computed: false,
        Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{

          "type": {
        Description: "",
        Computed: false,
      Required: true,
        Type: types.StringType,
      },
}),},
      "sequences": {
        Description:"",
        Computed: false,
         Optional: true,
       Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{

          "sequence_id": {
        Description: "",
        Computed: false,
      Required: true,
        Type: types.Int64Type,
      },
      "sequence_name": {
        Description: "",
        Computed: false,
      Required: true,
        Type: types.StringType,
      },
      "base_action": {
        Description: "",
        Computed: false,
      Required: true,
        Type: types.StringType,
      },
      "sequence_type": {
        Description: "",
        Computed: false,
      Required: true,
        Type: types.StringType,
      },
      "sequence_ip_type": {
        Description: "",
        Computed: false,
      Required: true,
        Type: types.StringType,
      },
      "match": {
        Description:"",
         Optional: true,
       Computed: false,
        Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{

          "entries": {
        Description:"",
        Computed: false,
         Optional: true,
       Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{

          "field": {
        Description: "",
        Computed: false,
      Required: true,
        Type: types.StringType,
      },
}),},
}),},
}),},
}),},
},
}, nil
}
func (d *vpnSiteListResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state vpnSiteListResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed order value from HashiCups
//
 vpnSiteListState := vpnSiteListResourceModel{}
  
{
	_, r, err := d.client.ConfigurationPolicyVPNListBuilderApi.GetListsById39(context.Background(), state.Vpn.ListId.Value).Execute()
	dataStr, err := ioutil.ReadAll(r.Body)
  fmt.Println(string(dataStr))
	data := map[string]interface{}{}
	json.Unmarshal(dataStr, &data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read HashiCups Coffees",
			err.Error(),
		)
		return
	}

		resp.Diagnostics.AddWarning(
			"test",
			string(dataStr),
		)
	// Map response body to model

	vpnSiteList := data

statevpnResource :=  vpnResource {

            Name:              types.String{Value: vpnSiteList["name"].(string)}, 
  Description:              types.String{Value: vpnSiteList["description"].(string)}, 
  ListId:              types.String{Value: vpnSiteList["listId"].(string)}, 
  LastUpdated:              types.Float64{Value: vpnSiteList["lastUpdated"].(float64)}, 
  Owner:              types.String{Value: vpnSiteList["owner"].(string)}, 
  ReferenceCount:              types.Float64{Value: vpnSiteList["referenceCount"].(float64)}, 
  Type:              types.String{Value: vpnSiteList["type"].(string)}, 
  ReadOnly:              types.Bool{Value: vpnSiteList["readOnly"].(bool)}, 
  Version:              types.String{Value: vpnSiteList["version"].(string)}, 
  InfoTag:              types.String{Value: vpnSiteList["infoTag"].(string)}, 
  IsActivatedByVsmart:              types.Bool{Value: vpnSiteList["isActivatedByVsmart"].(bool)}, 


        }
vpnSiteListState.Vpn = statevpnResource;
}

{
	_, r, err := d.client.ConfigurationPolicySiteListBuilderApi.GetListsById30(context.Background(), state.Site.ListId.Value).Execute()
	dataStr, err := ioutil.ReadAll(r.Body)
  fmt.Println(string(dataStr))
	data := map[string]interface{}{}
	json.Unmarshal(dataStr, &data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read HashiCups Coffees",
			err.Error(),
		)
		return
	}

		resp.Diagnostics.AddWarning(
			"test",
			string(dataStr),
		)
	// Map response body to model

	vpnSiteList := data

statesiteResource :=  siteResource {

            Name:              types.String{Value: vpnSiteList["name"].(string)}, 
  Description:              types.String{Value: vpnSiteList["description"].(string)}, 
  ListId:              types.String{Value: vpnSiteList["listId"].(string)}, 
  LastUpdated:              types.Float64{Value: vpnSiteList["lastUpdated"].(float64)}, 
  Owner:              types.String{Value: vpnSiteList["owner"].(string)}, 
  ReferenceCount:              types.Float64{Value: vpnSiteList["referenceCount"].(float64)}, 
  Type:              types.String{Value: vpnSiteList["type"].(string)}, 
  ReadOnly:              types.Bool{Value: vpnSiteList["readOnly"].(bool)}, 
  Version:              types.String{Value: vpnSiteList["version"].(string)}, 
  InfoTag:              types.String{Value: vpnSiteList["infoTag"].(string)}, 
  IsActivatedByVsmart:              types.Bool{Value: vpnSiteList["isActivatedByVsmart"].(bool)}, 


        }
vpnSiteListState.Site = statesiteResource;
}

{
	_, r, err := d.client.ConfigurationPolicyPrefixListBuilderApi.GetListsById27(context.Background(), state.Prefix.ListId.Value).Execute()
	dataStr, err := ioutil.ReadAll(r.Body)
  fmt.Println(string(dataStr))
	data := map[string]interface{}{}
	json.Unmarshal(dataStr, &data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read HashiCups Coffees",
			err.Error(),
		)
		return
	}

		resp.Diagnostics.AddWarning(
			"test",
			string(dataStr),
		)
	// Map response body to model

	vpnSiteList := data

stateprefixResource :=  prefixResource {

            Name:              types.String{Value: vpnSiteList["name"].(string)}, 
  Description:              types.String{Value: vpnSiteList["description"].(string)}, 
  ListId:              types.String{Value: vpnSiteList["listId"].(string)}, 
  LastUpdated:              types.Float64{Value: vpnSiteList["lastUpdated"].(float64)}, 
  Owner:              types.String{Value: vpnSiteList["owner"].(string)}, 
  ReferenceCount:              types.Float64{Value: vpnSiteList["referenceCount"].(float64)}, 
  Type:              types.String{Value: vpnSiteList["type"].(string)}, 
  ReadOnly:              types.Bool{Value: vpnSiteList["readOnly"].(bool)}, 
  Version:              types.String{Value: vpnSiteList["version"].(string)}, 
  InfoTag:              types.String{Value: vpnSiteList["infoTag"].(string)}, 
  IsActivatedByVsmart:              types.Bool{Value: vpnSiteList["isActivatedByVsmart"].(bool)}, 


        }
vpnSiteListState.Prefix = stateprefixResource;
}

{
	_, r, err := d.client.ConfigurationPolicyControlDefinitionBuilderApi.GetPolicyDefinition14(context.Background(), state.Control.DefinitionId.Value).Execute()
	dataStr, err := ioutil.ReadAll(r.Body)
  fmt.Println(string(dataStr))
	data := map[string]interface{}{}
	json.Unmarshal(dataStr, &data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read HashiCups Coffees",
			err.Error(),
		)
		return
	}

		resp.Diagnostics.AddWarning(
			"test",
			string(dataStr),
		)
	// Map response body to model

	vpnSiteList := data

statecontrolResource :=  controlResource {

            Name:              types.String{Value: vpnSiteList["name"].(string)}, 
  Description:              types.String{Value: vpnSiteList["description"].(string)}, 
  DefinitionId:              types.String{Value: vpnSiteList["definitionId"].(string)}, 
  Type:              types.String{Value: vpnSiteList["type"].(string)}, 


        }
vpnSiteListState.Control = statecontrolResource;
}




	state = vpnSiteListState

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}


func (r *vpnSiteListResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}
func (r *vpnSiteListResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan vpnSiteListResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan



  
{
      
body := map[string]interface{}{
  "name":        plan.Vpn.Name.Value,
  "type":        plan.Vpn.Type.Value,
}


	var entries []map[string]interface{}
	for _, item := range plan.Vpn.Entries {
		entries = append(entries, map[string]interface{}{
    // doing this will cause issues if there are multiple nested values
  "vpn":        item.Vpn.Value,
},


		)
	}

  body["entries"] = entries



  bodyStringed, _ := json.Marshal(&body)
  _, response, err := r.client.ConfigurationPolicyVPNListBuilderApi.CreatePolicyList39(context.Background()).XXSRFTOKEN(token).Body(body).Execute()
  if err != nil {
	  fmt.Fprintf(os.Stderr, "Error when calling `ConfigurationPolicyVPNListBuilderApi.CreatePolicyList39`: %v\n", err)
	  fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
  }

	responseBodyString, _ := ioutil.ReadAll(response.Body)
	// Create new order
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating order",
			"Could not create order, unexpected error: "+err.Error() + string(responseBodyString) + string(bodyStringed),

		)
		return
	}

  resp.Diagnostics.AddWarning("Response body string", string(responseBodyString))

	responseBody := map[string]interface{}{}


  fmt.Println(string(responseBodyString))

	err = json.Unmarshal(responseBodyString, &responseBody)

	// Map response body to schema and populate Computed attribute values


  plan.Vpn.ListId = types.String{Value: responseBody["listId"].(string)}


}

{
      
body := map[string]interface{}{
  "name":        plan.Site.Name.Value,
  "type":        plan.Site.Type.Value,
}


	var entries []map[string]interface{}
	for _, item := range plan.Site.Entries {
		entries = append(entries, map[string]interface{}{
    // doing this will cause issues if there are multiple nested values
  "siteId":        item.SiteId.Value,
},


		)
	}

  body["entries"] = entries



  bodyStringed, _ := json.Marshal(&body)
  _, response, err := r.client.ConfigurationPolicySiteListBuilderApi.CreatePolicyList30(context.Background()).Body(body).Execute()
  if err != nil {
	  fmt.Fprintf(os.Stderr, "Error when calling `ConfigurationPolicyVPNListBuilderApi.CreatePolicyList39`: %v\n", err)
	  fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
  }

	responseBodyString, _ := ioutil.ReadAll(response.Body)
	// Create new order
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating order",
			"Could not create order, unexpected error: "+err.Error() + string(responseBodyString) + string(bodyStringed),

		)
		return
	}

  resp.Diagnostics.AddWarning("Response body string", string(responseBodyString))

	responseBody := map[string]interface{}{}


  fmt.Println(string(responseBodyString))

	err = json.Unmarshal(responseBodyString, &responseBody)

	// Map response body to schema and populate Computed attribute values


  plan.Site.ListId = types.String{Value: responseBody["listId"].(string)}


}

{
      
body := map[string]interface{}{
  "name":        plan.Prefix.Name.Value,
  "type":        plan.Prefix.Type.Value,
}


	var entries []map[string]interface{}
	for _, item := range plan.Prefix.Entries {
		entries = append(entries, map[string]interface{}{
    // doing this will cause issues if there are multiple nested values
  "ipPrefix":        item.IpPrefix.Value,
},


		)
	}

  body["entries"] = entries



  bodyStringed, _ := json.Marshal(&body)
  _, response, err := r.client.ConfigurationPolicyPrefixListBuilderApi.CreatePolicyList27(context.Background()).Body(body).Execute()
  if err != nil {
	  fmt.Fprintf(os.Stderr, "Error when calling `ConfigurationPolicyVPNListBuilderApi.CreatePolicyList39`: %v\n", err)
	  fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
  }

	responseBodyString, _ := ioutil.ReadAll(response.Body)
	// Create new order
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating order",
			"Could not create order, unexpected error: "+err.Error() + string(responseBodyString) + string(bodyStringed),

		)
		return
	}

  resp.Diagnostics.AddWarning("Response body string", string(responseBodyString))

	responseBody := map[string]interface{}{}


  fmt.Println(string(responseBodyString))

	err = json.Unmarshal(responseBodyString, &responseBody)

	// Map response body to schema and populate Computed attribute values


  plan.Prefix.ListId = types.String{Value: responseBody["listId"].(string)}


}

{
      
body := map[string]interface{}{
  "name":        plan.Control.Name.Value,
  "type":        plan.Control.Type.Value,
  "description":        plan.Control.Description.Value,
}


	var sequences []map[string]interface{}
	for _, item := range plan.Control.Sequences {
		sequences = append(sequences, map[string]interface{}{
    // doing this will cause issues if there are multiple nested values
  "sequenceId":        item.SequenceId.Value,
  "sequenceName":        item.SequenceName.Value,
  "baseAction":        item.BaseAction.Value,
  "sequenceType":        item.SequenceType.Value,
  "sequenceIpType":        item.SequenceIpType.Value,
},


		)
	}

  body["sequences"] = sequences



  bodyStringed, _ := json.Marshal(&body)
body = CustomControl(body, plan)//pre-script
  _, response, err := r.client.ConfigurationPolicyControlDefinitionBuilderApi.CreatePolicyDefinition14(context.Background()).Body(body).Execute()
  if err != nil {
	  fmt.Fprintf(os.Stderr, "Error when calling `ConfigurationPolicyVPNListBuilderApi.CreatePolicyList39`: %v\n", err)
	  fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
  }

	responseBodyString, _ := ioutil.ReadAll(response.Body)
	// Create new order
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating order",
			"Could not create order, unexpected error: "+err.Error() + string(responseBodyString) + string(bodyStringed),

		)
		return
	}

  resp.Diagnostics.AddWarning("Response body string", string(responseBodyString))

	responseBody := map[string]interface{}{}


  fmt.Println(string(responseBodyString))

	err = json.Unmarshal(responseBodyString, &responseBody)

	// Map response body to schema and populate Computed attribute values


  plan.Control.DefinitionId = types.String{Value: responseBody["definitionId"].(string)}


}
  
	


	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *vpnSiteListResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state vpnSiteListResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Delete existing order
{
_, err := r.client.ConfigurationPolicyControlDefinitionBuilderApi.DeletePolicyDefinition14(context.Background(), state.Control.DefinitionId.Value).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ConfigurationPolicyVPNListBuilderApi.DeletePolicyList39`: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting HashiCups Order",
			"Could not delete order, unexpected error: "+err.Error(),
		)
		return
	}
}

{
_, err := r.client.ConfigurationPolicyVPNListBuilderApi.DeletePolicyList39(context.Background(), state.Vpn.ListId.Value).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ConfigurationPolicyVPNListBuilderApi.DeletePolicyList39`: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting HashiCups Order",
			"Could not delete order, unexpected error: "+err.Error(),
		)
		return
	}
}

{
_, err := r.client.ConfigurationPolicySiteListBuilderApi.DeletePolicyList30(context.Background(), state.Site.ListId.Value).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ConfigurationPolicyVPNListBuilderApi.DeletePolicyList39`: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting HashiCups Order",
			"Could not delete order, unexpected error: "+err.Error(),
		)
		return
	}
}

{
_, err := r.client.ConfigurationPolicyPrefixListBuilderApi.DeletePolicyList27(context.Background(), state.Prefix.ListId.Value).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ConfigurationPolicyVPNListBuilderApi.DeletePolicyList39`: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting HashiCups Order",
			"Could not delete order, unexpected error: "+err.Error(),
		)
		return
	}
}


}
