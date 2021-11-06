package mikrotik

import (
	"context"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIpAddressCreate,
		ReadContext:   resourceIpAddressRead,
		UpdateContext: resourceIpAddressUpdate,
		DeleteContext: resourceIpAddressDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			".id": {
				Type: schema.TypeString,
			},
			"address": {
				Type:     schema.TypeString,
				Required: true,
			},
			"comment": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"disabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"expired": {
				Type: schema.TypeBool,
			},
			"group": {
				Type: schema.TypeString,
			},
			"name": {
				Type: schema.TypeString,
			},
		},
	}
}

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	user := prepareUser(d)

	c := m.(*client.Mikrotik)

	usr, err := c.AddUser(user)

	if err != nil {
		return diag.FromErr(err)
	}

	return userToData(usr, d)
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Mikrotik)

	user, err := c.FindUser(d.Id())

	// Clear the state if the error represents that the resource no longer exists
	_, resourceMissing := err.(*client.NotFound)
	if resourceMissing && err != nil {
		d.SetId("")
		return nil
	}

	// Make sure all other errors are propagated
	if err != nil {
		return diag.FromErr(err)
	}

	return userToData(user, d)
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Mikrotik)

	User := prepareUser(d)
	User.Id = d.Id()

	usr, err := c.UpdateUser(User)

	if err != nil {
		return diag.FromErr(err)
	}

	return userToData(usr, d)
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Mikrotik)

	err := c.DeleteUser(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func userToData(user *client.User, d *schema.ResourceData) diag.Diagnostics {
	values := map[string]interface{}{
		"address":  user.Address,
		"comment":  user.Comment,
		"disabled": user.Disabled,
		"expired":  user.Expired,
		"group":    user.Group,
		"name":     user.Name,
	}

	d.SetId(user.Id)

	var diags diag.Diagnostics

	for key, value := range values {
		if err := d.Set(key, value); err != nil {
			diags = append(diags, diag.Errorf("failed to set %s: %v", key, err)...)
		}
	}

	return diags
}

func prepareUser(d *schema.ResourceData) *client.User {
	user := new(client.User)

	user.Comment = d.Get("comment").(string)
	user.Address = d.Get("address").(string)
	user.Disabled = d.Get("disabled").(bool)
	user.Expired = d.Get("expired").(bool)
	user.Group = d.Get("group").(string)
	user.Name = d.Get("name").(string)

	return user
}
