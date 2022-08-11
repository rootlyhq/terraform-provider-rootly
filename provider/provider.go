package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	"github.com/rootlyhq/terraform-provider-rootly/client"
)

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown

	// Customize the content of descriptions when output. For example you can add defaults on
	// to the exported descriptions if present.
	// schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
	// 	desc := s.Description
	// 	if s.Default != nil {
	// 		desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
	// 	}
	// 	return strings.TrimSpace(desc)
	// }
}

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"api_host": {
					Description: "The Rootly API host. Defaults to https://api.rootly.com",
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("ROOTLY_API_URL", "https://api.rootly.com"),
				},
				"api_token": {
					Description: "The Rootly API Token. Generate it from your account at https://rootly.com/account",
					Type:        schema.TypeString,
					Optional:    true,
					Sensitive:   true,
					DefaultFunc: schema.EnvDefaultFunc("ROOTLY_API_TOKEN", nil),
				},
			},
			DataSourcesMap: map[string]*schema.Resource{
				"rootly_causes": dataSourceCauses(),
				"rootly_severities": dataSourceSeverities(),
			},
			ResourcesMap: map[string]*schema.Resource{
				"rootly_cause": resourceCause(),
				"rootly_custom_field": resourceCustomField(),
				"rootly_custom_field_option": resourceCustomFieldOption(),
				"rootly_environment": resourceEnvironment(),
				"rootly_functionality": resourceFunctionality(),
				"rootly_incident_role": resourceIncidentRole(),
				"rootly_incident_type": resourceIncidentType(),
				"rootly_service": resourceService(),
				"rootly_severity": resourceSeverity(),
				"rootly_team": resourceTeam(),
				"rootly_workflow_group": resourceWorkflowGroup(),
				"rootly_workflow_incident": resourceWorkflowIncident(),
				"rootly_workflow_action_item": resourceWorkflowActionItem(),
				"rootly_workflow_alert": resourceWorkflowAlert(),
				"rootly_workflow_pulse": resourceWorkflowPulse(),
				"rootly_workflow_post_mortem": resourceWorkflowPostMortem(),
				"rootly_workflow_task_add_action_item": resourceWorkflowTaskAddActionItem(),
				"rootly_workflow_task_add_role": resourceWorkflowTaskAddRole(),
				"rootly_workflow_task_add_slack_bookmark": resourceWorkflowTaskAddSlackBookmark(),
				"rootly_workflow_task_add_team": resourceWorkflowTaskAddTeam(),
				"rootly_workflow_task_add_to_timeline": resourceWorkflowTaskAddToTimeline(),
				"rootly_workflow_task_archive_slack_channels": resourceWorkflowTaskArchiveSlackChannels(),
				"rootly_workflow_task_attach_datadog_dashboards": resourceWorkflowTaskAttachDatadogDashboards(),
				"rootly_workflow_task_auto_assign_role_opsgenie": resourceWorkflowTaskAutoAssignRoleOpsgenie(),
				"rootly_workflow_task_auto_assign_role_pagerduty": resourceWorkflowTaskAutoAssignRolePagerduty(),
				"rootly_workflow_task_auto_assign_role_victor_ops": resourceWorkflowTaskAutoAssignRoleVictorOps(),
				"rootly_workflow_task_call_people": resourceWorkflowTaskCallPeople(),
				"rootly_workflow_task_create_airtable_table_record": resourceWorkflowTaskCreateAirtableTableRecord(),
				"rootly_workflow_task_create_asana_subtask": resourceWorkflowTaskCreateAsanaSubtask(),
				"rootly_workflow_task_create_asana_task": resourceWorkflowTaskCreateAsanaTask(),
				"rootly_workflow_task_create_confluence_page": resourceWorkflowTaskCreateConfluencePage(),
				"rootly_workflow_task_create_datadog_notebook": resourceWorkflowTaskCreateDatadogNotebook(),
				"rootly_workflow_task_create_dropbox_paper_page": resourceWorkflowTaskCreateDropboxPaperPage(),
				"rootly_workflow_task_create_github_issue": resourceWorkflowTaskCreateGithubIssue(),
				"rootly_workflow_task_create_google_calendar_event": resourceWorkflowTaskCreateGoogleCalendarEvent(),
				"rootly_workflow_task_update_google_calendar_event": resourceWorkflowTaskUpdateGoogleCalendarEvent(),
				"rootly_workflow_task_create_google_docs_page": resourceWorkflowTaskCreateGoogleDocsPage(),
				"rootly_workflow_task_create_google_meeting": resourceWorkflowTaskCreateGoogleMeeting(),
				"rootly_workflow_task_create_incident": resourceWorkflowTaskCreateIncident(),
				"rootly_workflow_task_create_jira_issue": resourceWorkflowTaskCreateJiraIssue(),
				"rootly_workflow_task_create_jira_subtask": resourceWorkflowTaskCreateJiraSubtask(),
				"rootly_workflow_task_create_linear_issue": resourceWorkflowTaskCreateLinearIssue(),
				"rootly_workflow_task_create_linear_subtask_issue": resourceWorkflowTaskCreateLinearSubtaskIssue(),
				"rootly_workflow_task_create_microsoft_teams_meeting": resourceWorkflowTaskCreateMicrosoftTeamsMeeting(),
				"rootly_workflow_task_create_notion_page": resourceWorkflowTaskCreateNotionPage(),
				"rootly_workflow_task_create_service_now_incident": resourceWorkflowTaskCreateServiceNowIncident(),
				"rootly_workflow_task_create_shortcut_story": resourceWorkflowTaskCreateShortcutStory(),
				"rootly_workflow_task_create_shortcut_task": resourceWorkflowTaskCreateShortcutTask(),
				"rootly_workflow_task_create_trello_card": resourceWorkflowTaskCreateTrelloCard(),
				"rootly_workflow_task_create_webex_meeting": resourceWorkflowTaskCreateWebexMeeting(),
				"rootly_workflow_task_create_zendesk_ticket": resourceWorkflowTaskCreateZendeskTicket(),
				"rootly_workflow_task_create_zoom_meeting": resourceWorkflowTaskCreateZoomMeeting(),
				"rootly_workflow_task_get_github_commits": resourceWorkflowTaskGetGithubCommits(),
				"rootly_workflow_task_get_gitlab_commits": resourceWorkflowTaskGetGitlabCommits(),
				"rootly_workflow_task_get_pulses": resourceWorkflowTaskGetPulses(),
				"rootly_workflow_task_http_client": resourceWorkflowTaskHttpClient(),
				"rootly_workflow_task_invite_to_slack_channel_opsgenie": resourceWorkflowTaskInviteToSlackChannelOpsgenie(),
				"rootly_workflow_task_invite_to_slack_channel_pagerduty": resourceWorkflowTaskInviteToSlackChannelPagerduty(),
				"rootly_workflow_task_invite_to_slack_channel": resourceWorkflowTaskInviteToSlackChannel(),
				"rootly_workflow_task_invite_to_slack_channel_victor_ops": resourceWorkflowTaskInviteToSlackChannelVictorOps(),
				"rootly_workflow_task_page_opsgenie_on_call_responders": resourceWorkflowTaskPageOpsgenieOnCallResponders(),
				"rootly_workflow_task_page_pagerduty_on_call_responders": resourceWorkflowTaskPagePagerdutyOnCallResponders(),
				"rootly_workflow_task_page_victor_ops_on_call_responders": resourceWorkflowTaskPageVictorOpsOnCallResponders(),
				"rootly_workflow_task_print": resourceWorkflowTaskPrint(),
				"rootly_workflow_task_publish_incident": resourceWorkflowTaskPublishIncident(),
				"rootly_workflow_task_redis_client": resourceWorkflowTaskRedisClient(),
				"rootly_workflow_task_rename_slack_channel": resourceWorkflowTaskRenameSlackChannel(),
				"rootly_workflow_task_run_command_heroku": resourceWorkflowTaskRunCommandHeroku(),
				"rootly_workflow_task_send_email": resourceWorkflowTaskSendEmail(),
				"rootly_workflow_task_send_slack_message": resourceWorkflowTaskSendSlackMessage(),
				"rootly_workflow_task_send_sms": resourceWorkflowTaskSendSms(),
				"rootly_workflow_task_snapshot_datadog_graph": resourceWorkflowTaskSnapshotDatadogGraph(),
				"rootly_workflow_task_snapshot_grafana_dashboard": resourceWorkflowTaskSnapshotGrafanaDashboard(),
				"rootly_workflow_task_snapshot_looker_look": resourceWorkflowTaskSnapshotLookerLook(),
				"rootly_workflow_task_snapshot_new_relic_graph": resourceWorkflowTaskSnapshotNewRelicGraph(),
				"rootly_workflow_task_tweet_twitter_message": resourceWorkflowTaskTweetTwitterMessage(),
				"rootly_workflow_task_update_airtable_table_record": resourceWorkflowTaskUpdateAirtableTableRecord(),
				"rootly_workflow_task_update_asana_task": resourceWorkflowTaskUpdateAsanaTask(),
				"rootly_workflow_task_update_github_issue": resourceWorkflowTaskUpdateGithubIssue(),
				"rootly_workflow_task_update_incident": resourceWorkflowTaskUpdateIncident(),
				"rootly_workflow_task_update_jira_issue": resourceWorkflowTaskUpdateJiraIssue(),
				"rootly_workflow_task_update_linear_issue": resourceWorkflowTaskUpdateLinearIssue(),
				"rootly_workflow_task_update_service_now_incident": resourceWorkflowTaskUpdateServiceNowIncident(),
				"rootly_workflow_task_update_shortcut_story": resourceWorkflowTaskUpdateShortcutStory(),
				"rootly_workflow_task_update_shortcut_task": resourceWorkflowTaskUpdateShortcutTask(),
				"rootly_workflow_task_update_slack_channel_topic": resourceWorkflowTaskUpdateSlackChannelTopic(),
				"rootly_workflow_task_update_status": resourceWorkflowTaskUpdateStatus(),
				"rootly_workflow_task_update_trello_card": resourceWorkflowTaskUpdateTrelloCard(),
				"rootly_workflow_task_update_zendesk_ticket": resourceWorkflowTaskUpdateZendeskTicket(),
			},
		}

		p.ConfigureContextFunc = configure(version, p)

		return p
	}
}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(c context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		host := d.Get("api_host").(string)
		token := d.Get("api_token").(string)

		// Warning or errors can be collected in a slice type
		var diags diag.Diagnostics

		cli, err := client.NewClient(host, token, RootlyUserAgent(version))
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to create Rootly client",
				Detail:   "Unable to authenticate user for authenticated Rootly client",
			})

			return nil, diags
		}

		return cli, diags
	}
}

func RootlyUserAgent(version string) string {
	return fmt.Sprintf("Rootly Terraform Provider/%s (+https://www.terraform.io) Terraform Plugin SDK/%s", version, meta.SDKVersionString())
}
