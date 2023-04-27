// Code generated by cmd/lexgen (see Makefile's lexgen); DO NOT EDIT.

package atproto

// schema: com.atproto.sync.getRepo

import (
	"bytes"
	"context"

	"github.com/bluesky-social/indigo/xrpc"
)

// SyncGetRepo calls the XRPC method "com.atproto.sync.getRepo".
//
// did: The DID of the repo.
// earliest: The earliest commit in the commit range (not inclusive)
// latest: The latest commit in the commit range (inclusive)
func SyncGetRepo(ctx context.Context, c *xrpc.Client, did string, earliest string, latest string) ([]byte, error) {
	buf := new(bytes.Buffer)

	params := map[string]interface{}{
		"did":      did,
		"earliest": earliest,
		"latest":   latest,
	}
	if err := c.Do(ctx, xrpc.Query, "", "com.atproto.sync.getRepo", params, nil, buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
