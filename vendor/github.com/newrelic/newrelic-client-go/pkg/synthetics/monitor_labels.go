package synthetics

import (
	"fmt"
	"strings"
)

// GetMonitorLabels is used to retrieve all labels for a given Synthetics monitor.
func (s *Synthetics) GetMonitorLabels(monitorID string) ([]*MonitorLabel, error) {
	url := fmt.Sprintf("/v4/monitors/%s/labels", monitorID)

	resp := getMonitorLabelsResponse{}

	_, err := s.client.Get(url, nil, &resp)
	if err != nil {
		return []*MonitorLabel{}, err
	}

	return resp.Labels, nil
}

// AddMonitorLabel is used to add a label to a given monitor.
func (s *Synthetics) AddMonitorLabel(monitorID, labelKey, labelValue string) error {
	url := fmt.Sprintf("/v4/monitors/%s/labels", monitorID)

	data := fmt.Sprintf("%s:%s", strings.Title(labelKey), strings.Title(labelValue))

	// We use RawPost here due to the Syntheics API's lack of support for JSON on
	// this call.  The values must be POSTed as bare key:value word string.
	_, err := s.client.RawPost(url, nil, data, nil)
	if err != nil {
		return err
	}

	return nil
}

// DeleteMonitorLabel deletes a key:value label from the given Syntheics monitor.
func (s *Synthetics) DeleteMonitorLabel(monitorID, labelKey, labelValue string) error {
	url := fmt.Sprintf("/v4/monitors/%s/labels/%s:%s", monitorID, strings.Title(labelKey), strings.Title(labelValue))

	_, err := s.client.Delete(url, nil, nil)
	if err != nil {
		return err
	}

	return nil
}

type getMonitorLabelsResponse struct {
	Labels []*MonitorLabel `json:"labels"`
	Count  int             `json:"count"`
}
