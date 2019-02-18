package main

import (
	"fmt"

	clientcmd "github.com/bluele/rkvs/pkg/client/cmd"
	nodecmd "github.com/bluele/rkvs/pkg/node/cmd"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "rkvs",
	Short: "rkvs node",
}

func init() {
	rootCmd.AddCommand(
		nodecmd.GetStartCMD(),
		clientcmd.GetServersCMD(),
		clientcmd.GetKVSCMD(),
	)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
