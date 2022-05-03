// Copyright © 2022 Obol Labs Inc.
//
// This program is free software: you can redistribute it and/or modify it
// under the terms of the GNU General Public License as published by the Free
// Software Foundation, either version 3 of the License, or (at your option)
// any later version.
//
// This program is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of  MERCHANTABILITY or
// FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for
// more details.
//
// You should have received a copy of the GNU General Public License along with
// this program.  If not, see <http://www.gnu.org/licenses/>.

package dkg_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"
	"testing"
	"time"

	"github.com/coinbase/kryptology/pkg/signatures/bls/bls_sig"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"

	"github.com/obolnetwork/charon/app/log"
	"github.com/obolnetwork/charon/cluster"
	"github.com/obolnetwork/charon/cmd"
	"github.com/obolnetwork/charon/dkg"
	"github.com/obolnetwork/charon/p2p"
	"github.com/obolnetwork/charon/tbls"
	"github.com/obolnetwork/charon/testutil"
	"github.com/obolnetwork/charon/testutil/keystore"
)

func TestDKG(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	const (
		nodes = 3
		vals  = 2
	)

	// Start bootnode
	bootnode, errChan := startBootnode(ctx, t)

	// Setup
	lock, keys, _ := cluster.NewForT(t, vals, 3, 3, 0)

	dir, err := os.MkdirTemp("", "")
	require.NoError(t, err)

	conf := dkg.Config{
		DataDir: dir,
		P2P: p2p.Config{
			UDPBootnodes: []string{bootnode},
		},
		Log:     log.DefaultConfig(),
		TestDef: &lock.Definition,
	}

	// Run dkg for each node
	var eg errgroup.Group
	for i := 0; i < nodes; i++ {
		conf := conf
		conf.DataDir = path.Join(dir, fmt.Sprintf("node%d", i))
		conf.P2P.TCPAddrs = []string{testutil.AvailableAddr(t).String()}
		conf.P2P.UDPAddr = testutil.AvailableAddr(t).String()

		require.NoError(t, os.MkdirAll(conf.DataDir, 0o755))
		err := crypto.SaveECDSA(p2p.KeyPath(conf.DataDir), keys[i])
		require.NoError(t, err)

		eg.Go(func() error {
			return dkg.Run(ctx, conf)
		})
		if i == 0 {
			// Allow node0 some time to startup, this just mitigates startup races and backoffs but isn't required.
			time.Sleep(time.Millisecond * 100)
		}
	}

	// Wait until complete

	runChan := make(chan error, 1)
	go func() {
		runChan <- eg.Wait()
	}()

	select {
	case err := <-errChan:
		cancel()
		testutil.SkipIfBindErr(t, err)
		require.NoError(t, err)
	case err := <-runChan:
		cancel()
		testutil.SkipIfBindErr(t, err)
		require.NoError(t, err)
	}

	// Read generated lock and keystores from disk
	var (
		shares = make([][]*bls_sig.SecretKey, vals)
		locks  []cluster.Lock
	)
	for i := 0; i < nodes; i++ {
		dataDir := path.Join(dir, fmt.Sprintf("node%d", i))
		keyShares, err := keystore.LoadKeys(dataDir)
		require.NoError(t, err)
		require.Len(t, keyShares, vals)

		for i, key := range keyShares {
			shares[i] = append(shares[i], key)
		}

		lockFile, err := os.ReadFile(path.Join(dataDir, "cluster_lock.json"))
		require.NoError(t, err)

		var lock cluster.Lock
		require.NoError(t, json.Unmarshal(lockFile, &lock))
		locks = append(locks, lock)
	}

	// Ensure locks hashes are identical.
	var hash []byte
	for i, lock := range locks {
		h, err := lock.HashTreeRoot()
		require.NoError(t, err)
		if i == 0 {
			hash = h[:]
		} else {
			require.Equal(t, hash, h[:])
		}
	}

	// 	Ensure keystores can generate valid tbls aggregate signature.
	for i := 0; i < vals; i++ {
		var sigs []*bls_sig.PartialSignature
		for j := 0; j < nodes; j++ {
			sig, err := tbls.Sign(shares[i][j], []byte("data"))
			require.NoError(t, err)
			sigs = append(sigs, &bls_sig.PartialSignature{
				Identifier: byte(j),
				Signature:  sig.Value,
			})
		}
		_, err := tbls.Aggregate(sigs)
		require.NoError(t, err)
	}
}

// startBootnode starts a charon bootnode and returns its http ENR endpoint.
func startBootnode(ctx context.Context, t *testing.T) (string, <-chan error) {
	t.Helper()

	dir, err := os.MkdirTemp("", "")
	require.NoError(t, err)

	addr := testutil.AvailableAddr(t).String()

	errChan := make(chan error, 1)
	go func() {
		errChan <- cmd.RunBootnode(ctx, cmd.BootnodeConfig{
			DataDir:  dir,
			HTTPAddr: addr,
			P2PConfig: p2p.Config{
				UDPAddr:  testutil.AvailableAddr(t).String(),
				TCPAddrs: []string{testutil.AvailableAddr(t).String()},
			},
			LogConfig: log.Config{
				Level:  "error",
				Format: "console",
			},
			AutoP2PKey: true,
			P2PRelay:   true,
		})
	}()

	endpoint := "http://" + addr + "/enr"

	// Wait for bootnode to become available.
	for ctx.Err() == nil {
		_, err := http.Get(endpoint)
		if err == nil {
			break
		}
		time.Sleep(time.Millisecond * 100)
	}

	return endpoint, errChan
}