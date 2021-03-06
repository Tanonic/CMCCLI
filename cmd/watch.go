// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
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
	"github.com/Jeffail/gabs"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strings"
)

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "list your watched currencies",
	Long: `With just "list" you can list your saved currencies. With "list add [symbol of coin e.g. BTC]". With "list rm [BTC] you can remove it again from your watchlist"`,
	TraverseChildren: true,
	Run: func(cmd *cobra.Command, args []string) {
		listWatched()
	},
}

var cmdAdd = &cobra.Command{
	Use:   "add [BTC]",
	Short: "add a currency to the watchlist",
	Long: `Add a currency to your personal watchlist by supplying the symbol of the currency, e.g. BTC.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		symbol := strings.Join(args, "")
		addWatch(symbol)
		fmt.Println("added " + symbol + " to your watchlist")
	},
}

var cmdRmWatch = &cobra.Command{
	Use:   "rm",
	Short: "watch rm [symbol]",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("please add the currency you want to remove, example: btc")
			os.Exit(-3)
		}
		removeWatched(strings.Join(args, ""))
	},
}

func listWatched(){
	if _, err := os.Stat("watchlist.json"); !os.IsNotExist(err) { //if file already exists
		//parse watchlist.json file with gabs
		jsonFile, err := gabs.ParseJSONFile("watchlist.json")
		if err != nil {
			log.Fatal(err)
		}

		var sym []string
		slice := sym
		//append all the symbol names to a string array for the api request
		children, _ := jsonFile.S("symbol").Children()
		for _, child := range children {
			slice = append(slice, strings.ToUpper(child.Data().(string))) //make all symbols uppercase because the api requires the symbols to be uppercase
		}
		GetCoinData(slice)
	} else { //if file doesn't exist yet
		fmt.Println("You first have to add something to your watchlist before listing. Use \"watch add [BTC] to add a currency to your watchlist\"")
	}
}

func removeWatched(symbol string) {
	if _, err := os.Stat("watchlist.json"); !os.IsNotExist(err) { //if file already exists
		//parse watchlist.json file with gabs
		jsonFile, err := gabs.ParseJSONFile("watchlist.json")
		if err != nil {
			log.Fatal(err)
		}

		var index int
		children,_ := jsonFile.S("symbol").Children()
		for _, child := range children {
			if child.Data().(string) == symbol {
				jsonFile.ArrayRemove(index, "symbol")
				writeWatchFile(jsonFile)
				fmt.Println("removed " + symbol + " from your watchlist")
				os.Exit(-3)
			} else {
				index++
			}
		}
		fmt.Println(symbol + " is not in your watchlist")
	} else { //if file doesn't exist yet
		fmt.Println("You don't have a watchlist yet, to add something use \"watch add [BTC]\"")
	}
}

func addWatch(symbol string) {
	if _, err := os.Stat("watchlist.json"); !os.IsNotExist(err) { //if file already exists, overwrite it with new old and new data merged
		jsonData, err := gabs.ParseJSONFile("watchlist.json")
		if err != nil {
			log.Fatal(err)
		}
		jsonData.ArrayAppend(symbol, "symbol")
		writeWatchFile(jsonData)
	} else { //if file doesn't exist yet, create it and write to it
		jsonObj := gabs.New()
		jsonObj.Array("symbol")
		jsonObj.ArrayAppend(symbol ,"symbol")
		writeWatchFile(jsonObj)
	}
}

func writeWatchFile(jsonData *gabs.Container) {
	f, err := os.Create("watchlist.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	fmt.Fprintln(f, jsonData.StringIndent("", "  "))
}

func init() {
	rootCmd.AddCommand(watchCmd)
	watchCmd.AddCommand(cmdAdd)
	watchCmd.AddCommand(cmdRmWatch)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// portCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// portCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
