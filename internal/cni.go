package internal

import "github.com/containerd/go-cni"

func NewCNI(interfacePrefix, pluginDir, pluginConfDir string) (cni.CNI, error) {
	c, err := cni.New(
		cni.WithPluginConfDir(pluginConfDir),
		cni.WithPluginDir([]string{pluginDir}),
		cni.WithInterfacePrefix(interfacePrefix),
	)
	if err != nil {
		return nil, err
	}

	err = c.Load(
		cni.WithLoNetwork,
		cni.WithDefaultConf,
	)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func Labels(netns string) map[string]string {
	return map[string]string{
		"NETWORK_ID":         netns,
		"IgnoreUnknown": "1",
	}
}
