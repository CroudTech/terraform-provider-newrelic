package newrelic

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/newrelic/newrelic-client-go/newrelic"
	"github.com/newrelic/newrelic-client-go/pkg/errors"
)

func resourceNewRelicAlertPolicyChannel() *schema.Resource {
	return &schema.Resource{
		Create: resourceNewRelicAlertPolicyChannelCreate,
		Read:   resourceNewRelicAlertPolicyChannelRead,
		// Update: Not currently supported in API
		Delete: resourceNewRelicAlertPolicyChannelDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"policy_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"channel_id": {
				Type:          schema.TypeInt,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"channel_ids"},
				Deprecated:    "use `config` block instead",
			},
			"channel_ids": {
				Type:          schema.TypeList,
				Optional:      true,
				ForceNew:      true,
				MinItems:      1,
				ConflictsWith: []string{"channel_id"},
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
		},
	}
}

func resourceNewRelicAlertPolicyChannelCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient

	log.Print("\n\n*******************\n\n")

	policyID := d.Get("policy_id").(int)
	channelID := d.Get("channel_id").(int)
	channelIDs := d.Get("channel_ids").([]interface{})

	ids := expandChannelIDs(channelIDs)

	log.Printf("CHANNEL ID: %T - %+v", channelID, channelID)
	log.Printf("CHANNEL IDs: %T - %+v", channelIDs, channelIDs)
	log.Printf("CHANNEL EXP IDs: %T - %+v", ids, ids)

	serializedID := serializeIDs(append([]int{policyID}, ids...))

	log.Printf("[INFO] Creating New Relic alert policy channel %s", serializedID)

	resp, err := client.Alerts.UpdatePolicyChannels(policyID, ids)

	log.Printf("\n\n RESPONSE? %+v  \n\n", *resp)

	if err != nil {
		return err
	}

	d.SetId(serializedID)

	return nil
}

func resourceNewRelicAlertPolicyChannelRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient

	ids, err := parseIDs(d.Id(), 2)
	if err != nil {
		return err
	}

	policyID := ids[0]
	channelID := ids[1]

	log.Printf("[INFO] Reading New Relic alert policy channel %s", d.Id())

	exists, err := policyChannelExists(client, policyID, channelID)
	if err != nil {
		return err
	}

	if !exists {
		d.SetId("")
		return nil
	}

	d.Set("policy_id", policyID)
	d.Set("channel_id", channelID)

	return nil
}

func resourceNewRelicAlertPolicyChannelDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient

	ids, err := parseIDs(d.Id(), 2)
	if err != nil {
		return err
	}

	policyID := ids[0]
	channelID := ids[1]

	log.Printf("[INFO] Deleting New Relic alert policy channel %s", d.Id())

	exists, err := policyChannelExists(client, policyID, channelID)
	if err != nil {
		return err
	}

	if exists {
		if _, err := client.Alerts.DeletePolicyChannel(policyID, channelID); err != nil {
			if _, ok := err.(*errors.NotFound); ok {
				return nil
			}
			return err
		}
	}

	return nil
}

func policyChannelExists(client *newrelic.NewRelic, policyID int, channelID int) (bool, error) {
	channel, err := client.Alerts.GetChannel(channelID)
	if err != nil {
		if _, ok := err.(*errors.NotFound); ok {
			return false, nil
		}

		return false, err
	}

	for _, id := range channel.Links.PolicyIDs {
		if id == policyID {
			return true, nil
		}
	}

	return false, nil
}
