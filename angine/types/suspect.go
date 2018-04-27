package types

import (
	"bytes"
	"encoding/json"

	"github.com/Baptist-Publication/chorus/types"
)

type SuspectTx struct {
	Suspect   *Hypocrite `json:"suspect"`
	PubKey    []byte     `json:"pubkey"`
	Signature []byte     `json:"signature"`
}

func IsSuspectTx(tx []byte) bool {
	return bytes.Equal(types.TxTagAngineEcoSuspect, tx[:3])
}

func (tx *SuspectTx) ToBytes() ([]byte, error) {
	return json.Marshal(tx)
}

func (tx *SuspectTx) FromBytes(bs []byte) error {
	return json.Unmarshal(bs, tx)
}
