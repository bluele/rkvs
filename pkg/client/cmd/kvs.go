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

func GetKVSCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use: "kvs",
	}
	cmd.PersistentFlags().String(flagAddress, "127.0.0.1:10000", "Raft listen address")
	cmd.AddCommand(GetPingCMD(), GetWriteCMD(), GetReadCMD())
	return cmd
}

func getKVSClient() proto.KVSClient {
	return client.NewKVSClient(
		client.NewCommonConnector(
			viper.GetString(flagAddress),
			grpc.WithInsecure(),
		),
	)
}

func GetPingCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "ping",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			viper.BindPFlags(cmd.Flags())
			ctx := context.Background()
			_, err := getKVSClient().Ping(ctx, &proto.KVSRequestPing{})
			return err
		},
	}
	return cmd
}

func GetWriteCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "write",
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			viper.BindPFlags(cmd.Flags())
			ctx := context.Background()
			key, value := args[0], args[1]
			_, err := getKVSClient().Write(ctx, &proto.KVSRequestWrite{Key: []byte(key), Value: []byte(value)})
			return err
		},
	}
	return cmd
}

func GetReadCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "read",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			viper.BindPFlags(cmd.Flags())
			ctx := context.Background()
			key := args[0]
			res, err := getKVSClient().Read(ctx, &proto.KVSRequestRead{Key: []byte(key)})
			if err != nil {
				return err
			}
			fmt.Println(string(res.Value))
			return nil
		},
	}
	return cmd
}
