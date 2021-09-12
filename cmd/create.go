package cmd

import (
	"context"
	"fmt"
	"github.com/containerd/go-cni"
	"github.com/j0shgrant/cnictl/internal"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strconv"
	"time"
)

var createCmd = &cobra.Command{
	Use:     "create <NETNS> <ID>",
	Aliases: []string{"c"},
	Short:   "Create a new CNI network",
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		// Load config
		netnsDir := cmd.Flag("netns-dir").Value.String()
		pluginDir := cmd.Flag("plugin-dir").Value.String()
		pluginConfDir := cmd.Flag("plugin-conf-dir").Value.String()
		timeoutStr := cmd.Flag("timeout").Value.String()
		timeout, _ := strconv.Atoi(timeoutStr)

		netns := args[0]
		id := args[1]

		// Handle context
		ctx, cancelFunc := context.WithTimeout(context.Background(), time.Duration(timeout) * time.Second)
		defer cancelFunc()

		// Create CNI
		c, err := internal.NewCNI("eth", pluginDir, pluginConfDir)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		// Load the cni configuration
		if err := c.Load(cni.WithLoNetwork, cni.WithDefaultConf); err != nil {
			log.Fatalf("failed to load cni configuration: %v", err)
		}

		res, err := c.Setup(ctx, id, netnsDir + "/" + netns, cni.WithLabels(internal.Labels(id)))
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		res.

		fmt.Printf("Created CNI network [%s]\n", id)
		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
