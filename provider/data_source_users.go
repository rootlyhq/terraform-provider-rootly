// Hand-maintained. Unlike the singular data source, missing emails are omitted
// rather than an error, so a for_each over emails may include not-yet-provisioned users.

package provider

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rootlyhq/terraform-provider-rootly/v5/client"
	rootlygo "github.com/rootlyhq/terraform-provider-rootly/v5/schema"
)

func dataSourceUsers() *schema.Resource {
	return &schema.Resource{
		Description: "Lists Rootly users, optionally filtered to a set of emails. Emails that do not resolve to an existing user are omitted rather than causing an error.",
		ReadContext: dataSourceUsersRead,
		Schema: map[string]*schema.Schema{
			"emails": &schema.Schema{
				Type:        schema.TypeSet,
				Description: "Optional set of emails to filter by (case-insensitive). When omitted, all users are returned.",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"users": &schema.Schema{
				Type:        schema.TypeList,
				Description: "The matching users.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"email": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"on_call_role_id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"role_id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceUsersRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	var wanted map[string]struct{}
	if raw, ok := d.GetOk("emails"); ok {
		set := raw.(*schema.Set)
		wanted = make(map[string]struct{}, set.Len())
		for _, e := range set.List() {
			wanted[strings.ToLower(e.(string))] = struct{}{}
		}
	}

	// filter[email] accepts only one email, so page through all and filter here.
	// Sort by created_at (immutable) for stable pagination.
	pageSize := 100
	sort := rootlygo.ListUsersParamsSortCreatedAt
	users := make([]map[string]interface{}, 0)

	for page := 1; ; page++ {
		pageNumber := page
		params := &rootlygo.ListUsersParams{
			PageNumber: &pageNumber,
			PageSize:   &pageSize,
			Sort:       &sort,
		}

		items, err := c.ListUsers(params)
		if err != nil {
			return diag.FromErr(err)
		}
		if len(items) == 0 {
			break
		}

		for _, it := range items {
			u, ok := it.(*client.User)
			if !ok {
				continue
			}
			if wanted != nil {
				if _, want := wanted[strings.ToLower(u.Email)]; !want {
					continue
				}
			}
			users = append(users, map[string]interface{}{
				"id":              u.ID,
				"email":           u.Email,
				"on_call_role_id": u.OnCallRoleId(),
				"role_id":         u.RoleId(),
			})
		}

		if wanted != nil && len(users) == len(wanted) {
			break
		}

		if len(items) < pageSize {
			break
		}
	}

	if err := d.Set("users", users); err != nil {
		return diag.FromErr(err)
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return nil
}
