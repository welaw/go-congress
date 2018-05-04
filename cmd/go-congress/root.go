package gocongress

import (
	"fmt"

	"github.com/spf13/cobra"
)

func MakeRootCmd() *cobra.Command {

	var goUsaAddr string
	var upstreamName string

	var rootCmd = &cobra.Command{
		Use:   "go-congress",
		Short: "An upstream for the United States Congress",
		Long: `go-congress is a welaw upstream for the United States Congress.

Laws are retrieved from fdsys.gov and votes from congress.org.

Website: 	http://welaw.org
Source:		http://github.com/welaw/welaw
`,
	}

	rootCmd.PersistentFlags().StringVar(&goUsaAddr, "grpc.addr", ":8091", "GoUSGPO GRPC listen address")

	rootCmd.PersistentFlags().StringVar(&upstreamName, "upstream.name", "congress", "The upstream identifier")
	//rootCmd.PersistentFlags().StringVar(&ballotGRPCAddr, "grpc.addr", ":8091", "Publish Service GRPC listen address")

	// server
	rootCmd.AddCommand(makeServeCmd(&goUsaAddr, &upstreamName))
	// client
	rootCmd.AddCommand(makeSendVoteCmd(&goUsaAddr))
	rootCmd.AddCommand(makeSendLawCmd(&goUsaAddr))
	rootCmd.AddCommand(makeStatusCmd(&goUsaAddr))

	return rootCmd
}

func requiredArgs(req int) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) < req {
			return fmt.Errorf("required args: %d", req)
		}
		return nil
	}
}
