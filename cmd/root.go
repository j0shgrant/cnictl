package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "cnictl",
	Short: "cnictl is a CLI utility for managing local CNI networks",
	Long:  "cnictl is a CLI utility for managing local CNI networks using https://github.com/containerd/go-cni.",
}

var netnsDir, pluginDir, pluginConfDir string
var timeout int

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&netnsDir, "netns-dir", "/var/run/netns", "Path to CNI plugin binaries directory.")
	rootCmd.PersistentFlags().StringVar(&pluginDir, "plugin-dir", "/opt/cni/bin", "Path to CNI plugin binaries directory.")
	rootCmd.PersistentFlags().StringVar(&pluginConfDir, "plugin-conf-dir", "/etc/cni/net.d", "Path to CNI plugin binaries directory.")
	rootCmd.PersistentFlags().IntVar(&timeout, "timeout", 30, "Maximum duration to wait for CNI operation before timing out.")
}