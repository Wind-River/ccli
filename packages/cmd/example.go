package cmd

import (
	"fmt"
	"wrs/catalog/ccli/packages/config"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

func Example(configFile *config.ConfigData) *cobra.Command {
	var verboseFlag bool
	rootCmd := &cobra.Command{
		Use:   "ccli",
		Short: "Ccli is used to interact with the Software Parts Catalog.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if verboseFlag {
				zerolog.SetGlobalLevel(0)
			}
			fmt.Println("Please provide a command to be executed. Refer to the following examples and use help for more information.")
			exampleString :=
				`	$ ccli add part openssl-1.1.1n.yml
	$ ccli add profile profile_openssl-1.1.1n.yml
	$ ccli query "{part(id:\"aR25sd-V8dDvs2-p3Gfae\"){file_verification_code}}"
	$ ccli export part id sdl3ga-naTs42g5-rbow2A -o file.yml
	$ ccli export template security -o file.yml
	$ ccli update openssl-1.1.1n.v4.yml
	$ ccli upload openssl-1.1.1n.tar.gz
	$ ccli find part busybox
	$ ccli find sha256 2493347f59c03...
	$ ccli find profile security werS12-da54FaSff-9U2aef
	$ ccli delete adjb23-A4D3faTa-d95Xufs
	$ ccli ping
	$ ccli `
			fmt.Printf("%s\n", exampleString)
			return nil
		},
	}
	rootCmd.PersistentFlags().BoolVarP(&verboseFlag, "verbose", "v", false, "To Execute commands in verbose mode")
	return rootCmd
}
