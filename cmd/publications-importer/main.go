package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/Tarick/naca-publications/internal/application/importer"
	"github.com/Tarick/naca-publications/internal/version"
	"github.com/Tarick/naca-publications/pkg/apiclient"

	"github.com/spf13/cobra"
)

func main() {
	var publicationsAPIURL string
	// rootCmd represents the base command when called without any subcommands
	rootCmd := &cobra.Command{
		Use:     "publications-importer",
		Short:   "Publications importer",
		Long:    `Publication importer is used to import Publisher and their publications information. Requires running APIs, url to Publications API and accepts JSON filename as parameter.`,
		Example: `publications-importer --url http://publications publications.json`,
		// Positional arg - one filename of the feed with entries
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// Client for RSS Feeds API
			if len(args) > 1 {
				fmt.Println("Number of parameters more than 1, we accept only one filename: ", args)
				os.Exit(1)
			}
			// Open filename and pass to importer
			if fp, err := os.Open(args[0]); err == nil {
				defer fp.Close()
				bytes, _ := ioutil.ReadAll(fp)
				ip := importer.Importer{
					APIClient: apiclient.New(publicationsAPIURL),
				}
				err := ip.RunImport(bytes)
				if err != nil {
					fmt.Println("Error running import: ", err)
					os.Exit(1)
				}
			} else if os.IsNotExist(err) {
				fmt.Printf("Path '%s' does not exist", args[0])
				os.Exit(1)
			} else {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}
	rootCmd.Flags().StringVar(&publicationsAPIURL, "url", "", "base URL to publications api, e.g. http://publication-api:8080")
	rootCmd.MarkFlagRequired("url")

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version number of application",
		Long:  `Software version`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("NACA Publications importer version:", version.Version, ",build on:", version.BuildTime)
		},
	}
	rootCmd.AddCommand(versionCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
