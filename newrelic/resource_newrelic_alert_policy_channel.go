package newrelic

import (
	"fmt"
	"log"
	"time"

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
			// State: func(d *schema.ResourceData, data interface{}) ([]*schema.ResourceData, error) {
			// 	log.Print("\n\n\n *********************** \n\n")
			// 	log.Printf("IMPORT FUNCTION: %+v \n", d.Get("channel_id"))
			// 	log.Printf("IMPORT DATA: %+v \n\n\n", data)

			// 	return []*schema.ResourceData{}, nil
			// },
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
				Deprecated:    "use `channel_ids` argument instead",
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

	policyID := d.Get("policy_id").(int)
	channelID := d.Get("channel_id").(int)
	channelIDs := d.Get("channel_ids").([]interface{})

	if channelID == 0 && len(channelIDs) == 0 {
		return fmt.Errorf("must provide channel_id or channel_ids for resource newrelic_alert_policy_channel")
	}

	var ids []int

	if channelID != 0 {
		ids = []int{channelID}
	} else {
		ids = expandChannelIDs(channelIDs)
	}

	serializedID := serializeIDs(append([]int{policyID}, ids...))

	log.Printf("[INFO] Creating New Relic alert policy channel %s", serializedID)

	_, err := client.Alerts.UpdatePolicyChannels(policyID, ids)

	if err != nil {
		return err
	}

	d.SetId(serializedID)

	return nil
}

func resourceNewRelicAlertPolicyChannelRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient

	ids, err := parseHashedIDs(d.Id())
	if err != nil {
		return err
	}

	policyID := ids[0]
	parsedChannelIDs := ids[1:]

	log.Printf("[INFO] Reading New Relic alert policy channel %s", d.Id())

	log.Print("\n\n\n *********************** \n\n")
	log.Printf("STATE BEFORE: %+v \n", d.State())
	log.Printf("IS NEW:       %+v \n", d.IsNewResource())

	exists, err := policyChannelsExist(client, policyID, parsedChannelIDs)

	log.Printf("EXISTS?: %+v \n", exists)
	log.Printf("ERROR?: %+v \n", err)

	if err != nil {
		return err
	}

	if !exists {
		d.SetId("")
		return nil
	}

	d.Set("policy_id", policyID)

	channelID, channelIDOk := d.GetOk("channel_id")
	channelIDs, channelIDsOk := d.GetOk("channel_ids")

	gotExists, gotOkExists := d.GetOkExists("channel_id")
	hasChange := d.HasChange("channel_id")

	log.Printf("hasChange:        %+v \n", hasChange)
	log.Printf("gotExists:        %+v - %+v \n", gotOkExists, gotExists)
	log.Printf("channelIDOk:      %+v - %+v \n", channelIDOk, channelID)
	log.Printf("channelIDsOk:     %+v - %+v \n", channelIDsOk, channelIDs)
	log.Printf("parsedChannelIDs: %+v - %+v \n", parsedChannelIDs, len(parsedChannelIDs))
	log.Printf("SHOULD SET IDs:   %+v \n", channelIDsOk && len(parsedChannelIDs) > 0)

	// if channelIDOk && len(parsedChannelIDs) == 1 {
	// 	d.Set("channel_id", parsedChannelIDs[0])
	// }

	if channelIDsOk && len(parsedChannelIDs) > 0 {
		d.Set("channel_ids", parsedChannelIDs)
	}

	// // If importing resource, prefer `channel_ids` attribute
	// if !channelIDOk && !channelIDsOk && len(parsedChannelIDs) > 0 {
	// 	d.Set("channel_ids", parsedChannelIDs)
	// }

	log.Printf("\n\nSTATE AFTER: %+v \n", d.State())
	log.Print("\n\n *********************** \n\n")
	time.Sleep(6 * time.Second)

	return nil
}

func resourceNewRelicAlertPolicyChannelDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient

	ids, err := parseHashedIDs(d.Id())
	if err != nil {
		return err
	}

	policyID := ids[0]
	channelIDs := ids[1:]

	log.Printf("[INFO] Deleting New Relic alert policy channel %s", d.Id())

	exists, err := policyChannelsExist(client, policyID, channelIDs)
	if err != nil {
		return err
	}

	if exists {
		for _, id := range channelIDs {
			if _, err := client.Alerts.DeletePolicyChannel(policyID, id); err != nil {
				if _, ok := err.(*errors.NotFound); ok {
					return nil
				}
				return err
			}
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

func policyChannelsExist(client *newrelic.NewRelic, policyID int, channelIDs []int) (bool, error) {
	for _, id := range channelIDs {
		channelExists, err := policyChannelExists(client, policyID, id)

		if err != nil {
			return false, err
		}

		if !channelExists {
			return false, nil
		}
	}

	return true, nil
}
