package service

import (
	"testing"
)

func TestDeliveryAlerts(t *testing.T) {
	msg := "{\"username\": \"Chef Automate\", \"icon_url\": \"https://delivery-assets.chef.io/img/cheficon48x48.png\", \"attachments\": [{\"color\": \"#58B957\", \"fields\": [{\"value\": \"easytrace-web\", \"title\": \"Project:\"}, {\"short\": \"true\", \"value\": \"c33483458\", \"title\": \"Change submitted by:\"}], \"fallback\": \"<https://pchfdel1v.standardbank.co.za/e/chopchop/#/organizations/VirtualChannels/projects/easytrace-web/changes/3e082f5f-a637-42d6-8ce6-47b7d0a5fb4a/review|INCORRECT-BRANCH-FIX>\nVerify Passed. Change is ready for review.\"}], \"text\": \"<https://pchfdel1v.standardbank.co.za/e/chopchop/#/organizations/VirtualChannels/projects/easytrace-web/changes/3e082f5f-a637-42d6-8ce6-47b7d0a5fb4a/review|INCORRECT-BRANCH-FIX>\nVerify Passed. Change is ready for review.\"}"
	SendDeliveryAlert(msg)
}
