// Code generated by github.com/whyrusleeping/cbor-gen. DO NOT EDIT.

package repo

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

func (t *SignedRoot) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write([]byte{162}); err != nil {
		return err
	}

	// t.Root (cid.Cid) (struct)
	if len("root") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"root\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("root"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("root")); err != nil {
		return err
	}

	if err := cbg.WriteCid(cw, t.Root); err != nil {
		return xerrors.Errorf("failed to write cid field t.Root: %w", err)
	}

	// t.Sig ([]uint8) (slice)
	if len("sig") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"sig\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("sig"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("sig")); err != nil {
		return err
	}

	if len(t.Sig) > cbg.ByteArrayMaxLen {
		return xerrors.Errorf("Byte array in field t.Sig was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajByteString, uint64(len(t.Sig))); err != nil {
		return err
	}

	if _, err := cw.Write(t.Sig[:]); err != nil {
		return err
	}
	return nil
}

func (t *SignedRoot) UnmarshalCBOR(r io.Reader) (err error) {
	*t = SignedRoot{}

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
		return fmt.Errorf("SignedRoot: map struct too large (%d)", extra)
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
		// t.Root (cid.Cid) (struct)
		case "root":

			{

				c, err := cbg.ReadCid(cr)
				if err != nil {
					return xerrors.Errorf("failed to read cid field t.Root: %w", err)
				}

				t.Root = c

			}
			// t.Sig ([]uint8) (slice)
		case "sig":

			maj, extra, err = cr.ReadHeader()
			if err != nil {
				return err
			}

			if extra > cbg.ByteArrayMaxLen {
				return fmt.Errorf("t.Sig: byte array too large (%d)", extra)
			}
			if maj != cbg.MajByteString {
				return fmt.Errorf("expected byte array")
			}

			if extra > 0 {
				t.Sig = make([]uint8, extra)
			}

			if _, err := io.ReadFull(cr, t.Sig[:]); err != nil {
				return err
			}

		default:
			// Field doesn't exist on this type, so ignore it
			cbg.ScanForLinks(r, func(cid.Cid) {})
		}
	}

	return nil
}
func (t *Meta) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write([]byte{163}); err != nil {
		return err
	}

	// t.Datastore (string) (string)
	if len("datastore") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"datastore\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("datastore"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("datastore")); err != nil {
		return err
	}

	if len(t.Datastore) > cbg.MaxLength {
		return xerrors.Errorf("Value in field t.Datastore was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len(t.Datastore))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string(t.Datastore)); err != nil {
		return err
	}

	// t.Did (string) (string)
	if len("did") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"did\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("did"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("did")); err != nil {
		return err
	}

	if len(t.Did) > cbg.MaxLength {
		return xerrors.Errorf("Value in field t.Did was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len(t.Did))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string(t.Did)); err != nil {
		return err
	}

	// t.Version (int64) (int64)
	if len("version") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"version\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("version"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("version")); err != nil {
		return err
	}

	if t.Version >= 0 {
		if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.Version)); err != nil {
			return err
		}
	} else {
		if err := cw.WriteMajorTypeHeader(cbg.MajNegativeInt, uint64(-t.Version-1)); err != nil {
			return err
		}
	}
	return nil
}

func (t *Meta) UnmarshalCBOR(r io.Reader) (err error) {
	*t = Meta{}

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
		return fmt.Errorf("Meta: map struct too large (%d)", extra)
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
		// t.Datastore (string) (string)
		case "datastore":

			{
				sval, err := cbg.ReadString(cr)
				if err != nil {
					return err
				}

				t.Datastore = string(sval)
			}
			// t.Did (string) (string)
		case "did":

			{
				sval, err := cbg.ReadString(cr)
				if err != nil {
					return err
				}

				t.Did = string(sval)
			}
			// t.Version (int64) (int64)
		case "version":
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
						return fmt.Errorf("int64 negative oveflow")
					}
					extraI = -1 - extraI
				default:
					return fmt.Errorf("wrong type for int64 field: %d", maj)
				}

				t.Version = int64(extraI)
			}

		default:
			// Field doesn't exist on this type, so ignore it
			cbg.ScanForLinks(r, func(cid.Cid) {})
		}
	}

	return nil
}
func (t *Commit) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write([]byte{164}); err != nil {
		return err
	}

	// t.AuthToken (string) (string)
	if len("auth_token") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"auth_token\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("auth_token"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("auth_token")); err != nil {
		return err
	}

	if t.AuthToken == nil {
		if _, err := cw.Write(cbg.CborNull); err != nil {
			return err
		}
	} else {
		if len(*t.AuthToken) > cbg.MaxLength {
			return xerrors.Errorf("Value in field t.AuthToken was too long")
		}

		if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len(*t.AuthToken))); err != nil {
			return err
		}
		if _, err := io.WriteString(w, string(*t.AuthToken)); err != nil {
			return err
		}
	}

	// t.Data (cid.Cid) (struct)
	if len("data") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"data\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("data"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("data")); err != nil {
		return err
	}

	if err := cbg.WriteCid(cw, t.Data); err != nil {
		return xerrors.Errorf("failed to write cid field t.Data: %w", err)
	}

	// t.Meta (cid.Cid) (struct)
	if len("meta") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"meta\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("meta"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("meta")); err != nil {
		return err
	}

	if err := cbg.WriteCid(cw, t.Meta); err != nil {
		return xerrors.Errorf("failed to write cid field t.Meta: %w", err)
	}

	// t.Prev (cid.Cid) (struct)
	if len("prev") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"prev\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("prev"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("prev")); err != nil {
		return err
	}

	if err := cbg.WriteCid(cw, t.Prev); err != nil {
		return xerrors.Errorf("failed to write cid field t.Prev: %w", err)
	}

	return nil
}

func (t *Commit) UnmarshalCBOR(r io.Reader) (err error) {
	*t = Commit{}

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
		return fmt.Errorf("Commit: map struct too large (%d)", extra)
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
		// t.AuthToken (string) (string)
		case "auth_token":

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

					t.AuthToken = (*string)(&sval)
				}
			}
			// t.Data (cid.Cid) (struct)
		case "data":

			{

				c, err := cbg.ReadCid(cr)
				if err != nil {
					return xerrors.Errorf("failed to read cid field t.Data: %w", err)
				}

				t.Data = c

			}
			// t.Meta (cid.Cid) (struct)
		case "meta":

			{

				c, err := cbg.ReadCid(cr)
				if err != nil {
					return xerrors.Errorf("failed to read cid field t.Meta: %w", err)
				}

				t.Meta = c

			}
			// t.Prev (cid.Cid) (struct)
		case "prev":

			{

				c, err := cbg.ReadCid(cr)
				if err != nil {
					return xerrors.Errorf("failed to read cid field t.Prev: %w", err)
				}

				t.Prev = c

			}

		default:
			// Field doesn't exist on this type, so ignore it
			cbg.ScanForLinks(r, func(cid.Cid) {})
		}
	}

	return nil
}
