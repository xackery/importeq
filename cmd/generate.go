// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
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
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/xackery/importeq/model"
	"github.com/xackery/importeq/takeadump"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "importeq generate <outfile>",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := generateRun(cmd, args)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	generateCmd.Flags().StringP("dump", "d", ".", "takeadump csv directory")
	generateCmd.Flags().StringP("zone", "z", "", "zone shortname to parse")
}

func generateRun(cmd *cobra.Command, args []string) (err error) {
	start := time.Now()
	var newNpcs map[string]*model.NpcType
	npcs := map[string]*model.NpcType{}
	if len(args) < 1 {
		err = fmt.Errorf("must specify an outfile")
		return
	}

	outFile := args[0]
	fmt.Println("outfile:", outFile)
	fi, err := os.Stat(outFile)
	if err != nil {
		if !os.IsNotExist(err) {
			err = errors.Wrapf(err, "failed to stat %s", outFile)
			return
		}
	}
	if err == nil && fi.IsDir() {
		err = fmt.Errorf("%s is a directory", outFile)
		return
	}

	data, err := cmd.Flags().GetString("dump")

	if data != "" {
		newNpcs, err = generateTakeADump(cmd)
		if err != nil {
			err = errors.Wrap(err, "dump")
			return
		}
		for _, npc := range newNpcs {
			if npcs[npc.Name] == nil {
				npcs[npc.Name] = npc
			}
		}
	}

	outTxt := ""
	//outTxt += fmt.Sprintf("SELECT @npc_type_id := MAX(id) + 1 FROM npc_types WHERE id > %d*1000-1 AND id < (%d+1)*1000;\n", zoneID, zoneID)
	out := ""
	if len(npcs) == 0 {
		err = fmt.Errorf("no npcs were generated")
		return
	}

	for _, npc := range npcs {
		out, err = npc.GenerateInsert()
		if err != nil {
			return
		}
		outTxt += out
	}

	err = ioutil.WriteFile(outFile, []byte(outTxt), 0755)
	if err != nil {
		err = errors.Wrap(err, "failed to save data")
		return
	}
	fmt.Println("----Results----")
	fmt.Printf("%d total npcs extracted\n", len(npcs))
	fmt.Printf("completed in %0.2f seconds\n", time.Since(start).Seconds())
	return
}

func generateTakeADump(cmd *cobra.Command) (npcs map[string]*model.NpcType, err error) {
	filename, err := cmd.Flags().GetString("dump")
	if err != nil {
		err = errors.Wrap(err, "please specify dump dir with e.g -d=data/")
		return
	}

	zoneShortName, err := cmd.Flags().GetString("zone")
	if err != nil {
		err = errors.Wrap(err, "please specify a zone short name with e.g -z=ecommons")
		return
	}

	c, err := takeadump.New(filename, strings.ToLower(zoneShortName))
	if err != nil {
		return
	}
	npcs = c.Npcs
	return
}
