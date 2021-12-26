package symbiosis

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceTeamMember() *schema.Resource {
	return &schema.Resource{
    Description: `
    Manages team membership and invitations.
    `,
		CreateContext: resourceTeamMemberCreate,
		ReadContext:   resourceTeamMemberRead,
		DeleteContext: resourceTeamMemberDelete,
		UpdateContext: resourceTeamMemberUpdate,
		Schema: map[string]*schema.Schema{
			"email": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "User email to invite. Adding an team member will send the user an invitation. Deleting a team member will either delete the invitation or the user depending on whether the user has accepted the invitation.",
			},
			"accepted_invitation": {
				Type:        schema.TypeString,
        Computed:    true,
				Description: "Whether the user has accepted the invitation to the team.",
			},
			"role": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    false,
				Description: "User role. Valid values include [MEMBER, ADMIN].",
			},
		},
	}
}

type PostTeamMemberInput struct {
	Emails []string `json:"emails"`
	Role   string   `json:"role"`
}

type PutTeamMemberInput struct {
	Role   string   `json:"role"`
}

func resourceTeamMemberCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	api := meta.(*SymbiosisClient).symbiosisApi
  email := d.Get("email").(string)
	input := &PostTeamMemberInput{
		Emails: []string{
			email,
		},
		Role: d.Get("role").(string),
	}
	resp, err := api.R().SetBody(input).SetError(SymbiosisApiError{}).ForceContentType("application/json").Post("rest/v1/team/member/invite")
	if err != nil {
		return diag.FromErr(err)
	}
	if resp.StatusCode() != 200 {
		symbiosisErr := resp.Error().(*SymbiosisApiError)
		if symbiosisErr.Message != "" {
			return diag.FromErr(symbiosisErr)
		}
		return diag.FromErr(err)
	}
	d.SetId(email)
	var diags diag.Diagnostics
	return diags
}

func resourceTeamMemberUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
  api := meta.(*SymbiosisClient).symbiosisApi
  if d.HasChange("role") {
    input := &PutTeamMemberInput{
      Role: d.Get("role").(string),
    }
    _, err := api.R().SetBody(input).SetError(SymbiosisApiError{}).ForceContentType("application/json").Post(fmt.Sprintf("rest/v1/team/member/%v", d.Id()))
    if err != nil {
      return diag.FromErr(err)
    }
  }
  var diags diag.Diagnostics
  return diags
}

func resourceTeamMemberDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	api := meta.(*SymbiosisClient).symbiosisApi
	resp, err := api.R().SetError(SymbiosisApiError{}).ForceContentType("application/json").Delete(fmt.Sprintf("rest/v1/team/member/%v", d.Id()))
	if err != nil {
		return diag.FromErr(err)
	}
	if resp.StatusCode() != 200 {
		symbiosisErr := resp.Error().(*SymbiosisApiError)
		if symbiosisErr.Message != "" {
			return diag.FromErr(symbiosisErr)
		}
		return diag.FromErr(err)
	}
	var diags diag.Diagnostics
	return diags
}

func resourceTeamMemberRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*SymbiosisClient)
	member, err := client.describeTeamMember(d.Id())
  var diags diag.Diagnostics
	if err != nil {
		return diag.FromErr(err)
	}
	if member != nil {
    log.Printf("[DEBUG] member invite: %+v", member)
    d.Set("email", member.Email)
    d.Set("role", member.Role)
    d.Set("accepted_invitation", true)
    return diags
	}

  invitation, err := client.describeTeamMemberInvitation(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	if invitation != nil {
    log.Printf("[DEBUG] invite: %+v", invitation)
    d.Set("email", invitation.Email)
    d.Set("role", invitation.Role)
    d.Set("accepted_invitation", false)
    return diags
  }

  d.SetId("")
  return diags
}
