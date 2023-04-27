// Code generated by cmd/lexgen (see Makefile's lexgen); DO NOT EDIT.

package bsky

// schema: app.bsky.actor.searchActors

import (
	"context"

	"github.com/bluesky-social/indigo/xrpc"
)

// ActorSearchActors_Output is the output of a app.bsky.actor.searchActors call.
type ActorSearchActors_Output struct {
	Actors []*ActorDefs_ProfileView `json:"actors" cborgen:"actors"`
	Cursor *string                  `json:"cursor,omitempty" cborgen:"cursor,omitempty"`
}

// ActorSearchActors calls the XRPC method "app.bsky.actor.searchActors".
func ActorSearchActors(ctx context.Context, c *xrpc.Client, cursor string, limit int64, term string) (*ActorSearchActors_Output, error) {
	var out ActorSearchActors_Output

	params := map[string]interface{}{
		"cursor": cursor,
		"limit":  limit,
		"term":   term,
	}
	if err := c.Do(ctx, xrpc.Query, "", "app.bsky.actor.searchActors", params, nil, &out); err != nil {
		return nil, err
	}

	return &out, nil
}
