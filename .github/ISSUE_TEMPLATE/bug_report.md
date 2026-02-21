---
name: Bug report
about: Create a report to help us improve
title: ''
labels: ''
assignees: ''

---

Please note that the Terraform issue tracker is reserved for bug reports and feature requests. For general usage questions, please see the Terraform documentation: https://developer.hashicorp.com/terraform/docs

### Terraform and provider versions

Run `terraform version` to show the versions and include the output here.

### Affected resource(s)

Please list the affected resources as a list. For example:
- rootly_schedule
- rootly_schedule_rotation

### Terraform configuration

```hcl
# Paste your Terraform configuration here. For large Terraform config, please attach as a text file or link to a GitHub Gist.
```

IMPORTANT: Redact any sensitive values like access keys.

### Debug output

Please link to a GitHub Gist containing your complete debug output. Please do not paste the debug output in the issue. For information on how to enable debug logging, please see: https://developer.hashicorp.com/terraform/internals/debugging

IMPORTANT: Redact any sensitive values like access keys.

### Panic output

If you encountered a panic, please provide a link to a GitHub Gist containing the output of `crash.log`.

### Expected behavior

What was expected to happen?

### Actual behavior

What actually happened?

### Steps to reproduce

Please list the steps required to reproduce the issue. Also include information on any other software you are using to manage and apply Terraform configuration.
