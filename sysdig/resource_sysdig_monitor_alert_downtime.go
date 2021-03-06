package sysdig

import (
	"context"
	"fmt"
	"github.com/draios/terraform-provider-sysdig/sysdig/monitor"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
	"time"
)

func resourceSysdigMonitorAlertDowntime() *schema.Resource {
	timeout := 30 * time.Second

	return &schema.Resource{
		CreateContext: resourceSysdigAlertDowntimeCreate,
		UpdateContext: resourceSysdigAlertDowntimeUpdate,
		ReadContext:   resourceSysdigAlertDowntimeRead,
		DeleteContext: resourceSysdigAlertDowntimeDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(timeout),
			Update: schema.DefaultTimeout(timeout),
			Read:   schema.DefaultTimeout(timeout),
			Delete: schema.DefaultTimeout(timeout),
		},

		Schema: createAlertSchema(map[string]*schema.Schema{
			"entities_to_monitor": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"trigger_after_pct": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  100,
			},
		}),
	}
}

func resourceSysdigAlertDowntimeCreate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, err := i.(SysdigClients).sysdigMonitorClient()
	if err != nil {
		return diag.FromErr(err)
	}

	alert, err := downtimeAlertFromResourceData(data)
	if err != nil {
		return diag.FromErr(err)
	}

	alertCreated, err := client.CreateAlert(ctx, *alert)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(strconv.Itoa(alertCreated.ID))
	data.Set("version", alertCreated.Version)
	return nil
}

func resourceSysdigAlertDowntimeUpdate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, err := i.(SysdigClients).sysdigMonitorClient()
	if err != nil {
		return diag.FromErr(err)
	}

	alert, err := downtimeAlertFromResourceData(data)
	if err != nil {
		return diag.FromErr(err)
	}

	alert.ID, _ = strconv.Atoi(data.Id())

	_, err = client.UpdateAlert(ctx, *alert)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigAlertDowntimeRead(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, err := i.(SysdigClients).sysdigMonitorClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	alert, err := client.GetAlertById(ctx, id)

	if err != nil {
		data.SetId("")
		return nil
	}

	err = downtimeAlertToResourceData(&alert, data)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
func resourceSysdigAlertDowntimeDelete(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, err := i.(SysdigClients).sysdigMonitorClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.DeleteAlert(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func downtimeAlertFromResourceData(d *schema.ResourceData) (alert *monitor.Alert, err error) {
	alert, err = alertFromResourceData(d)
	if err != nil {
		return
	}

	alert.SegmentCondition = &monitor.SegmentCondition{Type: "ANY"}
	alert.Condition = fmt.Sprintf("avg(timeAvg(uptime)) <= %.2f", float64(1.0-(d.Get("trigger_after_pct").(int)/100)))

	entitiesRaw := d.Get("entities_to_monitor").([]interface{})
	for _, entityRaw := range entitiesRaw {
		alert.SegmentBy = append(alert.SegmentBy, entityRaw.(string))
	}

	return
}

func downtimeAlertToResourceData(alert *monitor.Alert, data *schema.ResourceData) (err error) {
	err = alertToResourceData(alert, data)
	if err != nil {
		return
	}

	var trigger_after_pct float64
	fmt.Sscanf(alert.Condition, "avg(timeAvg(uptime)) <= %.2f", &trigger_after_pct)
	trigger_after_pct = (1 - trigger_after_pct) * 100

	data.Set("trigger_after_pct", int(trigger_after_pct))
	data.Set("entities_to_monitor", alert.SegmentBy)

	return
}
