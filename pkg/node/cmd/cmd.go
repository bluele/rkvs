package cmd

import (
	"context"
	"strings"
	"time"

	"github.com/bluele/rkvs/pkg/client"
	"github.com/bluele/rkvs/pkg/config"
	"github.com/bluele/rkvs/pkg/node"
	"github.com/bluele/rkvs/pkg/proto"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

const (
	flagLocalID   = "id"
	flagAddress   = "address"
	flagBootstrap = "bootstrap"
	flagJoin      = "join"
)

func GetStartCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "start",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			viper.BindPFlags(cmd.Flags())
			logger := logrus.New()

			addr := viper.GetString(flagAddress)
			id := viper.GetString(flagLocalID)
			cfg := config.MakeDefaultConfig(addr, id)
			joinTo := viper.GetString(flagJoin)

			n, err := node.MakeNode(cfg, id, logger)
			if err != nil {
				return err
			}
			if viper.GetBool(flagBootstrap) {
				n.Bootstrap(id, addr)
			}

			ctx := context.Background()
			go func() {
				time.Sleep(time.Second)
				for _, node := range parseNodeInfos(joinTo) {
					cp := client.NewPoolingConnector(client.NewConnectionPool(), []string{node.Addr}, grpc.WithInsecure())
					cl := client.NewSystemClient(cp)

					_, err := cl.Join(ctx, &proto.SystemRequestJoin{
						Id:   []byte(id),
						Addr: []byte(addr),
					})
					if err != nil {
						panic(err)
					}
				}
			}()
			return n.Serve()
		},
	}
	cmd.Flags().String(flagLocalID, "node-id", "local id")
	cmd.Flags().String(flagAddress, "127.0.0.1:10000", "Raft listen address")
	cmd.Flags().Bool(flagBootstrap, false, "run as bootstrap mode")
	cmd.Flags().String(flagJoin, "", "join to address")
	return cmd
}

type nodeInfo struct {
	ID   string
	Addr string
}

func parseNodeInfos(parts string) (nodes []nodeInfo) {
	if parts == "" {
		return nil
	}
	for _, part := range strings.Split(parts, ",") {
		nodes = append(nodes, parseNodeInfo(part))
	}
	return nodes
}

func parseNodeInfo(part string) nodeInfo {
	cols := strings.Split(part, "@")
	return nodeInfo{ID: cols[0], Addr: cols[1]}
}
