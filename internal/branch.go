package internal

import "context"

// newBranch branch initialization method of node 'start' in graph 'awesomeeino'
func newBranch(ctx context.Context, input ProcessedResponse) (endNode string, err error) {
	if input.TriggerProfile {
		return "extractMessages", nil
	} else {
		return "OutputHandler", nil
	}
}
