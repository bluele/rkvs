package cmd

import (
	"context"
	"fmt"

	"github.com/bluele/rkvs/pkg/client"
	"github.com/bluele/rkvs/pkg/proto"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

const (
	flagAddress = "address"
)

func GetServersCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "servers",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			viper.BindPFlags(cmd.Flags())
			ctx := context.Background()
			res, err := getSystemClient().Servers(ctx, new(proto.SystemRequestServers))
			if err != nil {
				return err
			}
			for _, info := range res.Infos {
				fmt.Printf("suffrage=%v id=%v address=%v\n", info.Suffrage, info.Id, info.Address)
			}
			return nil
		},
	}
	cmd.Flags().String(flagAddress, "127.0.0.1:10000", "Raft listen address")
	return cmd
}

func getSystemClient() proto.SystemClient {
	return client.NewSystemClient(
		client.NewCommonConnector(
			viper.GetString(flagAddress),
			grpc.WithInsecure(),
		),
	)
}
