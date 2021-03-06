package cmd

import (
	"fmt"
	"net/url"
	"os"
	"strconv"

	"encoding/json"

	"github.com/skuid/skuid/force"
	"github.com/skuid/skuid/types"
	"github.com/spf13/cobra"
)

var outputDir string
var verbose bool
var module string

// pullCmd represents the pull command
var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull Skuid Pages from Salesforce into a local directory.",
	Long:  `Pull your Skuid Pages from your Salesforce org to a local directory.`,
	Run: func(cmd *cobra.Command, args []string) {

		//login to the Force API
		api, err := force.Login(
			appClientID,
			appClientSecret,
			sfdcHost,
			sfdcUsername,
			sfdcPassword,
			apiVersion,
		)

		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		//build the module query paramater
		var requestedModule = []string{""}

		if module != "" {
			requestedModule = []string{module}
		}

		query := url.Values{
			"module": requestedModule,
		}

		//query the API for all pages in the requested module
		result, err := api.Connection.Get("/skuid/api/v1/pages", query)

		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		var pages map[string]types.PullResponse

		//you have to unquote the string because what comes back
		//is escaped json
		unquoted, _ := strconv.Unquote(string(result))

		//unmarshal all pages into the type PullResponse
		err = json.Unmarshal([]byte(unquoted), &pages)

		if err != nil {
			fmt.Println(string(result))
			fmt.Println(err.Error())
			os.Exit(1)
		}

		for _, page := range pages {

			//write the page in the at rest format
			err = page.WriteAtRest(outputDir)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
		}

		fmt.Printf("Pages written to %s\n", outputDir)
	},
}

func init() {
	RootCmd.AddCommand(pullCmd)

	pullCmd.Flags().StringVarP(&outputDir, "dir", "d", "skuidpages", "Output directory for all pages.")
	pullCmd.Flags().StringVarP(&module, "module", "m", "", "Page module(s) to pull. Separated by a comma.")
}
