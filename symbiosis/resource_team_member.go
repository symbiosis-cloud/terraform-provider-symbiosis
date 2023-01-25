package symbiosis

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/symbiosis-cloud/symbiosis-go"
	"log"
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

func resourceTeamMemberCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*symbiosis.Client)
	email := d.Get("email").(string)

	_, err := client.Team.InviteMembers([]string{email}, symbiosis.UserRole(d.Get("role").(string)))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(email)
	var diags diag.Diagnostics
	return diags
}

func resourceTeamMemberUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*symbiosis.Client)
	if d.HasChange("role") {

		err := client.Team.ChangeRole(d.Id(), symbiosis.UserRole(d.Get("role").(string)))

		if err != nil {
			return diag.FromErr(err)
		}
	}
	var diags diag.Diagnostics
	return diags
}

func resourceTeamMemberDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*symbiosis.Client)

	err := client.Team.DeleteMember(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics
	return diags
}

func resourceTeamMemberRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*symbiosis.Client)

	member, err := client.Team.GetMemberByEmail(d.Id())
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

	invitation, err := client.Team.GetInvitationByEmail(d.Id())
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
