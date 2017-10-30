// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
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
	"fmt"
	"os"

	"github.com/kkirsche/gbust/libgbust"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var c libgbust.Config

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "gbust",
	Short: "gbust is a golang based directory searcher",
	Long: `gbust is a gobuster inspired directory brute forcer. It builds upon
the concepts outlined in gobuster to offer a more robust feature set as
commonly requested by the community.`,
	Run: func(cmd *cobra.Command, args []string) {
		logrus.SetLevel(logrus.InfoLevel)
		if c.Verbose {
			logrus.SetLevel(logrus.DebugLevel)
		}

		attacker, err := libgbust.NewAttacker(&c)
		if err != nil {
			logrus.WithError(err).Errorln("[!] failed to create attacker")
			return
		}
		attacker.Attack()
		attacker.Wg.Wait()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	RootCmd.Flags().IntVarP(&c.Goroutines, "goroutines", "g", 50, "goroutines to work in")
	RootCmd.Flags().StringSliceVarP(&c.Cookies, "cookies", "c", []string{}, "cookies to use for the connections")
	RootCmd.Flags().StringSliceVarP(&c.Wordlists, "wordlists", "w", []string{}, "wordlists to leverage")
	RootCmd.Flags().StringVarP(&c.RawURL, "url", "u", "", "url to brute force")
	RootCmd.Flags().BoolVarP(&c.Verbose, "verbose", "v", false, "enable verbose logging")
}
