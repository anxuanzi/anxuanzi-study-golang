package SystemInitialization

import (
	"az-ops/Global"
	"az-ops/Model"
	"az-ops/Utils"
	"errors"
	"net"
	"sort"
	"strconv"
	"time"

	"github.com/go-ping/ping"
	"github.com/gookit/gcli/v2"
	"github.com/pterm/pterm"
	"github.com/spf13/afero"
)

var InitSystemCommand = &gcli.Command{
	Name:     "systeminit",
	UseFor:   "Initialize the operating system (time settings, repository mirrors, network optimization, kernel parameter optimization, commonly used software installation).",
	Aliases:  []string{"sysinit", "sinit"},
	Config:   nil, //没有额外参数
	Examples: "{$binName} {$cmd}",
	Func:     app,
}

const second = time.Second

var configurations *Model.InitConfig

func app(c *gcli.Command, args []string) error {
	clear()
	configurations = &Model.InitConfig{
		CentosVersion:      0,
		YumMirror:          nil,
		Timezone:           "",
		KernelOptimization: false,
		BBR:                false,
		SELinux:            false,
		Firewall:           false,
		DNS1:               "",
		DNS2:               "",
		Software:           nil,
	}
	_ = pterm.DefaultBigText.WithLetters(
		pterm.NewLettersFromStringWithStyle("AZ", pterm.NewStyle(pterm.FgCyan)),
		pterm.NewLettersFromStringWithStyle(" OPS", pterm.NewStyle(pterm.FgLightMagenta))).
		Render()
	pterm.DefaultHeader.WithBackgroundStyle(pterm.NewStyle(pterm.BgLightBlue)).WithMargin(10).Println(
		"System Initialization Operator")
	loading, _ := pterm.DefaultSpinner.WithRemoveWhenDone(true).Start("Checking system requirements and preparing app...")
	//初始化文件系统资源
	Global.InitResources()
	//获取centos发布文件是否存在
	exists, _ := afero.Exists(Global.GetEtcFolder(), "redhat-release")
	//判断centos发布文件是否存在（确认系统为centos系列）
	if !exists {
		//系统不为centos系列，提示并退出
		pterm.Error.Println("Your OS is not Redhat CentOS series, this program only supports CentOS series!")
		loading.Fail("OS not supported!")
		return nil
	}
	//判断系统是否已经初始化过
	//获取安装锁文件
	exists, _ = afero.Exists(Global.GetEtcFolder(), "az-ops-sys-initiated.lock")
	//判断初始化锁存在
	if exists {
		//系统已经初始化过
		pterm.Error.Println("You already initiated your OS, not allow to do it again!")
		loading.Fail("initiated error!")
		return nil
	}
	loading.Stop()
	pterm.Info.Println("We are going to ask few questions about how to configure your OS" +
		"\nThis program will help you determine which YUM/DNF repository mirror to use." +
		"\nThis program will help you configure Time Zone and time sync." +
		"\nThis program could help you optimize system performance." +
		"\nThis program could help you change the SELinux status." +
		"\nThis program could help you change the Firewall status." +
		"\nThis program will help you determine which DNS to use." +
		"\nThis program will help you install all basic useful software or commands.")
	pterm.Println()
	introSpinner, _ := pterm.DefaultSpinner.WithRemoveWhenDone(true).Start("Program will start in 5 seconds...")
	time.Sleep(second)
	for i := 4; i > 0; i-- {
		if i > 1 {
			introSpinner.UpdateText("Program will start in " + strconv.Itoa(i) + " seconds...")
		} else {
			introSpinner.UpdateText("Program will start in " + strconv.Itoa(i) + " second...")
		}
		time.Sleep(second)
	}
	introSpinner.Stop()
	//清屏
	clear()
	_, osv := Utils.SingleChoicePrompt("Which CentOS version are you using?", []string{"CentOS-7", "CentOS-8"}, false)
	if osv == "CentOS-7" {
		configurations.CentosVersion = 7
	} else {
		configurations.CentosVersion = 8
	}
	clear()
	pterm.DefaultHeader.Println("YUM/DNF Mirror Configuration Wizard")
	res := Utils.YesNoPrompt("Do you want configure YUM/DNF repository mirror?")
	if res {
		err := configYumMirror()
		if err != nil {
			pterm.Error.Printf("\n something went wrong: %v\n", err)
			return nil
		}
	}
	pterm.Success.Println("YUM/DNF Settings DONE!")
	clear()
	pterm.DefaultHeader.Println("Time Configuration Wizard")
	res = Utils.YesNoPrompt("Do you want configure time settings?")
	if res {
		err := configTimeZone()
		if err != nil {
			pterm.Error.Printf("\n something went wrong: %v\n", err)
			return nil
		}
	}
	pterm.Success.Println("Time Settings DONE!")
	clear()
	pterm.DefaultHeader.Println("Kernel & Network Optimization")
	res = Utils.YesNoPrompt("Do you want configure Kernel optimization parameters?")
	if res {
		configurations.KernelOptimization = true
	}
	res = Utils.YesNoPrompt("Do you want enable BBR TCP optimization?")
	if res {
		configurations.BBR = true
	}
	pterm.Success.Println("Kernel Settings DONE!")
	clear()
	pterm.DefaultHeader.Println("System Settings")
	res = Utils.YesNoPrompt("Do you want enable SELinux?")
	if res {
		configurations.SELinux = true
	}
	res = Utils.YesNoPrompt("Do you want enable Firewall?")
	if res {
		configurations.Firewall = true
	}
	pterm.Success.Println("System Settings DONE!")
	clear()
	pterm.DefaultHeader.Println("Network Settings")
	res = Utils.YesNoPrompt("Do you want change your DNS server?")
	if res {
		err := configNetworks()
		if err != nil {
			pterm.Error.Printf("\n something went wrong: %v\n", err)
			return nil
		}
	}
	pterm.Success.Println("Network Settings DONE!")
	clear()
	pterm.DefaultHeader.Println("Utility Software Installation Checklist")
	softwareInstallChecklist()
	res = Utils.YesNoPrompt("Start Processing System Initialization?")
	if res {
		clear()
		err := processInitialization()
		if err != nil {
			pterm.Error.Printf("\n something went wrong: %v\n", err)
			return nil
		}
	}
	return nil
}
func configYumMirror() error {
	pterm.Info.Println("Start testing fastest mirror...")
	mirrorList := make([]*Model.Mirror, len(Global.GetMirrors()))
	//针对每一个镜像
	for index, mirror := range Global.GetMirrors() {
		spt, _ := pterm.DefaultSpinner.Start("Testing mirror " + mirror.Name + "(" + mirror.Region + ")...")
		pinger, err := ping.NewPinger(mirror.DomainName)
		if err != nil {
			spt.Fail("Error when creating Pinger")
			return err
		}
		pinger.SetPrivileged(true)
		pinger.Count = 5
		pinger.Timeout = time.Second * 3
		//pinger.Debug = true
		err = pinger.Run()
		if err != nil {
			spt.Fail("Something went wrong while trying to ping the domain name " + mirror.DomainName)
			return err
		}
		mirrorList[index] = &Model.Mirror{
			MirrorInfo: mirror,
			Stat:       pinger.Statistics(),
		}
		//if pinger.Statistics().PacketLoss >= 0.95 {
		//	spt.Fail("testing " + mirror.Name + "(" + mirror.Region + ") done! [poor connection]")
		//} else {
		spt.Success("testing " + mirror.Name + "(" + mirror.Region + ") done! [" + strconv.FormatInt(pinger.Statistics().AvgRtt.Milliseconds(), 10) + "ms]")
		//}
	}
	//测速完成，对结果进行排序
	sort.Slice(mirrorList, func(i, j int) bool {
		return mirrorList[i].Stat.AvgRtt.Milliseconds() < mirrorList[j].Stat.AvgRtt.Milliseconds()
	})
	options := make([]string, len(mirrorList))
	for i, v := range mirrorList {
		options[i] = v.MirrorInfo.Name + " (" + v.MirrorInfo.Region + ") [" + strconv.FormatInt(v.Stat.AvgRtt.Milliseconds(), 10) + "ms]"
	}
	pterm.Info.Println("Mirror speed test finished, please be ready to choose a mirror to use.")
	index, _ := Utils.SingleChoicePrompt("Please choose a mirror to use", options, false)
	configurations.YumMirror = mirrorList[index].MirrorInfo
	pterm.Success.Println("You choose mirror " + mirrorList[index].MirrorInfo.Name + " (" + mirrorList[index].MirrorInfo.Region + ") to use!")
	return nil
}
func configTimeZone() error {
	pterm.Info.Println("In this section you are going to chose a timezone and program will setup time sync for you.")
	_, tz := Utils.SingleChoicePrompt("Please select a timezone for this server", Global.GetTimezones(), false)
	configurations.Timezone = tz
	pterm.Info.Println("The NTP Sync will use default pool 'pool.ntp.org'.")
	return nil
}
func configNetworks() error {
	pterm.Info.Println("Now you are going to setup DNS server, step by step.")
	res := Utils.YesNoPrompt("Do you want benchmark DNS providers on this server?")
	rankedDns := make([]*Model.Dns, len(Global.GetDNS()))
	if res {
		pterm.Info.Println("Benchmark starting (using ping method)...")
		var err error
		rankedDns, err = dnsBenchmark()
		if err != nil {
			return err
		}
	} else {
		res = Utils.YesNoPrompt("Do you want enter DNS manually?")
		if res {
			firstDns := Utils.InputPrompt("Please enter the first DNS server", func(input string) error {
				if net.ParseIP(input) == nil {
					return errors.New("invalid number IP address format")
				}
				return nil
			})
			secondDns := Utils.InputPrompt("Please enter the second DNS server", func(input string) error {
				if net.ParseIP(input) == nil {
					return errors.New("invalid number IP address format")
				}
				return nil
			})
			pterm.Info.Println("Please confirm what you entered: ")
			pterm.Info.Println("DNS Server 1: " + firstDns)
			pterm.Info.Println("DNS Server 2: " + secondDns)
			confirm := Utils.YesNoPrompt("Are these DNS Server correct?")
			if confirm {
				configurations.DNS1 = firstDns
				configurations.DNS2 = secondDns
				return nil
			} else {
				clear()
				_ = configNetworks()
			}
		} else {
			for i, dnsInfo := range Global.GetDNS() {
				rankedDns[i] = &Model.Dns{
					DnsInfo: dnsInfo,
					Stat:    nil,
				}
			}
		}
	}
	options := make([]string, len(rankedDns))
	for index, dns := range rankedDns {
		if dns.Stat != nil {
			options[index] = dns.DnsInfo.ProviderName + " [" + strconv.FormatInt(dns.Stat.AvgRtt.Milliseconds(), 10) + "ms]"
		} else {
			options[index] = dns.DnsInfo.ProviderName
		}
	}
	selectedIndex, _ := Utils.SingleChoicePrompt("Please select your preferred DNS provider", options, true)
	pterm.Info.Println("Please confirm what you choose: ")
	pterm.Info.Println("DNS Provider: " + rankedDns[selectedIndex].DnsInfo.ProviderName)
	pterm.Info.Println("DNS Server 1: " + rankedDns[selectedIndex].DnsInfo.DNS1)
	pterm.Info.Println("DNS Server 2: " + rankedDns[selectedIndex].DnsInfo.DNS2)
	confirm := Utils.YesNoPrompt("Are these information correct?")
	if confirm {
		configurations.DNS1 = rankedDns[selectedIndex].DnsInfo.DNS1
		configurations.DNS2 = rankedDns[selectedIndex].DnsInfo.DNS2
		return nil
	} else {
		clear()
		_ = configNetworks()
	}
	return nil
}
func dnsBenchmark() ([]*Model.Dns, error) {
	rankedDns := make([]*Model.Dns, len(Global.GetDNS()))
	for index, dnsInfo := range Global.GetDNS() {
		spt, _ := pterm.DefaultSpinner.Start("Testing DNS provider " + dnsInfo.ProviderName + "...")
		pinger, err := ping.NewPinger(dnsInfo.DNS1)
		if err != nil {
			spt.Fail("Error when creating Pinger")
			return nil, err
		}
		pinger.SetPrivileged(true)
		pinger.Count = 5
		pinger.Timeout = time.Second * 3
		err = pinger.Run()
		if err != nil {
			spt.Fail("Something went wrong while trying to ping the DNS provider" + dnsInfo.ProviderName)
			return nil, err
		}
		rankedDns[index] = &Model.Dns{
			DnsInfo: dnsInfo,
			Stat:    pinger.Statistics(),
		}
		//if pinger.Statistics().PacketLoss > 0.95 {
		//	spt.Fail("testing provider " + dnsInfo.ProviderName + " done! [poor connection]")
		//} else {
		spt.Success("testing provider " + dnsInfo.ProviderName + " done! [" + strconv.FormatInt(pinger.Statistics().AvgRtt.Milliseconds(), 10) + "ms]")
		//}
	}
	sort.Slice(rankedDns, func(i, j int) bool {
		return rankedDns[i].Stat.AvgRtt.Milliseconds() < rankedDns[j].Stat.AvgRtt.Milliseconds()
	})
	return rankedDns, nil
}
func softwareInstallChecklist() {
	// Print info.
	pterm.Info.Println("Will install all of the following software.")
	from := pterm.NewRGB(0, 255, 255) // This RGB value is used as the gradients start point.
	to := pterm.NewRGB(255, 0, 255)   // This RGB value is used as the gradients end point.
	// For loop over the range of the terminal height.
	for i := 0; i < len(Global.GetBasicSoftwareList()); i++ {
		from.Fade(0, float32(pterm.GetTerminalHeight()-2), float32(i), to).Println(Global.GetBasicSoftwareList()[i])
	}
	configurations.Software = Global.GetBasicSoftwareList()
}
func clear() {
	print("\033[H\033[2J")
}
