package kube

// TODO
type KubeconfigClustersT struct {
	Name    string `yaml:"name"`
	Cluster KubeconfigClustersClusterT
}

type KubeconfigClustersClusterT struct {
	Server             string `yaml:"server"`
	InsecureSkipVerify bool   `yaml:"insecure-skip-tls-verify"`
}

// TODO
type KubeconfigContextT struct {
	Name    string                    `yaml:"name"`
	Context KubeconfigContextContextT `yaml:"context"`
}
type KubeconfigContextContextT struct {
	Cluster string `yaml:"cluster"`
	User    string `yaml:"user"`
}

// TODO
type KubeconfigUsersT struct {
	Name string               `yaml:"name"`
	User KubeconfigUsersUserT `yaml:"user"`
}

// TODO
type KubeconfigUsersUserT struct {
	Token string `yaml:"token"`
}

// TODO
type KubeconfigT struct {
	ApiVersion string                `yaml:"apiVersion"`
	Kind       string                `yaml:"kind"`
	Clusters   []KubeconfigClustersT `yaml:"clusters"`

	Contexts []KubeconfigContextT `yaml:"contexts"`

	Users []KubeconfigUsersT `yaml:"users"`

	CurrentContext string `yaml:"current-context"`
}
