package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/client"
)

func resourceSecret() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSecretCreate,
		ReadContext:   resourceSecretRead,
		UpdateContext: resourceSecretUpdate,
		DeleteContext: resourceSecretDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    false,
				Required:    true,
				Optional:    false,
				ForceNew:    false,
				Description: "The name of the secret",
			},

			"secret": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    false,
				Required:    true,
				Optional:    false,
				ForceNew:    false,
				Description: "The redacted secret",
			},

			"hashicorp_vault_mount": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    true,
				ForceNew:    false,
				Description: "The HashiCorp Vault secret mount path",
			},

			"hashicorp_vault_path": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    true,
				ForceNew:    false,
				Description: "The HashiCorp Vault secret path",
			},

			"hashicorp_vault_version": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Required:    false,
				Optional:    true,
				ForceNew:    false,
				Description: "The HashiCorp Vault secret version",
			},
		},
	}
}

func resourceSecretCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	tflog.Trace(ctx, fmt.Sprintf("Creating Secret"))

	s := &client.Secret{}

	if value, ok := d.GetOkExists("name"); ok {
		s.Name = value.(string)
	}
	if value, ok := d.GetOkExists("secret"); ok {
		s.Secret = value.(string)
	}
	if value, ok := d.GetOkExists("hashicorp_vault_mount"); ok {
		s.HashicorpVaultMount = value.(string)
	}
	if value, ok := d.GetOkExists("hashicorp_vault_path"); ok {
		s.HashicorpVaultPath = value.(string)
	}
	if value, ok := d.GetOkExists("hashicorp_vault_version"); ok {
		s.HashicorpVaultVersion = value.(int)
	}

	res, err := c.CreateSecret(s)
	if err != nil {
		return diag.Errorf("Error creating secret: %s", err.Error())
	}

	d.SetId(res.ID)
	tflog.Trace(ctx, fmt.Sprintf("created a secret resource: %s", d.Id()))

	return resourceSecretRead(ctx, d, meta)
}

func resourceSecretRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Reading Secret: %s", d.Id()))

	item, err := c.GetSecret(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("Secret (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error reading secret: %s", d.Id())
	}

	d.Set("name", item.Name)
	d.Set("secret", item.Secret)
	d.Set("hashicorp_vault_mount", item.HashicorpVaultMount)
	d.Set("hashicorp_vault_path", item.HashicorpVaultPath)
	d.Set("hashicorp_vault_version", item.HashicorpVaultVersion)

	return nil
}

func resourceSecretUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Updating Secret: %s", d.Id()))

	s := &client.Secret{}

	if d.HasChange("name") {
		s.Name = d.Get("name").(string)
	}
	if d.HasChange("secret") {
		s.Secret = d.Get("secret").(string)
	}
	if d.HasChange("hashicorp_vault_mount") {
		s.HashicorpVaultMount = d.Get("hashicorp_vault_mount").(string)
	}
	if d.HasChange("hashicorp_vault_path") {
		s.HashicorpVaultPath = d.Get("hashicorp_vault_path").(string)
	}
	if d.HasChange("hashicorp_vault_version") {
		s.HashicorpVaultVersion = d.Get("hashicorp_vault_version").(int)
	}

	_, err := c.UpdateSecret(d.Id(), s)
	if err != nil {
		return diag.Errorf("Error updating secret: %s", err.Error())
	}

	return resourceSecretRead(ctx, d, meta)
}

func resourceSecretDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Deleting Secret: %s", d.Id()))

	err := c.DeleteSecret(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream.
		// We just remove it from the state.
		if _, ok := err.(client.NotFoundError); ok && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("Secret (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting secret: %s", err.Error())
	}

	d.SetId("")

	return nil
}
