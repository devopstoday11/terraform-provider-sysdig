package sysdig_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"os"
	"testing"

	"github.com/draios/terraform-provider-sysdig/sysdig"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNotificationChannelDataSource(t *testing.T) {
	rText := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			if v := os.Getenv("SYSDIG_SECURE_API_TOKEN"); v == "" {
				t.Fatal("SYSDIG_SECURE_API_TOKEN must be set for acceptance tests")
			}
		},
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: notificationChannelEmailWithNameAndDatasource(rText),
			},
		},
	})
}

func notificationChannelEmailWithNameAndDatasource(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_notification_channel_email" "sample_email" {
	name = "%s"
	enabled = true
	recipients = ["root@localhost.com"]
	notify_when_ok = false
	notify_when_resolved = false
}

data "sysdig_secure_notification_channel" "sample_email" {
	depends_on = [sysdig_secure_notification_channel_email.sample_email]
	name = sysdig_secure_notification_channel_email.sample_email.name
}
`, name)
}
