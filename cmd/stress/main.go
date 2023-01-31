package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	mathrand "math/rand"
	"os"
	"sync"
	"time"

	"github.com/bluesky-social/indigo/api"
	atproto "github.com/bluesky-social/indigo/api/atproto"
	bsky "github.com/bluesky-social/indigo/api/bsky"
	"github.com/bluesky-social/indigo/carstore"
	cliutil "github.com/bluesky-social/indigo/cmd/gosky/util"
	"github.com/bluesky-social/indigo/repo"
	"github.com/bluesky-social/indigo/testing"
	"github.com/bluesky-social/indigo/util"
	"github.com/bluesky-social/indigo/xrpc"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-datastore"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	cbor "github.com/ipfs/go-ipld-cbor"
	"github.com/ipld/go-car"
	cli "github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()

	app.Commands = []*cli.Command{
		postingCmd,
		genRepoCmd,
	}

	app.RunAndExitOnError()
}

var postingCmd = &cli.Command{
	Name: "posting",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name: "quiet",
		},
		&cli.IntFlag{
			Name:  "count",
			Value: 100,
		},
		&cli.IntFlag{
			Name:  "concurrent",
			Value: 1,
		},
	},
	Action: func(cctx *cli.Context) error {
		atp, err := cliutil.GetATPClient(cctx, false)
		if err != nil {
			return err
		}

		ctx := context.TODO()

		buf := make([]byte, 6)
		rand.Read(buf)
		id := hex.EncodeToString(buf)

		var invite *string
		acc, err := atp.CreateAccount(ctx, fmt.Sprintf("user-%s@test.com", id), "user-"+id+".test", "password", invite)
		if err != nil {
			return err
		}

		quiet := cctx.Bool("quiet")

		atp.C.Auth = &xrpc.AuthInfo{
			Did:       acc.Did,
			AccessJwt: acc.AccessJwt,
			Handle:    acc.Handle,
		}

		count := cctx.Int("count")
		concurrent := cctx.Int("concurrent")

		var wg sync.WaitGroup
		for con := 0; con < concurrent; con++ {
			wg.Add(1)
			go func(worker int) {
				defer wg.Done()
				for i := 0; i < count; i++ {
					buf := make([]byte, 100)
					rand.Read(buf)

					res, err := atp.RepoCreateRecord(ctx, acc.Did, "app.bsky.feed.post", true, &api.PostRecord{
						Text:      hex.EncodeToString(buf),
						CreatedAt: time.Now().Format(time.RFC3339),
					})
					if err != nil {
						fmt.Printf("errored on worker %d loop %d: %s\n", worker, i, err)
						return
					}

					if !quiet {
						fmt.Println(res.Cid, res.Uri)
					}
				}
			}(con)
		}

		wg.Wait()

		return nil
	},
}

func randAction() string {
	v := mathrand.Intn(100)
	if v < 40 {
		return "post"
	} else if v < 60 {
		return "repost"
	} else if v < 80 {
		return "reply"
	} else {
		return "like"
	}
}

var genRepoCmd = &cli.Command{
	Name: "gen-repo",
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:  "len",
			Value: 50,
		},
	},
	Action: func(cctx *cli.Context) error {
		fname := cctx.Args().First()

		l := cctx.Int("len")

		membs := blockstore.NewBlockstore(datastore.NewMapDatastore())

		ctx := context.Background()

		r := repo.NewRepo(ctx, membs)

		words, err := testing.ReadWords()
		if err != nil {
			return err
		}

		var root cid.Cid
		for i := 0; i < l; i++ {
			switch randAction() {
			case "post":
				_, _, err := r.CreateRecord(ctx, "app.bsky.feed.post", &bsky.FeedPost{
					CreatedAt: time.Now().Format(util.ISO8601),
					Text:      testing.RandSentence(words, 200),
				})
				if err != nil {
					return err
				}
			case "repost":
				_, _, err := r.CreateRecord(ctx, "app.bsky.feed.repost", &bsky.FeedRepost{
					CreatedAt: time.Now().Format(util.ISO8601),
					Subject: &atproto.RepoStrongRef{
						Uri: testing.RandFakeAtUri("app.bsky.feed.post", ""),
						Cid: testing.RandFakeCid().String(),
					},
				})
				if err != nil {
					return err
				}
			case "reply":
				_, _, err := r.CreateRecord(ctx, "app.bsky.feed.post", &bsky.FeedPost{
					CreatedAt: time.Now().Format(util.ISO8601),
					Text:      testing.RandSentence(words, 200),
					Reply: &bsky.FeedPost_ReplyRef{
						Root: &atproto.RepoStrongRef{
							Uri: testing.RandFakeAtUri("app.bsky.feed.post", ""),
							Cid: testing.RandFakeCid().String(),
						},
						Parent: &atproto.RepoStrongRef{
							Uri: testing.RandFakeAtUri("app.bsky.feed.post", ""),
							Cid: testing.RandFakeCid().String(),
						},
					},
				})
				if err != nil {
					return err
				}
			case "like":
				_, _, err := r.CreateRecord(ctx, "app.bsky.feed.vote", &bsky.FeedVote{
					CreatedAt: time.Now().Format(util.ISO8601),
					Direction: "up",
					Subject: &atproto.RepoStrongRef{
						Uri: testing.RandFakeAtUri("app.bsky.feed.post", ""),
						Cid: testing.RandFakeCid().String(),
					},
				})
				if err != nil {
					return err
				}
			}

			nroot, err := r.Commit(ctx)
			if err != nil {
				return err
			}

			root = nroot
		}

		fi, err := os.Create(fname)
		if err != nil {
			return err
		}
		defer fi.Close()

		h := &car.CarHeader{
			Roots:   []cid.Cid{root},
			Version: 1,
		}
		hb, err := cbor.DumpObject(h)
		if err != nil {
			return err
		}

		_, err = carstore.LdWrite(fi, hb)
		if err != nil {
			return err
		}

		kc, _ := membs.AllKeysChan(ctx)
		for k := range kc {
			blk, err := membs.Get(ctx, k)
			if err != nil {
				return err
			}

			_, err = carstore.LdWrite(fi, k.Bytes(), blk.RawData())
			if err != nil {
				return err
			}
		}

		return nil
	},
}
