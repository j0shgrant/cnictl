package cmd

import (
	"context"
	"fmt"
	"github.com/containerd/go-cni"
	"github.com/j0shgrant/cnictl/internal"
	"github.com/spf13/cobra"
	"os"
	"strconv"
	"time"
)

var removeCmd = &cobra.Command{
	Use:     "remove <NETNS> <ID>",
	Aliases: []string{"rm"},
	Short:   "Remove a CNI network",
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

		err = c.Remove(ctx, id, netnsDir + "/" + netns, cni.WithLabels(internal.Labels(id)))
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		fmt.Printf("Removed CNI network [%s]\n", id)
		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
