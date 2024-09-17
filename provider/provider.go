package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	"github.com/rootlyhq/terraform-provider-rootly/v2/client"
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
	// 		desc += fmt.Sprintf(" Defaults to %v.", s.Default)
	// 	}
	// 	return strings.TrimSpace(desc)
	// }
}

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"api_host": {
					Description: "The Rootly API host. Defaults to https://api.rootly.com. Can also be sourced from the ROOTLY_API_URL environment variable.",
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("ROOTLY_API_URL", "https://api.rootly.com"),
				},
				"api_token": {
					Description: "The Rootly API Token. Generate it from your account at https://rootly.com/account. It must be provided but can also be sourced from the ROOTLY_API_TOKEN environment variable.",
					Type:        schema.TypeString,
					Optional:    true,
					Sensitive:   true,
					DefaultFunc: schema.EnvDefaultFunc("ROOTLY_API_TOKEN", nil),
				},
			},
			DataSourcesMap: map[string]*schema.Resource{
				"rootly_alert_urgency":                    dataSourceAlertUrgency(),
				"rootly_authorization":                    dataSourceAuthorization(),
				"rootly_cause":                            dataSourceCause(),
				"rootly_custom_form":                      dataSourceCustomForm(),
				"rootly_environment":                      dataSourceEnvironment(),
				"rootly_form_field_option":                dataSourceFormFieldOption(),
				"rootly_form_field_placement_condition":   dataSourceFormFieldPlacementCondition(),
				"rootly_form_field_placement":             dataSourceFormFieldPlacement(),
				"rootly_form_field_position":              dataSourceFormFieldPosition(),
				"rootly_form_field":                       dataSourceFormField(),
				"rootly_form_set_condition":               dataSourceFormSetCondition(),
				"rootly_form_set":                         dataSourceFormSet(),
				"rootly_functionality":                    dataSourceFunctionality(),
				"rootly_workflow_group":                   dataSourceWorkflowGroup(),
				"rootly_workflow":                         dataSourceWorkflow(),
				"rootly_heartbeat":                        dataSourceHeartbeat(),
				"rootly_incident_permission_set_boolean":  dataSourceIncidentPermissionSetBoolean(),
				"rootly_incident_permission_set_resource": dataSourceIncidentPermissionSetResource(),
				"rootly_incident_permission_set":          dataSourceIncidentPermissionSet(),
				"rootly_incident_post_mortem":             dataSourceIncidentPostMortem(),
				"rootly_incident_role":                    dataSourceIncidentRole(),
				"rootly_incident_type":                    dataSourceIncidentType(),
				"rootly_incident":                         dataSourceIncident(),
				"rootly_role":                             dataSourceRole(),
				"rootly_service":                          dataSourceService(),
				"rootly_severity":                         dataSourceSeverity(),
				"rootly_status_page":                      dataSourceStatusPage(),
				"rootly_team":                             dataSourceTeam(),
				"rootly_user":                             dataSourceUser(),
				"rootly_webhooks_endpoint":                dataSourceWebhooksEndpoint(),
				"rootly_custom_field":                     dataSourceCustomField(),
				"rootly_custom_field_option":              dataSourceCustomFieldOption(),
				"rootly_causes":                           dataSourceCauses(),
				"rootly_custom_fields":                    dataSourceCustomFields(),
				"rootly_custom_field_options":             dataSourceCustomFieldOptions(),
				"rootly_environments":                     dataSourceEnvironments(),
				"rootly_functionalities":                  dataSourceFunctionalities(),
				"rootly_ip_ranges":                        dataSourceIpRanges(),
				"rootly_incident_types":                   dataSourceIncidentTypes(),
				"rootly_incident_roles":                   dataSourceIncidentRoles(),
				"rootly_retrospective_configuration":      dataSourceRetrospectiveConfiguration(),
				"rootly_workflow_task":                    dataSourceWorkflowTask(),
				"rootly_teams":                            dataSourceTeams(),
				"rootly_severities":                       dataSourceSeverities(),
				"rootly_services":                         dataSourceServices(),
			},
			ResourcesMap: map[string]*schema.Resource{
				"rootly_alert_urgency":                                    resourceAlertUrgency(),
				"rootly_alert_group":                                      resourceAlertGroup(),
				"rootly_authorization":                                    resourceAuthorization(),
				"rootly_cause":                                            resourceCause(),
				"rootly_custom_form":                                      resourceCustomForm(),
				"rootly_environment":                                      resourceEnvironment(),
				"rootly_escalation_policy":                                resourceEscalationPolicy(),
				"rootly_form_field_option":                                resourceFormFieldOption(),
				"rootly_form_field_placement_condition":                   resourceFormFieldPlacementCondition(),
				"rootly_form_field_placement":                             resourceFormFieldPlacement(),
				"rootly_form_field_position":                              resourceFormFieldPosition(),
				"rootly_form_field":                                       resourceFormField(),
				"rootly_form_set_condition":                               resourceFormSetCondition(),
				"rootly_form_set":                                         resourceFormSet(),
				"rootly_functionality":                                    resourceFunctionality(),
				"rootly_workflow_custom_field_selection":                  resourceWorkflowCustomFieldSelection(),
				"rootly_workflow_form_field_condition":                    resourceWorkflowFormFieldCondition(),
				"rootly_workflow_group":                                   resourceWorkflowGroup(),
				"rootly_heartbeat":                                        resourceHeartbeat(),
				"rootly_incident_permission_set_boolean":                  resourceIncidentPermissionSetBoolean(),
				"rootly_incident_permission_set_resource":                 resourceIncidentPermissionSetResource(),
				"rootly_incident_permission_set":                          resourceIncidentPermissionSet(),
				"rootly_incident_role_task":                               resourceIncidentRoleTask(),
				"rootly_incident_role":                                    resourceIncidentRole(),
				"rootly_incident_type":                                    resourceIncidentType(),
				"rootly_on_call_shadow":                                   resourceOnCallShadow(),
				"rootly_override_shift":                                   resourceOverrideShift(),
				"rootly_playbook_task":                                    resourcePlaybookTask(),
				"rootly_playbook":                                         resourcePlaybook(),
				"rootly_role":                                             resourceRole(),
				"rootly_schedule_rotation_active_day":                     resourceScheduleRotationActiveDay(),
				"rootly_schedule_rotation_user":                           resourceScheduleRotationUser(),
				"rootly_schedule_rotation":                                resourceScheduleRotation(),
				"rootly_schedule":                                         resourceSchedule(),
				"rootly_service":                                          resourceService(),
				"rootly_severity":                                         resourceSeverity(),
				"rootly_status_page_template":                             resourceStatusPageTemplate(),
				"rootly_status_page":                                      resourceStatusPage(),
				"rootly_team":                                             resourceTeam(),
				"rootly_webhooks_endpoint":                                resourceWebhooksEndpoint(),
				"rootly_custom_field":                                     resourceCustomField(),
				"rootly_custom_field_option":                              resourceCustomFieldOption(),
				"rootly_dashboard":                                        resourceDashboard(),
				"rootly_dashboard_panel":                                  resourceDashboardPanel(),
				"rootly_retrospective_configuration":                      resourceRetrospectiveConfiguration(),
				"rootly_retrospective_process":                            resourceRetrospectiveProcess(),
				"rootly_retrospective_step":                               resourceRetrospectiveStep(),
				"rootly_post_mortem_template":                             resourcePostmortemTemplate(),
				"rootly_secret":                                           resourceSecret(),
				"rootly_workflow_incident":                                resourceWorkflowIncident(),
				"rootly_workflow_action_item":                             resourceWorkflowActionItem(),
				"rootly_workflow_alert":                                   resourceWorkflowAlert(),
				"rootly_workflow_pulse":                                   resourceWorkflowPulse(),
				"rootly_workflow_simple":                                  resourceWorkflowSimple(),
				"rootly_workflow_post_mortem":                             resourceWorkflowPostMortem(),
				"rootly_workflow_task_add_action_item":                    resourceWorkflowTaskAddActionItem(),
				"rootly_workflow_task_update_action_item":                 resourceWorkflowTaskUpdateActionItem(),
				"rootly_workflow_task_add_role":                           resourceWorkflowTaskAddRole(),
				"rootly_workflow_task_add_slack_bookmark":                 resourceWorkflowTaskAddSlackBookmark(),
				"rootly_workflow_task_add_team":                           resourceWorkflowTaskAddTeam(),
				"rootly_workflow_task_add_to_timeline":                    resourceWorkflowTaskAddToTimeline(),
				"rootly_workflow_task_archive_slack_channels":             resourceWorkflowTaskArchiveSlackChannels(),
				"rootly_workflow_task_attach_datadog_dashboards":          resourceWorkflowTaskAttachDatadogDashboards(),
				"rootly_workflow_task_auto_assign_role_opsgenie":          resourceWorkflowTaskAutoAssignRoleOpsgenie(),
				"rootly_workflow_task_auto_assign_role_rootly":            resourceWorkflowTaskAutoAssignRoleRootly(),
				"rootly_workflow_task_auto_assign_role_pagerduty":         resourceWorkflowTaskAutoAssignRolePagerduty(),
				"rootly_workflow_task_update_pagerduty_incident":          resourceWorkflowTaskUpdatePagerdutyIncident(),
				"rootly_workflow_task_create_pagerduty_status_update":     resourceWorkflowTaskCreatePagerdutyStatusUpdate(),
				"rootly_workflow_task_create_pagertree_alert":             resourceWorkflowTaskCreatePagertreeAlert(),
				"rootly_workflow_task_update_pagertree_alert":             resourceWorkflowTaskUpdatePagertreeAlert(),
				"rootly_workflow_task_auto_assign_role_victor_ops":        resourceWorkflowTaskAutoAssignRoleVictorOps(),
				"rootly_workflow_task_call_people":                        resourceWorkflowTaskCallPeople(),
				"rootly_workflow_task_create_airtable_table_record":       resourceWorkflowTaskCreateAirtableTableRecord(),
				"rootly_workflow_task_create_asana_subtask":               resourceWorkflowTaskCreateAsanaSubtask(),
				"rootly_workflow_task_create_asana_task":                  resourceWorkflowTaskCreateAsanaTask(),
				"rootly_workflow_task_create_confluence_page":             resourceWorkflowTaskCreateConfluencePage(),
				"rootly_workflow_task_create_datadog_notebook":            resourceWorkflowTaskCreateDatadogNotebook(),
				"rootly_workflow_task_create_dropbox_paper_page":          resourceWorkflowTaskCreateDropboxPaperPage(),
				"rootly_workflow_task_create_github_issue":                resourceWorkflowTaskCreateGithubIssue(),
				"rootly_workflow_task_create_gitlab_issue":                resourceWorkflowTaskCreateGitlabIssue(),
				"rootly_workflow_task_create_outlook_event":               resourceWorkflowTaskCreateOutlookEvent(),
				"rootly_workflow_task_create_google_calendar_event":       resourceWorkflowTaskCreateGoogleCalendarEvent(),
				"rootly_workflow_task_update_google_docs_page":            resourceWorkflowTaskUpdateGoogleDocsPage(),
				"rootly_workflow_task_update_google_calendar_event":       resourceWorkflowTaskUpdateGoogleCalendarEvent(),
				"rootly_workflow_task_create_sharepoint_page":             resourceWorkflowTaskCreateSharepointPage(),
				"rootly_workflow_task_create_google_docs_page":            resourceWorkflowTaskCreateGoogleDocsPage(),
				"rootly_workflow_task_create_google_docs_permissions":     resourceWorkflowTaskCreateGoogleDocsPermissions(),
				"rootly_workflow_task_remove_google_docs_permissions":     resourceWorkflowTaskRemoveGoogleDocsPermissions(),
				"rootly_workflow_task_create_quip_page":                   resourceWorkflowTaskCreateQuipPage(),
				"rootly_workflow_task_create_google_meeting":              resourceWorkflowTaskCreateGoogleMeeting(),
				"rootly_workflow_task_create_go_to_meeting":               resourceWorkflowTaskCreateGoToMeeting(),
				"rootly_workflow_task_create_incident":                    resourceWorkflowTaskCreateIncident(),
				"rootly_workflow_task_create_incident_postmortem":         resourceWorkflowTaskCreateIncidentPostmortem(),
				"rootly_workflow_task_create_jira_issue":                  resourceWorkflowTaskCreateJiraIssue(),
				"rootly_workflow_task_create_jira_subtask":                resourceWorkflowTaskCreateJiraSubtask(),
				"rootly_workflow_task_create_linear_issue":                resourceWorkflowTaskCreateLinearIssue(),
				"rootly_workflow_task_create_linear_subtask_issue":        resourceWorkflowTaskCreateLinearSubtaskIssue(),
				"rootly_workflow_task_create_linear_issue_comment":        resourceWorkflowTaskCreateLinearIssueComment(),
				"rootly_workflow_task_create_microsoft_teams_meeting":     resourceWorkflowTaskCreateMicrosoftTeamsMeeting(),
				"rootly_workflow_task_create_notion_page":                 resourceWorkflowTaskCreateNotionPage(),
				"rootly_workflow_task_update_notion_page":                 resourceWorkflowTaskUpdateNotionPage(),
				"rootly_workflow_task_create_service_now_incident":        resourceWorkflowTaskCreateServiceNowIncident(),
				"rootly_workflow_task_create_shortcut_story":              resourceWorkflowTaskCreateShortcutStory(),
				"rootly_workflow_task_create_shortcut_task":               resourceWorkflowTaskCreateShortcutTask(),
				"rootly_workflow_task_create_trello_card":                 resourceWorkflowTaskCreateTrelloCard(),
				"rootly_workflow_task_create_webex_meeting":               resourceWorkflowTaskCreateWebexMeeting(),
				"rootly_workflow_task_create_zendesk_ticket":              resourceWorkflowTaskCreateZendeskTicket(),
				"rootly_workflow_task_create_zendesk_jira_link":           resourceWorkflowTaskCreateZendeskJiraLink(),
				"rootly_workflow_task_create_clickup_task":                resourceWorkflowTaskCreateClickupTask(),
				"rootly_workflow_task_create_motion_task":                 resourceWorkflowTaskCreateMotionTask(),
				"rootly_workflow_task_create_zoom_meeting":                resourceWorkflowTaskCreateZoomMeeting(),
				"rootly_workflow_task_get_github_commits":                 resourceWorkflowTaskGetGithubCommits(),
				"rootly_workflow_task_get_gitlab_commits":                 resourceWorkflowTaskGetGitlabCommits(),
				"rootly_workflow_task_get_pulses":                         resourceWorkflowTaskGetPulses(),
				"rootly_workflow_task_get_alerts":                         resourceWorkflowTaskGetAlerts(),
				"rootly_workflow_task_http_client":                        resourceWorkflowTaskHttpClient(),
				"rootly_workflow_task_invite_to_slack_channel_opsgenie":   resourceWorkflowTaskInviteToSlackChannelOpsgenie(),
				"rootly_workflow_task_invite_to_slack_channel_rootly":     resourceWorkflowTaskInviteToSlackChannelRootly(),
				"rootly_workflow_task_invite_to_slack_channel_pagerduty":  resourceWorkflowTaskInviteToSlackChannelPagerduty(),
				"rootly_workflow_task_invite_to_slack_channel":            resourceWorkflowTaskInviteToSlackChannel(),
				"rootly_workflow_task_invite_to_slack_channel_victor_ops": resourceWorkflowTaskInviteToSlackChannelVictorOps(),
				"rootly_workflow_task_page_opsgenie_on_call_responders":   resourceWorkflowTaskPageOpsgenieOnCallResponders(),
				"rootly_workflow_task_create_opsgenie_alert":              resourceWorkflowTaskCreateOpsgenieAlert(),
				"rootly_workflow_task_update_opsgenie_alert":              resourceWorkflowTaskUpdateOpsgenieAlert(),
				"rootly_workflow_task_update_opsgenie_incident":           resourceWorkflowTaskUpdateOpsgenieIncident(),
				"rootly_workflow_task_page_rootly_on_call_responders":     resourceWorkflowTaskPageRootlyOnCallResponders(),
				"rootly_workflow_task_page_pagerduty_on_call_responders":  resourceWorkflowTaskPagePagerdutyOnCallResponders(),
				"rootly_workflow_task_page_victor_ops_on_call_responders": resourceWorkflowTaskPageVictorOpsOnCallResponders(),
				"rootly_workflow_task_update_victor_ops_incident":         resourceWorkflowTaskUpdateVictorOpsIncident(),
				"rootly_workflow_task_print":                              resourceWorkflowTaskPrint(),
				"rootly_workflow_task_publish_incident":                   resourceWorkflowTaskPublishIncident(),
				"rootly_workflow_task_redis_client":                       resourceWorkflowTaskRedisClient(),
				"rootly_workflow_task_rename_slack_channel":               resourceWorkflowTaskRenameSlackChannel(),
				"rootly_workflow_task_change_slack_channel_privacy":       resourceWorkflowTaskChangeSlackChannelPrivacy(),
				"rootly_workflow_task_run_command_heroku":                 resourceWorkflowTaskRunCommandHeroku(),
				"rootly_workflow_task_send_email":                         resourceWorkflowTaskSendEmail(),
				"rootly_workflow_task_send_dashboard_report":              resourceWorkflowTaskSendDashboardReport(),
				"rootly_workflow_task_create_slack_channel":               resourceWorkflowTaskCreateSlackChannel(),
				"rootly_workflow_task_send_slack_message":                 resourceWorkflowTaskSendSlackMessage(),
				"rootly_workflow_task_send_sms":                           resourceWorkflowTaskSendSms(),
				"rootly_workflow_task_send_whatsapp_message":              resourceWorkflowTaskSendWhatsappMessage(),
				"rootly_workflow_task_snapshot_datadog_graph":             resourceWorkflowTaskSnapshotDatadogGraph(),
				"rootly_workflow_task_snapshot_grafana_dashboard":         resourceWorkflowTaskSnapshotGrafanaDashboard(),
				"rootly_workflow_task_snapshot_looker_look":               resourceWorkflowTaskSnapshotLookerLook(),
				"rootly_workflow_task_snapshot_new_relic_graph":           resourceWorkflowTaskSnapshotNewRelicGraph(),
				"rootly_workflow_task_tweet_twitter_message":              resourceWorkflowTaskTweetTwitterMessage(),
				"rootly_workflow_task_update_airtable_table_record":       resourceWorkflowTaskUpdateAirtableTableRecord(),
				"rootly_workflow_task_update_asana_task":                  resourceWorkflowTaskUpdateAsanaTask(),
				"rootly_workflow_task_update_github_issue":                resourceWorkflowTaskUpdateGithubIssue(),
				"rootly_workflow_task_update_gitlab_issue":                resourceWorkflowTaskUpdateGitlabIssue(),
				"rootly_workflow_task_update_incident":                    resourceWorkflowTaskUpdateIncident(),
				"rootly_workflow_task_update_incident_postmortem":         resourceWorkflowTaskUpdateIncidentPostmortem(),
				"rootly_workflow_task_update_jira_issue":                  resourceWorkflowTaskUpdateJiraIssue(),
				"rootly_workflow_task_update_linear_issue":                resourceWorkflowTaskUpdateLinearIssue(),
				"rootly_workflow_task_update_service_now_incident":        resourceWorkflowTaskUpdateServiceNowIncident(),
				"rootly_workflow_task_update_shortcut_story":              resourceWorkflowTaskUpdateShortcutStory(),
				"rootly_workflow_task_update_shortcut_task":               resourceWorkflowTaskUpdateShortcutTask(),
				"rootly_workflow_task_update_slack_channel_topic":         resourceWorkflowTaskUpdateSlackChannelTopic(),
				"rootly_workflow_task_update_status":                      resourceWorkflowTaskUpdateStatus(),
				"rootly_workflow_task_update_trello_card":                 resourceWorkflowTaskUpdateTrelloCard(),
				"rootly_workflow_task_update_clickup_task":                resourceWorkflowTaskUpdateClickupTask(),
				"rootly_workflow_task_update_motion_task":                 resourceWorkflowTaskUpdateMotionTask(),
				"rootly_workflow_task_update_zendesk_ticket":              resourceWorkflowTaskUpdateZendeskTicket(),
				"rootly_workflow_task_update_attached_alerts":             resourceWorkflowTaskUpdateAttachedAlerts(),
				"rootly_workflow_task_trigger_workflow":                   resourceWorkflowTaskTriggerWorkflow(),
				"rootly_workflow_task_send_slack_blocks":                  resourceWorkflowTaskSendSlackBlocks(),
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
