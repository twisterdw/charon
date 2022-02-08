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

package cmd

import (
	"context"
	"os"
	"os/signal"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/obolnetwork/charon/app"
	"github.com/obolnetwork/charon/p2p"
)

func newRunCmd(runFunc func(context.Context, app.Config) error) *cobra.Command {
	var conf app.Config

	cmd := &cobra.Command{
		Use:   "run",
		Short: "Runs the Charon middleware",
		Long:  "Starts the long-running Charon middleware process to perform distributed validator duties.",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := signal.NotifyContext(cmd.Context(), os.Interrupt)
			defer cancel()

			return runFunc(ctx, conf)
		},
	}

	bindRunFlags(cmd.Flags(), &conf)
	bindGeneralFlags(cmd.Flags(), &conf.DataDir)
	bindP2PFlags(cmd.Flags(), &conf.P2P)

	return cmd
}

func bindRunFlags(flags *pflag.FlagSet, config *app.Config) {
	flags.StringVar(&config.ManifestFile, "manifest-file", "./charon/manifest.json", "The path to the manifest file defining distributed validator cluster")
	flags.StringVar(&config.BeaconNodeAddr, "beacon-node-endpoint", "http://localhost/", "Beacon node endpoint URL")
	flags.StringVar(&config.ValidatorAPIAddr, "validator-api-address", "127.0.0.1:3500", "Listening address (ip and port) for validator-facing traffic proxying the beacon-node API")
	flags.StringVar(&config.MonitoringAddr, "monitoring-address", "127.0.0.1:8088", "Listening address (ip and port) for the monitoring API (prometheus, pprof)")
	flags.StringVar(&config.JaegerAddr, "jaegar-address", "", "Listening address for Jaegar tracing")
}

func bindGeneralFlags(flags *pflag.FlagSet, dataDir *string) {
	flags.StringVar(dataDir, "data-dir", "./charon/data", "The directory where charon will store all its internal data")
}

func bindP2PFlags(flags *pflag.FlagSet, config *p2p.Config) {
	flags.StringVar(&config.DBPath, "p2p-udp-peerdb", "", "Path to store a discv5 peer database. Empty default results in in-memory database.")
	flags.StringSliceVar(&config.UDPBootnodes, "p2p-udp-bootnodes", nil, "Comma-separated list of discv5 bootnode URLs or ENRs. Example URL: enode://<hex node id>@10.3.58.6:30303?discport=30301")
	flags.StringVar(&config.UDPAddr, "p2p-udp-address", "127.0.0.1:30309", "Listening UDP address (ip and port) for Discv5 discovery")
	flags.StringSliceVar(&config.TCPAddrs, "p2p-tcp-address", []string{"127.0.0.1:13900"}, "Comma-separated list of listening TCP addresses (ip and port) for LibP2P traffic")
	flags.StringVar(&config.Allowlist, "p2p-allowlist", "", "Comma-separated list of CIDR subnets for allowing only certain peer connections. Example: 192.168.0.0/16 would permit connections to peers on your local network only. The default is to accept all connections.")
	flags.StringVar(&config.Denylist, "p2p-denylist", "", "Comma-separated list of CIDR subnets for disallowing certain peer connections. Example: 192.168.0.0/16 would disallow connections to peers on your local network. The default is to accept all connections.")
}