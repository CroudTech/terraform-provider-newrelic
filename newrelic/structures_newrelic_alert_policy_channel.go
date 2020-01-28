package newrelic

func expandChannelIDs(channelIDs []interface{}) []int {
	ids := make([]int, len(channelIDs))

	for i := range ids {
		ids[i] = channelIDs[i].(int)
	}

	return ids
}
