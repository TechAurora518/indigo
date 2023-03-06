// Code generated by github.com/whyrusleeping/cbor-gen. DO NOT EDIT.

package events

import (
	"fmt"
	"io"
	"math"
	"sort"

	cid "github.com/ipfs/go-cid"
	cbg "github.com/whyrusleeping/cbor-gen"
	xerrors "golang.org/x/xerrors"
)

var _ = xerrors.Errorf
var _ = cid.Undef
var _ = math.E
var _ = sort.Sort

func (t *EventHeader) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write([]byte{161}); err != nil {
		return err
	}

	// t.Op (int64) (int64)
	if len("op") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"op\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("op"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("op")); err != nil {
		return err
	}

	if t.Op >= 0 {
		if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.Op)); err != nil {
			return err
		}
	} else {
		if err := cw.WriteMajorTypeHeader(cbg.MajNegativeInt, uint64(-t.Op-1)); err != nil {
			return err
		}
	}
	return nil
}

func (t *EventHeader) UnmarshalCBOR(r io.Reader) (err error) {
	*t = EventHeader{}

	cr := cbg.NewCborReader(r)

	maj, extra, err := cr.ReadHeader()
	if err != nil {
		return err
	}
	defer func() {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
	}()

	if maj != cbg.MajMap {
		return fmt.Errorf("cbor input should be of type map")
	}

	if extra > cbg.MaxLength {
		return fmt.Errorf("EventHeader: map struct too large (%d)", extra)
	}

	var name string
	n := extra

	for i := uint64(0); i < n; i++ {

		{
			sval, err := cbg.ReadString(cr)
			if err != nil {
				return err
			}

			name = string(sval)
		}

		switch name {
		// t.Op (int64) (int64)
		case "op":
			{
				maj, extra, err := cr.ReadHeader()
				var extraI int64
				if err != nil {
					return err
				}
				switch maj {
				case cbg.MajUnsignedInt:
					extraI = int64(extra)
					if extraI < 0 {
						return fmt.Errorf("int64 positive overflow")
					}
				case cbg.MajNegativeInt:
					extraI = int64(extra)
					if extraI < 0 {
						return fmt.Errorf("int64 negative overflow")
					}
					extraI = -1 - extraI
				default:
					return fmt.Errorf("wrong type for int64 field: %d", maj)
				}

				t.Op = int64(extraI)
			}

		default:
			// Field doesn't exist on this type, so ignore it
			cbg.ScanForLinks(r, func(cid.Cid) {})
		}
	}

	return nil
}
func (t *RepoAppend) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write([]byte{169}); err != nil {
		return err
	}

	// t.Ops ([]*events.RepoOp) (slice)
	if len("ops") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"ops\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("ops"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("ops")); err != nil {
		return err
	}

	if len(t.Ops) > cbg.MaxLength {
		return xerrors.Errorf("Slice value in field t.Ops was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajArray, uint64(len(t.Ops))); err != nil {
		return err
	}
	for _, v := range t.Ops {
		if err := v.MarshalCBOR(cw); err != nil {
			return err
		}
	}

	// t.Seq (int64) (int64)
	if len("seq") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"seq\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("seq"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("seq")); err != nil {
		return err
	}

	if t.Seq >= 0 {
		if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.Seq)); err != nil {
			return err
		}
	} else {
		if err := cw.WriteMajorTypeHeader(cbg.MajNegativeInt, uint64(-t.Seq-1)); err != nil {
			return err
		}
	}

	// t.Prev (string) (string)
	if len("prev") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"prev\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("prev"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("prev")); err != nil {
		return err
	}

	if t.Prev == nil {
		if _, err := cw.Write(cbg.CborNull); err != nil {
			return err
		}
	} else {
		if len(*t.Prev) > cbg.MaxLength {
			return xerrors.Errorf("Value in field t.Prev was too long")
		}

		if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len(*t.Prev))); err != nil {
			return err
		}
		if _, err := io.WriteString(w, string(*t.Prev)); err != nil {
			return err
		}
	}

	// t.Repo (string) (string)
	if len("repo") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"repo\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("repo"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("repo")); err != nil {
		return err
	}

	if len(t.Repo) > cbg.MaxLength {
		return xerrors.Errorf("Value in field t.Repo was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len(t.Repo))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string(t.Repo)); err != nil {
		return err
	}

	// t.Time (string) (string)
	if len("time") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"time\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("time"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("time")); err != nil {
		return err
	}

	if len(t.Time) > cbg.MaxLength {
		return xerrors.Errorf("Value in field t.Time was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len(t.Time))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string(t.Time)); err != nil {
		return err
	}

	// t.Blobs ([]string) (slice)
	if len("blobs") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"blobs\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("blobs"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("blobs")); err != nil {
		return err
	}

	if len(t.Blobs) > cbg.MaxLength {
		return xerrors.Errorf("Slice value in field t.Blobs was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajArray, uint64(len(t.Blobs))); err != nil {
		return err
	}
	for _, v := range t.Blobs {
		if len(v) > cbg.MaxLength {
			return xerrors.Errorf("Value in field v was too long")
		}

		if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len(v))); err != nil {
			return err
		}
		if _, err := io.WriteString(w, string(v)); err != nil {
			return err
		}
	}

	// t.Event (string) (string)
	if len("event") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"event\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("event"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("event")); err != nil {
		return err
	}

	if len(t.Event) > cbg.MaxLength {
		return xerrors.Errorf("Value in field t.Event was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len(t.Event))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string(t.Event)); err != nil {
		return err
	}

	// t.Blocks ([]uint8) (slice)
	if len("blocks") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"blocks\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("blocks"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("blocks")); err != nil {
		return err
	}

	if len(t.Blocks) > cbg.ByteArrayMaxLen {
		return xerrors.Errorf("Byte array in field t.Blocks was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajByteString, uint64(len(t.Blocks))); err != nil {
		return err
	}

	if _, err := cw.Write(t.Blocks[:]); err != nil {
		return err
	}

	// t.Commit (string) (string)
	if len("commit") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"commit\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("commit"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("commit")); err != nil {
		return err
	}

	if len(t.Commit) > cbg.MaxLength {
		return xerrors.Errorf("Value in field t.Commit was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len(t.Commit))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string(t.Commit)); err != nil {
		return err
	}
	return nil
}

func (t *RepoAppend) UnmarshalCBOR(r io.Reader) (err error) {
	*t = RepoAppend{}

	cr := cbg.NewCborReader(r)

	maj, extra, err := cr.ReadHeader()
	if err != nil {
		return err
	}
	defer func() {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
	}()

	if maj != cbg.MajMap {
		return fmt.Errorf("cbor input should be of type map")
	}

	if extra > cbg.MaxLength {
		return fmt.Errorf("RepoAppend: map struct too large (%d)", extra)
	}

	var name string
	n := extra

	for i := uint64(0); i < n; i++ {

		{
			sval, err := cbg.ReadString(cr)
			if err != nil {
				return err
			}

			name = string(sval)
		}

		switch name {
		// t.Ops ([]*events.RepoOp) (slice)
		case "ops":

			maj, extra, err = cr.ReadHeader()
			if err != nil {
				return err
			}

			if extra > cbg.MaxLength {
				return fmt.Errorf("t.Ops: array too large (%d)", extra)
			}

			if maj != cbg.MajArray {
				return fmt.Errorf("expected cbor array")
			}

			if extra > 0 {
				t.Ops = make([]*RepoOp, extra)
			}

			for i := 0; i < int(extra); i++ {

				var v RepoOp
				if err := v.UnmarshalCBOR(cr); err != nil {
					return err
				}

				t.Ops[i] = &v
			}

			// t.Seq (int64) (int64)
		case "seq":
			{
				maj, extra, err := cr.ReadHeader()
				var extraI int64
				if err != nil {
					return err
				}
				switch maj {
				case cbg.MajUnsignedInt:
					extraI = int64(extra)
					if extraI < 0 {
						return fmt.Errorf("int64 positive overflow")
					}
				case cbg.MajNegativeInt:
					extraI = int64(extra)
					if extraI < 0 {
						return fmt.Errorf("int64 negative overflow")
					}
					extraI = -1 - extraI
				default:
					return fmt.Errorf("wrong type for int64 field: %d", maj)
				}

				t.Seq = int64(extraI)
			}
			// t.Prev (string) (string)
		case "prev":

			{
				b, err := cr.ReadByte()
				if err != nil {
					return err
				}
				if b != cbg.CborNull[0] {
					if err := cr.UnreadByte(); err != nil {
						return err
					}

					sval, err := cbg.ReadString(cr)
					if err != nil {
						return err
					}

					t.Prev = (*string)(&sval)
				}
			}
			// t.Repo (string) (string)
		case "repo":

			{
				sval, err := cbg.ReadString(cr)
				if err != nil {
					return err
				}

				t.Repo = string(sval)
			}
			// t.Time (string) (string)
		case "time":

			{
				sval, err := cbg.ReadString(cr)
				if err != nil {
					return err
				}

				t.Time = string(sval)
			}
			// t.Blobs ([]string) (slice)
		case "blobs":

			maj, extra, err = cr.ReadHeader()
			if err != nil {
				return err
			}

			if extra > cbg.MaxLength {
				return fmt.Errorf("t.Blobs: array too large (%d)", extra)
			}

			if maj != cbg.MajArray {
				return fmt.Errorf("expected cbor array")
			}

			if extra > 0 {
				t.Blobs = make([]string, extra)
			}

			for i := 0; i < int(extra); i++ {

				{
					sval, err := cbg.ReadString(cr)
					if err != nil {
						return err
					}

					t.Blobs[i] = string(sval)
				}
			}

			// t.Event (string) (string)
		case "event":

			{
				sval, err := cbg.ReadString(cr)
				if err != nil {
					return err
				}

				t.Event = string(sval)
			}
			// t.Blocks ([]uint8) (slice)
		case "blocks":

			maj, extra, err = cr.ReadHeader()
			if err != nil {
				return err
			}

			if extra > cbg.ByteArrayMaxLen {
				return fmt.Errorf("t.Blocks: byte array too large (%d)", extra)
			}
			if maj != cbg.MajByteString {
				return fmt.Errorf("expected byte array")
			}

			if extra > 0 {
				t.Blocks = make([]uint8, extra)
			}

			if _, err := io.ReadFull(cr, t.Blocks[:]); err != nil {
				return err
			}
			// t.Commit (string) (string)
		case "commit":

			{
				sval, err := cbg.ReadString(cr)
				if err != nil {
					return err
				}

				t.Commit = string(sval)
			}

		default:
			// Field doesn't exist on this type, so ignore it
			cbg.ScanForLinks(r, func(cid.Cid) {})
		}
	}

	return nil
}
func (t *RepoOp) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write([]byte{163}); err != nil {
		return err
	}

	// t.Cid (cid.Cid) (struct)
	if len("cid") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"cid\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("cid"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("cid")); err != nil {
		return err
	}

	if t.Cid == nil {
		if _, err := cw.Write(cbg.CborNull); err != nil {
			return err
		}
	} else {
		if err := cbg.WriteCid(cw, *t.Cid); err != nil {
			return xerrors.Errorf("failed to write cid field t.Cid: %w", err)
		}
	}

	// t.Path (string) (string)
	if len("path") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"path\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("path"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("path")); err != nil {
		return err
	}

	if len(t.Path) > cbg.MaxLength {
		return xerrors.Errorf("Value in field t.Path was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len(t.Path))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string(t.Path)); err != nil {
		return err
	}

	// t.Action (string) (string)
	if len("action") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"action\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("action"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("action")); err != nil {
		return err
	}

	if len(t.Action) > cbg.MaxLength {
		return xerrors.Errorf("Value in field t.Action was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len(t.Action))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string(t.Action)); err != nil {
		return err
	}
	return nil
}

func (t *RepoOp) UnmarshalCBOR(r io.Reader) (err error) {
	*t = RepoOp{}

	cr := cbg.NewCborReader(r)

	maj, extra, err := cr.ReadHeader()
	if err != nil {
		return err
	}
	defer func() {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
	}()

	if maj != cbg.MajMap {
		return fmt.Errorf("cbor input should be of type map")
	}

	if extra > cbg.MaxLength {
		return fmt.Errorf("RepoOp: map struct too large (%d)", extra)
	}

	var name string
	n := extra

	for i := uint64(0); i < n; i++ {

		{
			sval, err := cbg.ReadString(cr)
			if err != nil {
				return err
			}

			name = string(sval)
		}

		switch name {
		// t.Cid (cid.Cid) (struct)
		case "cid":

			{

				b, err := cr.ReadByte()
				if err != nil {
					return err
				}
				if b != cbg.CborNull[0] {
					if err := cr.UnreadByte(); err != nil {
						return err
					}

					c, err := cbg.ReadCid(cr)
					if err != nil {
						return xerrors.Errorf("failed to read cid field t.Cid: %w", err)
					}

					t.Cid = &c
				}

			}
			// t.Path (string) (string)
		case "path":

			{
				sval, err := cbg.ReadString(cr)
				if err != nil {
					return err
				}

				t.Path = string(sval)
			}
			// t.Action (string) (string)
		case "action":

			{
				sval, err := cbg.ReadString(cr)
				if err != nil {
					return err
				}

				t.Action = string(sval)
			}

		default:
			// Field doesn't exist on this type, so ignore it
			cbg.ScanForLinks(r, func(cid.Cid) {})
		}
	}

	return nil
}
func (t *InfoFrame) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write([]byte{162}); err != nil {
		return err
	}

	// t.Info (string) (string)
	if len("info") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"info\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("info"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("info")); err != nil {
		return err
	}

	if len(t.Info) > cbg.MaxLength {
		return xerrors.Errorf("Value in field t.Info was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len(t.Info))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string(t.Info)); err != nil {
		return err
	}

	// t.Message (string) (string)
	if len("message") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"message\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("message"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("message")); err != nil {
		return err
	}

	if len(t.Message) > cbg.MaxLength {
		return xerrors.Errorf("Value in field t.Message was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len(t.Message))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string(t.Message)); err != nil {
		return err
	}
	return nil
}

func (t *InfoFrame) UnmarshalCBOR(r io.Reader) (err error) {
	*t = InfoFrame{}

	cr := cbg.NewCborReader(r)

	maj, extra, err := cr.ReadHeader()
	if err != nil {
		return err
	}
	defer func() {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
	}()

	if maj != cbg.MajMap {
		return fmt.Errorf("cbor input should be of type map")
	}

	if extra > cbg.MaxLength {
		return fmt.Errorf("InfoFrame: map struct too large (%d)", extra)
	}

	var name string
	n := extra

	for i := uint64(0); i < n; i++ {

		{
			sval, err := cbg.ReadString(cr)
			if err != nil {
				return err
			}

			name = string(sval)
		}

		switch name {
		// t.Info (string) (string)
		case "info":

			{
				sval, err := cbg.ReadString(cr)
				if err != nil {
					return err
				}

				t.Info = string(sval)
			}
			// t.Message (string) (string)
		case "message":

			{
				sval, err := cbg.ReadString(cr)
				if err != nil {
					return err
				}

				t.Message = string(sval)
			}

		default:
			// Field doesn't exist on this type, so ignore it
			cbg.ScanForLinks(r, func(cid.Cid) {})
		}
	}

	return nil
}
func (t *ErrorFrame) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write([]byte{162}); err != nil {
		return err
	}

	// t.Error (string) (string)
	if len("error") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"error\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("error"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("error")); err != nil {
		return err
	}

	if len(t.Error) > cbg.MaxLength {
		return xerrors.Errorf("Value in field t.Error was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len(t.Error))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string(t.Error)); err != nil {
		return err
	}

	// t.Message (string) (string)
	if len("message") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"message\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("message"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("message")); err != nil {
		return err
	}

	if len(t.Message) > cbg.MaxLength {
		return xerrors.Errorf("Value in field t.Message was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len(t.Message))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string(t.Message)); err != nil {
		return err
	}
	return nil
}

func (t *ErrorFrame) UnmarshalCBOR(r io.Reader) (err error) {
	*t = ErrorFrame{}

	cr := cbg.NewCborReader(r)

	maj, extra, err := cr.ReadHeader()
	if err != nil {
		return err
	}
	defer func() {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
	}()

	if maj != cbg.MajMap {
		return fmt.Errorf("cbor input should be of type map")
	}

	if extra > cbg.MaxLength {
		return fmt.Errorf("ErrorFrame: map struct too large (%d)", extra)
	}

	var name string
	n := extra

	for i := uint64(0); i < n; i++ {

		{
			sval, err := cbg.ReadString(cr)
			if err != nil {
				return err
			}

			name = string(sval)
		}

		switch name {
		// t.Error (string) (string)
		case "error":

			{
				sval, err := cbg.ReadString(cr)
				if err != nil {
					return err
				}

				t.Error = string(sval)
			}
			// t.Message (string) (string)
		case "message":

			{
				sval, err := cbg.ReadString(cr)
				if err != nil {
					return err
				}

				t.Message = string(sval)
			}

		default:
			// Field doesn't exist on this type, so ignore it
			cbg.ScanForLinks(r, func(cid.Cid) {})
		}
	}

	return nil
}
