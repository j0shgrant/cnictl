package cmd

import (
	"context"
	"crypto/sha512"
	"fmt"
	"github.com/containernetworking/cni/libcni"
	"github.com/spf13/cobra"
	"os"
	"strconv"
	"time"
)

var addCmd = &cobra.Command{
	Use:   "add <NETNS>",
	Short: "Create a new CNI network",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		// Load config
		pluginDir := cmd.Flag("plugin-dir").Value.String()
		pluginConfDir := cmd.Flag("plugin-conf-dir").Value.String()
		timeoutStr := cmd.Flag("timeout").Value.String()
		timeout, _ := strconv.Atoi(timeoutStr)

		// Parse args
		netns := args[0]
		netConfName := args[1]
		s := sha512.Sum512([]byte(netns))
		containerId := fmt.Sprintf("cnictl-%x", s[:10])
		ifName := "eth0"

		// Handle context
		ctx, cancelFunc := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
		defer cancelFunc()

		// Add to network
		cniConfig := libcni.NewCNIConfig([]string{pluginDir}, nil)
		netConfig, err := libcni.ConfFromFile(pluginConfDir + "/" + netConfName)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		runtimeConfig := &libcni.RuntimeConf{
			ContainerID:    containerId,
			NetNS:          netns,
			IfName:         ifName,
			Args:           nil,
			CapabilityArgs: nil,
		}

		result, err := cniConfig.AddNetwork(ctx, netConfig, runtimeConfig)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		_ = result.Print()
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
