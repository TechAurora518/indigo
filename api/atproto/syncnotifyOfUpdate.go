// Code generated by cmd/lexgen (see Makefile's lexgen); DO NOT EDIT.

package atproto

// schema: com.atproto.sync.notifyOfUpdate

import (
	"context"

	"github.com/bluesky-social/indigo/xrpc"
)

// SyncNotifyOfUpdate calls the XRPC method "com.atproto.sync.notifyOfUpdate".
//
// hostname: Hostname of the service that is notifying of update.
func SyncNotifyOfUpdate(ctx context.Context, c *xrpc.Client, hostname string) error {

	params := map[string]interface{}{
		"hostname": hostname,
	}
	if err := c.Do(ctx, xrpc.Query, "", "com.atproto.sync.notifyOfUpdate", params, nil, nil); err != nil {
		return err
	}

	return nil
}
