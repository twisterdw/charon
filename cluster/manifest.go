// Copyright © 2021 Obol Technologies Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cluster

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/drand/kyber"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/p2p/enr"
	"github.com/ethereum/go-ethereum/rlp"
	libp2pcrypto "github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiformats/go-multiaddr"

	"github.com/obolnetwork/charon/app/errors"
	"github.com/obolnetwork/charon/crypto"
)

// Manifest captures the public cryptographic and networking info required to connect to a DV cluster.
type Manifest struct {
	TSS     crypto.TBLSScheme     `json:"tss"`     // Threshold signature scheme params
	Members []crypto.BLSPubkeyHex `json:"members"` // DV consensus BLS pubkeys
	ENRs    []string              `json:"enrs"`    // Charon peer ENRs
}

// Pubkey returns the BLS public key of the distributed validator.
func (m *Manifest) Pubkey() kyber.Point {
	return m.TSS.Pubkey()
}

// ParsedENRs returns the decoded list of ENRs in a manifest.
func (m *Manifest) ParsedENRs() ([]enr.Record, error) {
	records := make([]enr.Record, 0, len(m.ENRs))

	for _, enr := range m.ENRs {
		record, err := DecodeENR(enr)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	return records, nil
}

// PeerIDs maps ENRs to libp2p peer IDs.
//
// TODO This can be computed at deserialization-time.
func (m *Manifest) PeerIDs() ([]peer.ID, error) {
	records, err := m.ParsedENRs()
	if err != nil {
		return nil, err
	}

	ids := make([]peer.ID, 0, len(records))

	for _, record := range records {

		info, err := PeerInfoFromENR(record)
		if err != nil {
			return nil, err
		}

		ids = append(ids, info.ID)
	}

	return ids, nil
}

// PeerInfoFromENR derives the libp2p peer ID from the secp256k1 public key encoded in the ENR.
func PeerInfoFromENR(record enr.Record) (peer.AddrInfo, error) {
	var pubkey enode.Secp256k1
	if err := record.Load(&pubkey); err != nil {
		return peer.AddrInfo{}, errors.Wrap(err, "pubkey from enr")
	}

	p2pPubkey := libp2pcrypto.Secp256k1PublicKey(pubkey)
	id, err := peer.IDFromPublicKey(&p2pPubkey)
	if err != nil {
		return peer.AddrInfo{}, errors.Wrap(err, "enr id from pubkey")
	}

	var ip enr.IPv4
	if err := record.Load(&ip); err != nil {
		return peer.AddrInfo{}, errors.Wrap(err, "ip from enr")
	}

	var port enr.TCP
	if err := record.Load(&port); err != nil {
		return peer.AddrInfo{}, errors.Wrap(err, "port from enr")
	}

	mAddrStr := fmt.Sprintf("/ip4/%s/tcp/%d", net.IP(ip).String(), port)
	addr, err := multiaddr.NewMultiaddr(mAddrStr)
	if err != nil {
		return peer.AddrInfo{}, err
	}

	return peer.AddrInfo{
		ID:    id,
		Addrs: []multiaddr.Multiaddr{addr},
	}, nil
}

func EncodeENR(record enr.Record) (string, error) {
	var buf bytes.Buffer
	if err := record.EncodeRLP(&buf); err != nil {
		return "", err
	}

	return "enr:" + base64.URLEncoding.EncodeToString(buf.Bytes()), nil
}

func DecodeENR(enrStr string) (enr.Record, error) {
	enrStr = strings.TrimPrefix(enrStr, "enr:")
	enrBytes, err := base64.URLEncoding.DecodeString(enrStr)
	if err != nil {
		return enr.Record{}, errors.Wrap(err, "base64 enr")
	}

	// TODO support hex encoding too?
	var record enr.Record

	rd := bytes.NewReader(enrBytes)
	if err := rlp.Decode(rd, &record); err != nil {
		return enr.Record{}, errors.Wrap(err, "rlp enr")
	}

	if rd.Len() > 0 {
		return enr.Record{}, errors.New("leftover garbage bytes in ENR")
	}

	return record, nil
}

// LoadManifest reads the manifest from the given file path.
func LoadManifest(file string) (Manifest, error) {
	buf, err := os.ReadFile(file)
	if err != nil {
		return Manifest{}, err
	}

	var res Manifest
	err = json.Unmarshal(buf, &res)
	if err != nil {
		return Manifest{}, errors.Wrap(err, "unmarshal manifest")
	}

	return res, nil
}