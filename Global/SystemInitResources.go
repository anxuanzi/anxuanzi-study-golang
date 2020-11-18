package Global

import "az-ops/Model"

var mirrors []*Model.MirrorInfo
var timezones []string
var basicSoftware []string
var dns []*Model.DnsInfo

func initMirrors() {
	mirrors = make([]*Model.MirrorInfo, 4)

	mirrors[0] = &Model.MirrorInfo{
		Name:         "xTom",
		DomainName:   "mirrors.xtom.com",
		RepoFileName: "xtom.repo",
		Region:       "USA-CA",
	}
	mirrors[1] = &Model.MirrorInfo{
		Name:         "xTom",
		DomainName:   "mirrors.xtom.com.hk",
		RepoFileName: "xtom-hk.repo",
		Region:       "CN-HK",
	}
	mirrors[2] = &Model.MirrorInfo{
		Name:         "Tuna",
		DomainName:   "mirrors.tuna.tsinghua.edu.cn",
		RepoFileName: "tuna.repo",
		Region:       "CN-BJ",
	}
	mirrors[3] = &Model.MirrorInfo{
		Name:         "Tencent",
		DomainName:   "mirrors.cloud.tencent.com",
		RepoFileName: "tencent.repo",
		Region:       "CN-BGP",
	}
}
func initTimezones() {
	timezones = append(timezones, "America/Chicago")
	timezones = append(timezones, "America/Los_Angeles")
	timezones = append(timezones, "America/New_York")
	timezones = append(timezones, "America/Phoenix")
	timezones = append(timezones, "Asia/Hong_Kong")
	timezones = append(timezones, "Asia/Shanghai")
	timezones = append(timezones, "Asia/Singapore")
	timezones = append(timezones, "Asia/Tokyo")
	timezones = append(timezones, "Australia/Adelaide")
	timezones = append(timezones, "Europe/Berlin")
	timezones = append(timezones, "Europe/London")
	timezones = append(timezones, "Europe/Paris")
	timezones = append(timezones, "Europe/Vienna")
}
func initBasicSoftwareList() {
	basicSoftware = append(basicSoftware, "epel-release")
	basicSoftware = append(basicSoftware, "curl")
	basicSoftware = append(basicSoftware, "wget")
	basicSoftware = append(basicSoftware, "telnet")
	basicSoftware = append(basicSoftware, "vim")
	basicSoftware = append(basicSoftware, "screen")
	basicSoftware = append(basicSoftware, "make")
	//basicSoftware = append(basicSoftware, "gcc")
	basicSoftware = append(basicSoftware, "net-tools")
	//basicSoftware = append(basicSoftware, "perl")
	//basicSoftware = append(basicSoftware, "ruby")
	//basicSoftware = append(basicSoftware, "kernel-devel")
}

func initDNSList() {
	dns = append(dns, &Model.DnsInfo{
		ProviderName: "OpenDNS",
		DNS1:         "208.67.222.222",
		DNS2:         "208.67.220.220",
	})
	dns = append(dns, &Model.DnsInfo{
		ProviderName: "Cloudflare",
		DNS1:         "1.1.1.1",
		DNS2:         "1.0.0.1",
	})
	dns = append(dns, &Model.DnsInfo{
		ProviderName: "Google",
		DNS1:         "8.8.8.8",
		DNS2:         "8.8.4.4",
	})
	dns = append(dns, &Model.DnsInfo{
		ProviderName: "Quad9",
		DNS1:         "9.9.9.9",
		DNS2:         "149.112.112.112",
	})
	dns = append(dns, &Model.DnsInfo{
		ProviderName: "AliDNS",
		DNS1:         "223.5.5.5",
		DNS2:         "223.6.6.6",
	})
	dns = append(dns, &Model.DnsInfo{
		ProviderName: "DNSPod",
		DNS1:         "119.29.29.29",
		DNS2:         "119.28.28.28",
	})
	dns = append(dns, &Model.DnsInfo{
		ProviderName: "114DNS",
		DNS1:         "114.114.114.114",
		DNS2:         "114.114.115.115",
	})
}

func GetMirrors() []*Model.MirrorInfo {
	return mirrors
}
func GetTimezones() []string {
	return timezones
}
func GetBasicSoftwareList() []string {
	return basicSoftware
}
func GetDNS() []*Model.DnsInfo {
	return dns
}
