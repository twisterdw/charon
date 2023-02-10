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

// Command compose provides a tool to run, test, debug local charon clusters
// using docker-compose.
//
//	It consists of multiple steps:
//	 - compose new: Creates a new config.json that defines what will be composed.
//	 - compose define: Creates a docker-compose.yml that executes `charon create dkg` if keygen==dkg.
//	 - compose lock: Creates a docker-compose.yml that executes `charon create cluster` or `charon dkg`.
//	 - compose run: Creates a docker-compose.yml that executes `charon run`.
package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/obolnetwork/charon/app/log"
	"github.com/obolnetwork/charon/cluster"
	"github.com/obolnetwork/charon/testutil/compose"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cobra.CheckErr(newRootCmd().ExecuteContext(ctx))
}

func newRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "compose",
		Short: "Charon Compose - Run, test, and debug a developer-focussed insecure local charon cluster using docker-compose",
	}

	root.AddCommand(newNewCmd())
	root.AddCommand(newCleanCmd())
	root.AddCommand(newAutoCmd())
	root.AddCommand(newDockerCmd(
		"define",
		"Creates a docker-compose.yml that executes `charon create dkg` if keygen==dkg",
		compose.Define,
	))
	root.AddCommand(newDockerCmd(
		"lock",
		"Creates a docker-compose.yml that executes `charon create cluster` or `charon dkg`",
		compose.Lock,
	))
	root.AddCommand(newDockerCmd(
		"run",
		"Creates a docker-compose.yml that executes `charon run`",
		compose.Run,
	))

	return root
}

// newDockerCmd returns a cobra command that generates docker-compose.yml files and executes it.
func newDockerCmd(use string, short string, runFunc compose.RunFunc) *cobra.Command {
	cmd := &cobra.Command{
		Use:   use,
		Short: short,
	}

	up := addUpFlag(cmd.Flags())
	dir := addDirFlag(cmd.Flags())
	cmd.RunE = func(cmd *cobra.Command, _ []string) error {
		_, err := compose.NewRunnerFunc(use, *dir, *up, runFunc)(cmd.Context())
		if err != nil {
			log.Error(cmd.Context(), "Fatal error", err)
		}

		return err
	}

	return cmd
}

func newAutoCmd() *cobra.Command {
	var conf compose.AutoConfig

	cmd := &cobra.Command{
		Use:   "auto",
		Short: "Convenience function that runs `compose define && compose lock && compose run`",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := compose.Auto(cmd.Context(), conf)
			if err != nil {
				log.Error(cmd.Context(), "auto command fatal error", err)
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&conf.Dir, "compose-dir", ".", "Directory to use for compose artifacts")
	cmd.Flags().DurationVar(&conf.AlertTimeout, "alert-timeout", 0, "Timeout to collect alerts before shutdown. Zero disables timeout.")
	cmd.Flags().BoolVar(&conf.SudoPerms, "sudo-perms", false, "Enables changing all compose artefacts file permissions using sudo.")
	cmd.Flags().BoolVar(&conf.PrintYML, "print-yml", false, "Print generated docker-compose.yml files.")

	return cmd
}

func newNewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "new",
		Short: "Creates a new config.json that defines what will be composed",
	}

	conf := compose.NewDefaultConfig()

	dir := addDirFlag(cmd.Flags())
	keygen := cmd.Flags().String("keygen", string(conf.KeyGen), "Key generation process: create, split, dkg")
	buildLocal := cmd.Flags().Bool("build-local", conf.BuildLocal, "Enables building a local charon container from source. Note this requires the CHARON_REPO env var.")
	beaconNode := cmd.Flags().String("beacon-node", conf.BeaconNode, "Beacon node URL endpoint or 'mock' for simnet.")
	extRelay := cmd.Flags().String("external-relay", "", "Optional external relay HTTP url.")
	splitKeys := cmd.Flags().String("split-keys-dir", conf.SplitKeysDir, "Directory containing keys to split for keygen==create, or empty not to split.")
	featureSet := cmd.Flags().String("feature-set", conf.FeatureSet, "Minimum feature set to enable: alpha, beta, stable")
	numVals := cmd.Flags().Int("num-validators", conf.NumValidators, "Number of distributed validators.")
	vcTypes := cmd.Flags().StringSlice("validator-types", conf.VCStrings(), "Validator types to include.")
	nodes := cmd.Flags().Int("nodes", conf.NumNodes, "Number of charon nodes in the cluster.")
	insecureKeys := cmd.Flags().Bool("insecure-keys", conf.InsecureKeys, "To generate keys quickly.")

	cmd.RunE = func(cmd *cobra.Command, _ []string) error {
		conf.KeyGen = compose.KeyGen(*keygen)
		conf.BuildLocal = *buildLocal
		conf.BeaconNode = *beaconNode
		conf.SplitKeysDir = *splitKeys
		conf.FeatureSet = *featureSet
		conf.ExternalRelay = *extRelay
		conf.NumValidators = *numVals
		conf.NumNodes = *nodes
		conf.Threshold = cluster.Threshold(conf.NumNodes)
		conf.InsecureKeys = *insecureKeys

		if conf.BuildLocal {
			conf.ImageTag = "local"
		}

		var vcs []compose.VCType
		for _, vc := range *vcTypes {
			vcs = append(vcs, compose.VCType(vc))
		}
		conf.VCs = vcs

		ctx := log.WithTopic(cmd.Context(), "new")
		if err := compose.New(ctx, *dir, conf); err != nil {
			log.Error(ctx, "Fatal error", err)
			return err
		}

		return nil
	}

	return cmd
}

func newCleanCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clean",
		Short: "Convenience function that cleans the compose directory",
	}

	dir := addDirFlag(cmd.Flags())

	cmd.RunE = func(cmd *cobra.Command, _ []string) error {
		return compose.Clean(cmd.Context(), *dir)
	}

	return cmd
}

func addDirFlag(flags *pflag.FlagSet) *string {
	return flags.String("compose-dir", ".", "Directory to use for compose artifacts")
}

func addUpFlag(flags *pflag.FlagSet) *bool {
	return flags.Bool("up", true, "Execute `docker-compose up` when compose command completes")
}
