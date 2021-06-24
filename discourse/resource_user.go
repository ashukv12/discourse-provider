package discourse

import (
	"context"	
	"terraform-provider-discourse/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"log"
	"time"
	"strings"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceUserRead,
		CreateContext: resourceUserCreate,
		UpdateContext: resourceUserUpdate,
		DeleteContext: resourceUserDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceUserImporter,
		},
		Schema: map[string]*schema.Schema{
			"username": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"user_id": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"admin" : &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"email" : &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"active" : &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default: true,
			},
		},
	}
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	time.Sleep(60*time.Second)
	var diags diag.Diagnostics
	apiClient := m.(*client.Client)
	email := d.Id()
	retryErr := resource.Retry(2*time.Minute, func() *resource.RetryError {
		user, err := apiClient.GetUser(email)
		if err != nil {
			if apiClient.IsRetry(err) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		d.Set("user_id", user.Id)
		d.Set("username", user.Username)
		d.Set("name", user.Name)
		d.Set("admin", user.Admin)
		d.Set("active", user.Active)
		d.Set("email", user.Email)
		return nil
	})
	if retryErr!=nil {
		if strings.Contains(retryErr.Error(), "User Does Not Exist")==true {
			d.SetId("")
			return diags
		}
		return diag.FromErr(retryErr)
	}
	return diags
}

func resourceUserCreate(ctx context.Context,d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	apiClient := m.(*client.Client)
	email := d.Get("email").(string)
	var err error
	retryErr := resource.Retry(2*time.Minute, func() *resource.RetryError {
		if err = apiClient.NewUser(email); err != nil {
			if apiClient.IsRetry(err) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if retryErr != nil {
		time.Sleep(2 * time.Second)
		return diag.FromErr(retryErr)
	}
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(email)
	return diags
}

func resourceUserUpdate(ctx context.Context,d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var _ diag.Diagnostics
	apiClient := m.(*client.Client)
	var diags diag.Diagnostics
	if d.HasChange("email") {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "User not allowed to change email",
			Detail:   "User not allowed to change email",
		})
		return diags
	}
	user := client.User{
		Email:   d.Get("email").(string),
		Username:  d.Get("username").(string),
		Name:      d.Get("name").(string),
	}
	if d.HasChange("active") && d.Get("active").(bool) == true {
		err  := apiClient.ActivateUser(d.Get("user_id").(int))
		if err != nil{
			log.Println("[ERROR]: ",err)
			return diag.FromErr(err)
		}
	}  else if d.HasChange("active") && d.Get("active").(bool) == false {
		err := apiClient.DeactivateUser(d.Get("user_id").(int))
		if err != nil{
			log.Println("[ERROR]: ",err)
			return diag.FromErr(err)
		}
	}
	var err error
	retryErr := resource.Retry(2*time.Minute, func() *resource.RetryError {
		if err = apiClient.UpdateUser(&user); err != nil {
			if apiClient.IsRetry(err) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if retryErr != nil {
		time.Sleep(2 * time.Second)
		return diag.FromErr(retryErr)
	}
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func resourceUserDelete(ctx context.Context,d *schema.ResourceData, m interface{}) diag.Diagnostics {
 	var diags diag.Diagnostics
 	apiClient := m.(*client.Client)
 	username := d.Get("username").(string)
	 var err error
	retryErr := resource.Retry(2*time.Minute, func() *resource.RetryError {
		if err = apiClient.DeleteUser(username); err != nil {
			if apiClient.IsRetry(err) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if retryErr != nil {
		time.Sleep(2 * time.Second)
		return diag.FromErr(retryErr)
	}
	if err != nil {
		return diag.FromErr(err)
	}
 	d.SetId("")
 	return diags
}

func resourceUserImporter(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	apiClient := m.(*client.Client)
	email := d.Id()
	body, err := apiClient.GetUser(email)
	if err!=nil{
		return nil, err
	}
	d.Set("user_id", body.Id)
	d.Set("username", body.Username)
	d.Set("name", body.Name)
	d.Set("admin", body.Admin)
	d.Set("active", body.Active)
	d.Set("email", body.Email)
	return []*schema.ResourceData{d}, nil
}