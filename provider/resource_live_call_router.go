// DO NOT MODIFY: This file is generated by tools/generate.js. Any changes will be overwritten during the next build.

package provider

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/rootlyhq/terraform-provider-rootly/v2/client"
	"github.com/rootlyhq/terraform-provider-rootly/v2/tools"
)

func resourceLiveCallRouter() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLiveCallRouterCreate,
		ReadContext: resourceLiveCallRouterRead,
		UpdateContext: resourceLiveCallRouterUpdate,
		DeleteContext: resourceLiveCallRouterDelete,
		Importer: &schema.ResourceImporter {
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema {
			
			"kind": &schema.Schema {
				Type: schema.TypeString,
				Default: "voicemail",
				Required: false,
				Optional: true,
				ForceNew: false,
				Description: "The kind of the live_call_router. Value must be one of `voicemail`, `live`.",
		ValidateFunc: validation.StringInSlice([]string{"voicemail", "live"}, false),
				
			},
			

				"enabled": &schema.Schema {
					Type: schema.TypeBool,
					Default: true,
					Optional: true,
					
				},
				

			"name": &schema.Schema {
				Type: schema.TypeString,
				Computed: false,
				Required: true,
				Optional: false,
				ForceNew: false,
				Description: "The name of the live_call_router",
				
			},
			

			"country_code": &schema.Schema {
				Type: schema.TypeString,
				Default: "AU",
				Required: false,
				Optional: true,
				ForceNew: false,
				Description: "The country code of the live_call_router. Value must be one of `AU`, `CA`, `NL`, `NZ`, `GB`, `US`.",
		ValidateFunc: validation.StringInSlice([]string{"AU", "CA", "NL", "NZ", "GB", "US"}, false),
				
			},
			

			"phone_type": &schema.Schema {
				Type: schema.TypeString,
				Default: "local",
				Required: false,
				Optional: true,
				ForceNew: false,
				Description: "The phone type of the live_call_router. Value must be one of `local`, `toll_free`, `mobile`.",
		ValidateFunc: validation.StringInSlice([]string{"local", "toll_free", "mobile"}, false),
				
			},
			

			"phone_number": &schema.Schema {
				Type: schema.TypeString,
				Computed: false,
				Required: true,
				Optional: false,
				ForceNew: false,
				Description: "You can select a phone number using [generate_phone_number](#//api/v1/live_call_routers/generate_phone_number) API and pass that phone number here to register",
				
			},
			

			"voicemail_greeting": &schema.Schema {
				Type: schema.TypeString,
				Computed: true,
				Required: false,
				Optional: true,
				ForceNew: false,
				Description: "The voicemail greeting of the live_call_router",
				
			},
			

			"caller_greeting": &schema.Schema {
				Type: schema.TypeString,
				Computed: true,
				Required: false,
				Optional: true,
				ForceNew: false,
				Description: "The caller greeting message of the live_call_router",
				
			},
			

			"waiting_music_url": &schema.Schema {
				Type: schema.TypeString,
				Default: "https://storage.rootly.com/twilio/voicemail/ClockworkWaltz.mp3",
				Required: false,
				Optional: true,
				ForceNew: false,
				Description: "The waiting music URL of the live_call_router. Value must be one of `https://storage.rootly.com/twilio/voicemail/ClockworkWaltz.mp3`, `https://storage.rootly.com/twilio/voicemail/ith_brahms-116-4.mp3`, `https://storage.rootly.com/twilio/voicemail/Mellotroniac_-_Flight_Of_Young_Hearts_Flute.mp3`, `https://storage.rootly.com/twilio/voicemail/BusyStrings.mp3`, `https://storage.rootly.com/twilio/voicemail/oldDog_-_endless_goodbye_%28instr.%29.mp3`, `https://storage.rootly.com/twilio/voicemail/MARKOVICHAMP-Borghestral.mp3`, `https://storage.rootly.com/twilio/voicemail/ith_chopin-15-2.mp3`.",
		ValidateFunc: validation.StringInSlice([]string{"https://storage.rootly.com/twilio/voicemail/ClockworkWaltz.mp3", "https://storage.rootly.com/twilio/voicemail/ith_brahms-116-4.mp3", "https://storage.rootly.com/twilio/voicemail/Mellotroniac_-_Flight_Of_Young_Hearts_Flute.mp3", "https://storage.rootly.com/twilio/voicemail/BusyStrings.mp3", "https://storage.rootly.com/twilio/voicemail/oldDog_-_endless_goodbye_%28instr.%29.mp3", "https://storage.rootly.com/twilio/voicemail/MARKOVICHAMP-Borghestral.mp3", "https://storage.rootly.com/twilio/voicemail/ith_chopin-15-2.mp3"}, false),
				
			},
			

		"sent_to_voicemail_delay": &schema.Schema {
			Type: schema.TypeInt,
			Computed: true,
			Required: false,
			Optional: true,
			ForceNew: false,
			Description: "The delay (seconds) after which the caller in redirected to voicemail",
			
		},
		

			"should_redirect_to_voicemail_on_no_answer": &schema.Schema {
				Type: schema.TypeBool,
				Computed: true,
				Required: false,
				Optional: true,
				Description: "This prompts the caller to choose voicemail or connect live. Value must be one of true or false",
				
			},
			

		"escalation_level_delay_in_seconds": &schema.Schema {
			Type: schema.TypeInt,
			Computed: true,
			Required: false,
			Optional: true,
			ForceNew: false,
			Description: "This overrides the delay (seconds) in escalation levels",
			
		},
		

			"should_auto_resolve_alert_on_call_end": &schema.Schema {
				Type: schema.TypeBool,
				Computed: true,
				Required: false,
				Optional: true,
				Description: "This overrides the delay (seconds) in escalation levels. Value must be one of true or false",
				
			},
			

			"alert_urgency_id": &schema.Schema {
				Type: schema.TypeString,
				Computed: true,
				Required: false,
				Optional: true,
				ForceNew: false,
				Description: "This is used in escalation paths to determine who to page",
				
			},
			

			"calling_tree_prompt": &schema.Schema {
				Type: schema.TypeString,
				Computed: true,
				Required: false,
				Optional: true,
				ForceNew: false,
				Description: "The audio instructions callers will hear when they call this number, prompting them to select from available options to route their call",
				
			},
			

				"paging_targets": &schema.Schema {
					Type: schema.TypeList,
					Computed: true,
					Required: false,
					Optional: true,
					Description: "Paging targets that callers can select from when this live call router is configured as a phone tree.",
					DiffSuppressFunc: tools.EqualIgnoringOrder,
					Elem: &schema.Resource {
						Schema: map[string]*schema.Schema {
              
			"id": &schema.Schema {
				Type: schema.TypeString,
				Computed: true,
				Required: false,
				Optional: true,
				ForceNew: false,
				Description: "The ID of paging target",
				
			},
			

			"type": &schema.Schema {
				Type: schema.TypeString,
				Default: "service",
				Required: false,
				Optional: true,
				ForceNew: false,
				Description: "The type of the paging target. Value must be one of `service`, `team`, `escalation_policy`.",
		ValidateFunc: validation.StringInSlice([]string{"service", "team", "escalation_policy"}, false),
				
			},
			

			"alert_urgency_id": &schema.Schema {
				Type: schema.TypeString,
				Computed: true,
				Required: false,
				Optional: true,
				ForceNew: false,
				Description: "This is used in escalation paths to determine who to page",
				
			},
			
						},
					},
					
				},
				

			"escalation_policy_trigger_params": &schema.Schema {
				Type: schema.TypeMap,
				Elem: &schema.Schema {
					Type: schema.TypeString,
				},
				Computed: true,
				Required: false,
				Optional: true,
				Description: "",
			},
			
		},
	}
}

func resourceLiveCallRouterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	tflog.Trace(ctx, fmt.Sprintf("Creating LiveCallRouter"))

	s := &client.LiveCallRouter{}

	  if value, ok := d.GetOkExists("kind"); ok {
				s.Kind = value.(string)
			}
    if value, ok := d.GetOkExists("enabled"); ok {
				s.Enabled = tools.Bool(value.(bool))
			}
    if value, ok := d.GetOkExists("name"); ok {
				s.Name = value.(string)
			}
    if value, ok := d.GetOkExists("country_code"); ok {
				s.CountryCode = value.(string)
			}
    if value, ok := d.GetOkExists("phone_type"); ok {
				s.PhoneType = value.(string)
			}
    if value, ok := d.GetOkExists("phone_number"); ok {
				s.PhoneNumber = value.(string)
			}
    if value, ok := d.GetOkExists("voicemail_greeting"); ok {
				s.VoicemailGreeting = value.(string)
			}
    if value, ok := d.GetOkExists("caller_greeting"); ok {
				s.CallerGreeting = value.(string)
			}
    if value, ok := d.GetOkExists("waiting_music_url"); ok {
				s.WaitingMusicUrl = value.(string)
			}
    if value, ok := d.GetOkExists("sent_to_voicemail_delay"); ok {
				s.SentToVoicemailDelay = value.(int)
			}
    if value, ok := d.GetOkExists("should_redirect_to_voicemail_on_no_answer"); ok {
				s.ShouldRedirectToVoicemailOnNoAnswer = tools.Bool(value.(bool))
			}
    if value, ok := d.GetOkExists("escalation_level_delay_in_seconds"); ok {
				s.EscalationLevelDelayInSeconds = value.(int)
			}
    if value, ok := d.GetOkExists("should_auto_resolve_alert_on_call_end"); ok {
				s.ShouldAutoResolveAlertOnCallEnd = tools.Bool(value.(bool))
			}
    if value, ok := d.GetOkExists("alert_urgency_id"); ok {
				s.AlertUrgencyId = value.(string)
			}
    if value, ok := d.GetOkExists("calling_tree_prompt"); ok {
				s.CallingTreePrompt = value.(string)
			}
    if value, ok := d.GetOkExists("paging_targets"); ok {
				s.PagingTargets = value.([]interface{})
			}
    if value, ok := d.GetOkExists("escalation_policy_trigger_params"); ok {
				if valueList, ok := value.([]interface{}); ok && len(valueList) > 0 && valueList[0] != nil {
          if mapValue, ok := valueList[0].(map[string]interface{}); ok {
    				s.EscalationPolicyTriggerParams = mapValue
          }
        }
			}

	res, err := c.CreateLiveCallRouter(s)
	if err != nil {
		return diag.Errorf("Error creating live_call_router: %s", err.Error())
	}

	d.SetId(res.ID)
	tflog.Trace(ctx, fmt.Sprintf("created a live_call_router resource: %s", d.Id()))

	return resourceLiveCallRouterRead(ctx, d, meta)
}

func resourceLiveCallRouterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Reading LiveCallRouter: %s", d.Id()))

	item, err := c.GetLiveCallRouter(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream
		// We just remove it from the state.
		if errors.Is(err, client.NewNotFoundError("")) && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("LiveCallRouter (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error reading live_call_router: %s", d.Id())
	}

	d.Set("kind", item.Kind)
  d.Set("enabled", item.Enabled)
  d.Set("name", item.Name)
  d.Set("country_code", item.CountryCode)
  d.Set("phone_type", item.PhoneType)
  d.Set("phone_number", item.PhoneNumber)
  d.Set("voicemail_greeting", item.VoicemailGreeting)
  d.Set("caller_greeting", item.CallerGreeting)
  d.Set("waiting_music_url", item.WaitingMusicUrl)
  d.Set("sent_to_voicemail_delay", item.SentToVoicemailDelay)
  d.Set("should_redirect_to_voicemail_on_no_answer", item.ShouldRedirectToVoicemailOnNoAnswer)
  d.Set("escalation_level_delay_in_seconds", item.EscalationLevelDelayInSeconds)
  d.Set("should_auto_resolve_alert_on_call_end", item.ShouldAutoResolveAlertOnCallEnd)
  d.Set("alert_urgency_id", item.AlertUrgencyId)
  d.Set("calling_tree_prompt", item.CallingTreePrompt)
  
          if item.PagingTargets != nil {
              processedItems := make([]map[string]interface{}, 0)

              for _, c := range item.PagingTargets {
                  if rawItem, ok := c.(map[string]interface{}); ok {
                      // Create a new map with only the fields defined in the schema
                      processedItem := map[string]interface{}{
                          "id": rawItem["id"],
"type": rawItem["type"],
"alert_urgency_id": rawItem["alert_urgency_id"],
                      }
                      processedItems = append(processedItems, processedItem)
                  }
              }

              d.Set("paging_targets", processedItems)
          } else {
              d.Set("paging_targets", nil)
          }
        
  singleton_list := make([]interface{}, 1, 1)
          processedItem := map[string]interface{}{
            "id": item.EscalationPolicyTriggerParams["id"],
"type": item.EscalationPolicyTriggerParams["type"],
          }
          singleton_list[0] = processedItem
          d.Set("escalation_policy_trigger_params", singleton_list)
        

	return nil
}

func resourceLiveCallRouterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Updating LiveCallRouter: %s", d.Id()))

	s := &client.LiveCallRouter{}

	  if d.HasChange("kind") {
				s.Kind = d.Get("kind").(string)
			}
    if d.HasChange("enabled") {
				s.Enabled = tools.Bool(d.Get("enabled").(bool))
			}
    if d.HasChange("name") {
				s.Name = d.Get("name").(string)
			}
    if d.HasChange("country_code") {
				s.CountryCode = d.Get("country_code").(string)
			}
    if d.HasChange("phone_type") {
				s.PhoneType = d.Get("phone_type").(string)
			}
    if d.HasChange("phone_number") {
				s.PhoneNumber = d.Get("phone_number").(string)
			}
    if d.HasChange("voicemail_greeting") {
				s.VoicemailGreeting = d.Get("voicemail_greeting").(string)
			}
    if d.HasChange("caller_greeting") {
				s.CallerGreeting = d.Get("caller_greeting").(string)
			}
    if d.HasChange("waiting_music_url") {
				s.WaitingMusicUrl = d.Get("waiting_music_url").(string)
			}
    if d.HasChange("sent_to_voicemail_delay") {
				s.SentToVoicemailDelay = d.Get("sent_to_voicemail_delay").(int)
			}
    if d.HasChange("should_redirect_to_voicemail_on_no_answer") {
				s.ShouldRedirectToVoicemailOnNoAnswer = tools.Bool(d.Get("should_redirect_to_voicemail_on_no_answer").(bool))
			}
    if d.HasChange("escalation_level_delay_in_seconds") {
				s.EscalationLevelDelayInSeconds = d.Get("escalation_level_delay_in_seconds").(int)
			}
    if d.HasChange("should_auto_resolve_alert_on_call_end") {
				s.ShouldAutoResolveAlertOnCallEnd = tools.Bool(d.Get("should_auto_resolve_alert_on_call_end").(bool))
			}
    if d.HasChange("alert_urgency_id") {
				s.AlertUrgencyId = d.Get("alert_urgency_id").(string)
			}
    if d.HasChange("calling_tree_prompt") {
				s.CallingTreePrompt = d.Get("calling_tree_prompt").(string)
			}
    if d.HasChange("paging_targets") {
				s.PagingTargets = d.Get("paging_targets").([]interface{})
			}
    if d.HasChange("escalation_policy_trigger_params") {
      		tps := d.Get("escalation_policy_trigger_params").([]interface{})
      		for _, tpsi := range tps {
      			s.EscalationPolicyTriggerParams = tpsi.(map[string]interface{})
      		}
      	}
			

	_, err := c.UpdateLiveCallRouter(d.Id(), s)
	if err != nil {
		return diag.Errorf("Error updating live_call_router: %s", err.Error())
	}

	return resourceLiveCallRouterRead(ctx, d, meta)
}

func resourceLiveCallRouterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	tflog.Trace(ctx, fmt.Sprintf("Deleting LiveCallRouter: %s", d.Id()))

	err := c.DeleteLiveCallRouter(d.Id())
	if err != nil {
		// In the case of a NotFoundError, it means the resource may have been removed upstream.
		// We just remove it from the state.
		if errors.Is(err, client.NewNotFoundError("")) && !d.IsNewResource() {
			tflog.Warn(ctx, fmt.Sprintf("LiveCallRouter (%s) not found, removing from state", d.Id()))
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting live_call_router: %s", err.Error())
	}

	d.SetId("")

	return nil
}
