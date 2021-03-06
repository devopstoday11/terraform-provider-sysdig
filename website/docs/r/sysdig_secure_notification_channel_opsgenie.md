---
layout: "sysdig"
page_title: "Sysdig: sysdig_secure_notification_channel"
sidebar_current: "docs-sysdig-secure-notification-channel"
description: |-
  Creates a Sysdig Secure Notification Channel of type OpsGenie.
---

# sysdig\_secure\_notification\_channel\_opsgenie

Creates a Sysdig Secure Notification Channel of type OpsGenie.

~> **Note:** This resource is still experimental, and is subject of being changed.

## Example usage

```hcl
resource "sysdig_secure_notification_channel_opsgenie" "sample-opsgenie" {
	name                    = "Example Channel - OpsGenie"
	enabled                 = true
	api_key                 = "2349324-342354353-5324-23"
	notify_when_ok          = false
	notify_when_resolved    = false
}
```

## Argument Reference

* `name` - (Required) The name of the Notification Channel. Must be unique.

* `api_key` - (Required) Key for the API.

* `enabled` - (Optional) If false, the channel will not emit notifications. Default is true.

* `notify_when_ok` - (Optional) Send a new notification when the alert condition is 
    no longer triggered. Default is false.

* `notify_when_resolved` - (Optional) Send a new notification when the alert is manually 
    acknowledged by a user. Default is false.

* `send_test_notification` - (Optional) Send an initial test notification to check
    if the notification channel is working. Default is false.
