package Model

import "github.com/go-ping/ping"

type MirrorInfo struct {
	Name         string `json:"name"`
	DomainName   string `json:"domain_name"`
	RepoFileName string `json:"repo_file_name"`
	Region       string `json:"region"`
}

type Mirror struct {
	MirrorInfo *MirrorInfo      `json:"mirror_info"`
	Stat       *ping.Statistics `json:"stat"`
}

type InitConfig struct {
	CentosVersion      int         `json:"centos_version"`
	YumMirror          *MirrorInfo `json:"yum_mirror"`
	Timezone           string      `json:"time_zone"`
	KernelOptimization bool        `json:"kernel_optimization"`
	BBR                bool        `json:"bbr"`
	SELinux            bool        `json:"se_linux"`
	Firewall           bool        `json:"firewall"`
	DNS1               string      `json:"dns_1"`
	DNS2               string      `json:"dns_2"`
	Software           []string    `json:"software"`
}

type DnsInfo struct {
	ProviderName string `json:"provider_name"`
	DNS1         string `json:"dns_1"`
	DNS2         string `json:"dns_2"`
}
type Dns struct {
	DnsInfo *DnsInfo         `json:"dns_info"`
	Stat    *ping.Statistics `json:"stat"`
}
