package types

import (
	"encoding/hex"
	"fmt"
	"io"
	"time"

	"kuskcore/encoding/blockchain"
	"kuskcore/encoding/bufpool"
	"kuskcore/errors"
	"kuskcore/protocol/bc"
)

// BlockHeader defines information about a block and is used in the Kusk
type BlockHeader struct {
	Version           uint64  // The version of the block.
	Height            uint64  // The height of the block.
	PreviousBlockHash bc.Hash // The hash of the previous block.
	Timestamp         uint64  // The time of the block in seconds.
	BlockWitness
	SupLinks
	BlockCommitment
}

// Hash returns complete hash of the block header.
func (bh *BlockHeader) Hash() bc.Hash {
	h, _ := mapBlockHeader(bh)
	return h
}

// Time returns the time represented by the Timestamp in block header.
func (bh *BlockHeader) Time() time.Time {
	return time.Unix(int64(bh.Timestamp/1000), 0).UTC()
}

// MarshalText fulfills the json.Marshaler interface. This guarantees that
// block headers will get deserialized correctly when being parsed from HTTP
// requests.
func (bh *BlockHeader) MarshalText() ([]byte, error) {
	buf := bufpool.Get()
	defer bufpool.Put(buf)

	if _, err := bh.WriteTo(buf); err != nil {
		return nil, err
	}

	enc := make([]byte, hex.EncodedLen(buf.Len()))
	hex.Encode(enc, buf.Bytes())
	return enc, nil
}

// UnmarshalText fulfills the encoding.TextUnmarshaler interface.
func (bh *BlockHeader) UnmarshalText(text []byte) error {
	decoded := make([]byte, hex.DecodedLen(len(text)))
	if _, err := hex.Decode(decoded, text); err != nil {
		return err
	}

	serflag, err := bh.readFrom(blockchain.NewReader(decoded))
	if err != nil {
		return err
	}

	if serflag == SerBlockTransactions {
		return fmt.Errorf("unsupported serialization flags 0x%02x", serflag)
	}

	return nil
}

// WriteTo writes the block header to the input io.Writer
func (bh *BlockHeader) WriteTo(w io.Writer) (int64, error) {
	ew := errors.NewWriter(w)
	if err := bh.writeTo(ew, SerBlockHeader); err != nil {
		return 0, err
	}
	return ew.Written(), ew.Err()
}

func (bh *BlockHeader) readFrom(r *blockchain.Reader) (serflag uint8, err error) {
	var serflags [1]byte
	if _, err := io.ReadFull(r, serflags[:]); err != nil {
		return 0, err
	}

	serflag = serflags[0]
	switch serflag {
	case SerBlockHeader, SerBlockFull:
	case SerBlockTransactions:
		return
	default:
		return 0, fmt.Errorf("unsupported serialization flags 0x%x", serflags)
	}

	if bh.Version, err = blockchain.ReadVarint63(r); err != nil {
		return 0, err
	}

	if bh.Height, err = blockchain.ReadVarint63(r); err != nil {
		return 0, err
	}

	if _, err = bh.PreviousBlockHash.ReadFrom(r); err != nil {
		return 0, err
	}

	if bh.Timestamp, err = blockchain.ReadVarint63(r); err != nil {
		return 0, err
	}

	if _, err = blockchain.ReadExtensibleString(r, bh.BlockCommitment.readFrom); err != nil {
		return 0, err
	}

	if _, err = blockchain.ReadExtensibleString(r, bh.BlockWitness.readFrom); err != nil {
		return 0, err
	}

	if _, err = blockchain.ReadExtensibleString(r, bh.SupLinks.readFrom); err != nil {
		return 0, err
	}

	return
}

func (bh *BlockHeader) writeTo(w io.Writer, serflags uint8) (err error) {
	w.Write([]byte{serflags})
	if serflags == SerBlockTransactions {
		return nil
	}

	if _, err = blockchain.WriteVarint63(w, bh.Version); err != nil {
		return err
	}

	if _, err = blockchain.WriteVarint63(w, bh.Height); err != nil {
		return err
	}

	if _, err = bh.PreviousBlockHash.WriteTo(w); err != nil {
		return err
	}

	if _, err = blockchain.WriteVarint63(w, bh.Timestamp); err != nil {
		return err
	}

	if _, err = blockchain.WriteExtensibleString(w, nil, bh.BlockCommitment.writeTo); err != nil {
		return err
	}

	if _, err = blockchain.WriteExtensibleString(w, nil, bh.BlockWitness.writeTo); err != nil {
		return err
	}

	if _, err = blockchain.WriteExtensibleString(w, nil, bh.SupLinks.writeTo); err != nil {
		return err
	}

	return
}
