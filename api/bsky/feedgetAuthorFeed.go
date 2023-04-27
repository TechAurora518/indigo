// Code generated by cmd/lexgen (see Makefile's lexgen); DO NOT EDIT.

package bsky

import (
	"context"

	"github.com/bluesky-social/indigo/xrpc"
)

// schema: app.bsky.feed.getAuthorFeed

type FeedGetAuthorFeed_Output struct {
	Cursor *string                  `json:"cursor,omitempty" cborgen:"cursor,omitempty"`
	Feed   []*FeedDefs_FeedViewPost `json:"feed" cborgen:"feed"`
}

func FeedGetAuthorFeed(ctx context.Context, c *xrpc.Client, actor string, cursor string, limit int64) (*FeedGetAuthorFeed_Output, error) {
	var out FeedGetAuthorFeed_Output

	params := map[string]interface{}{
		"actor":  actor,
		"cursor": cursor,
		"limit":  limit,
	}
	if err := c.Do(ctx, xrpc.Query, "", "app.bsky.feed.getAuthorFeed", params, nil, &out); err != nil {
		return nil, err
	}

	return &out, nil
}
