package schemagen

import (
	"context"

	"github.com/whyrusleeping/gosky/xrpc"
)

// schema: com.atproto.session.delete

func SessionDelete(ctx context.Context, c *xrpc.Client) error {
	if err := c.Do(ctx, xrpc.Procedure, "", "com.atproto.session.delete", nil, nil, nil); err != nil {
		return err
	}

	return nil
}
