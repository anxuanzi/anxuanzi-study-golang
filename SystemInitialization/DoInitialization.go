package SystemInitialization

import (
	"az-ops/Global"
	"az-ops/Utils"
	_ "az-ops/resources/statik"
	"fmt"
	gonanoid "github.com/matoous/go-nanoid"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-cmd/cmd"
	"github.com/pterm/pterm"
	"github.com/spf13/afero"
)

func processInitialization() error {
	clear()
	intro()
	clear()
	yumCfg()
	time.Sleep(time.Second)
	timeCfg()
	time.Sleep(time.Second)
	kernelCfg()
	time.Sleep(time.Second)
	sysCfg()
	time.Sleep(time.Second)
	netCfg()
	time.Sleep(time.Second)
	installApps()
	time.Sleep(time.Second)
	donePage()
	return nil
}
func intro() {
	pterm.DefaultHeader.WithBackgroundStyle(pterm.NewStyle(pterm.BgLightBlue)).WithMargin(10).Println(
		"Configuration Processing")

	pterm.Info.Println("We are about to start configuration, all settings will be set based on your preference." +
		"\nIf the program make any change(s) with file(s), it will make a backup automatically." +
		"\nStarting configuration at: " + pterm.Green(time.Now().Format("02 Jan 2006 - 15:04:05 MST")))
	pterm.Println()
	introSpinner, _ := pterm.DefaultSpinner.WithRemoveWhenDone(true).Start("Waiting for 3 seconds...")
	time.Sleep(time.Second)
	for i := 2; i > 0; i-- {
		if i > 1 {
			introSpinner.UpdateText("Waiting for " + strconv.Itoa(i) + " seconds...")
		} else {
			introSpinner.UpdateText("Waiting for " + strconv.Itoa(i) + " second...")
		}
		time.Sleep(time.Second)
	}
	_ = introSpinner.Stop()
}
func yumCfg() {
	pterm.DefaultSection.Println("Configuring YUM/DNF")
	//1. update mirror
	//2. yum clean all
	//3. yum makecache
	if configurations.YumMirror == nil {
		pterm.Info.Println("Configuring YUM/DNF - SKIPPED")
		return
	}
	s, _ := pterm.DefaultSpinner.Start("Backup repo file")
	yumRepoFolder := afero.NewBasePathFs(Global.GetEtcFolder(), "yum.repos.d")
	err := yumRepoFolder.Rename("CentOS-Base.repo", "CentOS-Base.repo.bak")
	if err != nil {
		s.Fail("An error has occurred while trying to backup Repo file.")
		pterm.Error.Println("error details: ", err)
		os.Exit(0)
	}
	s.Success("Backup repo file")
	_ = s.Stop()
	s, _ = pterm.DefaultSpinner.Start("Creating new repo file")
	s.Success("Creating new repo file")
	_ = s.Stop()
	s, _ = pterm.DefaultSpinner.Start("Writing contents into new repo file")
	repoFile, err := Global.GetStatikFS().Open("/repos/centos" + strconv.FormatInt(int64(configurations.CentosVersion), 10) + "-" + configurations.YumMirror.RepoFileName)
	if err != nil {
		s.Fail("An error has occurred while trying to write Repo file")
		pterm.Error.Println("error details: ", err)
		os.Exit(0)
	}
	repoContent, err := ioutil.ReadAll(repoFile)
	if err != nil {
		s.Fail("An error has occurred while trying to write Repo file")
		pterm.Error.Println("error details: ", err)
		os.Exit(0)
	}
	defer repoFile.Close()
	createdRepo, err := yumRepoFolder.Create("CentOS-Base.repo")
	if err != nil {
		s.Fail("An error has occurred while trying to write Repo file")
		pterm.Error.Println("error details: ", err)
		os.Exit(0)
	}
	_, err = createdRepo.Write(repoContent)
	if err != nil {
		s.Fail("An error has occurred while trying to write Repo file")
		pterm.Error.Println("error details: ", err)
		os.Exit(0)
	}
	s.Success("Writing contents into new repo file")
	_ = s.Stop()
	s, _ = pterm.DefaultSpinner.Start("Cleaning YUM/DNF caches")
	if configurations.CentosVersion == 7 {
		Utils.ExecSimpleCmdWithResult("yum", "clean", "all", "-y", "&&", "rm", "-rf", "/var/cache/yum")
	} else {
		Utils.ExecSimpleCmdWithResult("dnf", "clean", "all", "-y", "&&", "rm", "-rf", "/var/cache/yum")
	}
	s.Success("Cleaning YUM/DNF caches")
	_ = s.Stop()
	pterm.Info.Println("Updating YUM/DNF (might take a while)")
	cmdOptions := cmd.Options{
		Buffered:  false,
		Streaming: true,
	}
	var envCmd *cmd.Cmd
	if configurations.CentosVersion == 7 {
		envCmd = cmd.NewCmdOptions(cmdOptions, "yum", "update", "-y")
	} else {
		envCmd = cmd.NewCmdOptions(cmdOptions, "dnf", "update", "-y")
	}
	doneChan := make(chan struct{})
	go func() {
		defer close(doneChan)
		for envCmd.Stdout != nil || envCmd.Stderr != nil {
			select {
			case line, open := <-envCmd.Stdout:
				if !open {
					envCmd.Stdout = nil
					continue
				}
				fmt.Println(line)
			case line, open := <-envCmd.Stderr:
				if !open {
					envCmd.Stderr = nil
					continue
				}
				_, _ = fmt.Fprintln(os.Stderr, line)
			}
		}
	}()
	// Run and wait for Cmd to return, discard Status
	<-envCmd.Start()
	// Wait for goroutine to print everything
	<-doneChan
	pterm.Success.Println("Updating YUM/DNF")
}
func timeCfg() {
	if configurations.Timezone != "" {
		pterm.DefaultSection.Println("Configuring Time Settings")
		a, _ := pterm.DefaultSpinner.Start("Checking required package")
		ex := Utils.ExecSimpleCmdWithResult("rpm", "-q", "chrony")
		if strings.Contains(ex.Stdout[len(ex.Stdout)-1], "is not installed") {
			//not installed
			a.Warning("Missing required package")
			_ = a.Stop()
			pterm.Info.Println("Installing required package")
			cmdOptions := cmd.Options{
				Buffered:  false,
				Streaming: true,
			}
			var envCmd *cmd.Cmd
			if configurations.CentosVersion == 7 {
				envCmd = cmd.NewCmdOptions(cmdOptions, "yum", "install", "chrony", "-y")
			} else {
				envCmd = cmd.NewCmdOptions(cmdOptions, "dnf", "update", "chrony", "-y")
			}
			doneChan := make(chan struct{})
			go func() {
				defer close(doneChan)
				for envCmd.Stdout != nil || envCmd.Stderr != nil {
					select {
					case line, open := <-envCmd.Stdout:
						if !open {
							envCmd.Stdout = nil
							continue
						}
						fmt.Println(line)
					case line, open := <-envCmd.Stderr:
						if !open {
							envCmd.Stderr = nil
							continue
						}
						_, _ = fmt.Fprintln(os.Stderr, line)
					}
				}
			}()
			// Run and wait for Cmd to return, discard Status
			<-envCmd.Start()
			// Wait for goroutine to print everything
			<-doneChan
			pterm.Info.Println("Starting chronyd")
			Utils.ExecSimpleCmdStreaming("systemctl", "start", "chronyd")
			pterm.Success.Println("Installing required package")
		} else {
			a.Success("Checking required package")
			_ = a.Stop()
		}
		s, _ := pterm.DefaultSpinner.Start("Changing timezone to " + configurations.Timezone)
		Utils.ExecSimpleCmdWithResult("timedatectl", "set-timezone", configurations.Timezone)
		s.Success("Changing timezone to " + configurations.Timezone)
		_ = s.Stop()
		s, _ = pterm.DefaultSpinner.Start("Checking config file")
		exist, err := afero.Exists(Global.GetEtcFolder(), "chrony.conf")
		if err != nil {
			s.Fail("An error has occurred while trying to find chrony.conf")
			pterm.Error.Println("error details: ", err)
			os.Exit(0)
		}
		if !exist {
			s.Fail("An error has occurred while trying to find chrony.conf")
			pterm.Error.Println("error details: ", "the config file does not exist!")
			os.Exit(0)
		}
		s.Success("Checking config file")
		_ = s.Stop()
		s, _ = pterm.DefaultSpinner.Start("Backup chrony config file")
		err = Global.GetEtcFolder().Rename("chrony.conf", "chrony.conf.bak")
		if err != nil {
			s.Fail("An error has occurred while trying to backup chrony config file.")
			pterm.Error.Println("error details: ", err)
			os.Exit(0)
		}
		s.Success("Backup chrony config file")
		_ = s.Stop()
		s, _ = pterm.DefaultSpinner.Start("Creating new chrony config file")
		s.Success("Creating new chrony config file")
		_ = s.Stop()
		s, _ = pterm.DefaultSpinner.Start("Writing contents into new chrony config file")
		chronyFile, err := Global.GetStatikFS().Open("/conf/chrony.conf")
		if err != nil {
			s.Fail("An error has occurred while trying to write chrony config file")
			pterm.Error.Println("error details: ", err)
			os.Exit(0)
		}
		chronyContent, err := ioutil.ReadAll(chronyFile)
		if err != nil {
			s.Fail("An error has occurred while trying to write chrony config file")
			pterm.Error.Println("error details: ", err)
			os.Exit(0)
		}
		defer chronyFile.Close()
		createdRepo, err := Global.GetEtcFolder().Create("chrony.conf")
		if err != nil {
			s.Fail("An error has occurred while trying to write chrony config file")
			pterm.Error.Println("error details: ", err)
			os.Exit(0)
		}
		_, err = createdRepo.Write(chronyContent)
		if err != nil {
			s.Fail("An error has occurred while trying to write chrony config file")
			pterm.Error.Println("error details: ", err)
			os.Exit(0)
		}
		s.Success("Writing contents into new chrony config file")
		_ = s.Stop()
		pterm.Info.Println("Restarting chrony service")
		Utils.ExecSimpleCmdStreaming("systemctl", "restart", "chronyd")
		Utils.ExecSimpleCmdStreaming("systemctl", "enable", "chronyd")
		pterm.Success.Println("Restarting chrony service")
	} else {
		pterm.Info.Println("Configuring Time Settings - SKIPPED")
	}
}
func kernelCfg() {
	pterm.DefaultSection.Println("Kernel & Network Optimization")
	if configurations.KernelOptimization {
		a, _ := pterm.DefaultSpinner.Start("Optimizing Kernel")
		Utils.ExecShellScript("kernel_optimization.sh")
		a.Success("Optimizing Kernel")
	} else {
		pterm.Info.Println("Optimizing Kernel - SKIPPED")
	}
	if configurations.BBR {
		a, _ := pterm.DefaultSpinner.Start("Configuring BBR")
		if configurations.CentosVersion == 7 {
			Utils.ExecShellScript("centos7_bbr.sh")
		}
		Utils.ExecShellScript("turn_on_bbr.sh")
		a.Success("Configuring BBR")
	} else {
		pterm.Info.Println("BBR Network Optimization - SKIPPED")
	}
}
func sysCfg() {
	pterm.DefaultSection.Println("System Settings Configuration")
	if !configurations.SELinux {
		a, _ := pterm.DefaultSpinner.Start("Configuring SELinux")
		Utils.ExecSimpleCmdWithResult("sed", "-i", "'s/SELINUX=enforcing/SELINUX=disabled/'", "/etc/selinux/config")
		Utils.ExecSimpleCmdWithResult("setenforce", "0")
		Utils.ExecSimpleCmdWithResult("grep", "SELINUX=disabled", "/etc/selinux/config")
		a.Success("Configuring SELinux")
	} else {
		pterm.Info.Println("Configuring SELinux - SKIPPED")
	}
	if !configurations.Firewall {
		a, _ := pterm.DefaultSpinner.Start("Configuring Firewall")
		Utils.ExecSimpleCmdWithResult("systemctl", "stop", "firewalld")
		Utils.ExecSimpleCmdWithResult("systemctl", "disable", "firewalld")
		a.Success("Configuring Firewall")
	} else {
		pterm.Info.Println("Configuring Firewall - SKIPPED")
	}
}
func netCfg() {
	pterm.DefaultSection.Println("Network Configuration")
	if configurations.DNS1 != "" && configurations.DNS2 != "" {
		a, _ := pterm.DefaultSpinner.Start("Configuring Network")
		err := ioutil.WriteFile("/etc/resolv.conf", []byte("nameserver "+configurations.DNS1+"\nnameserver "+configurations.DNS2), os.ModePerm)
		if err != nil {
			a.Fail("An error has occurred while trying to write DNS config file")
			pterm.Error.Println("error details: ", err)
			os.Exit(0)
		}
		a.Success("Configuring Network")
	} else {
		pterm.Info.Println("Configuring Network - SKIPPED")
	}
}
func installApps() {
	pterm.DefaultSection.Println("Installing Common Software")
	if configurations.Software != nil {
		apps := strings.Join(configurations.Software, " ")
		if configurations.CentosVersion == 7 {
			Utils.ExecSimpleCmdStreaming("yum", "install", apps, "-y")
		} else {
			Utils.ExecSimpleCmdStreaming("dnf", "install", apps, "-y")
		}
		pterm.Success.Println("Installing Common Software")
	} else {
		pterm.Info.Println("Installing Common Software - SKIPPED")
	}
}

func donePage() {
	f, _ := Global.GetEtcFolder().Create("az-ops-sys-initiated.lock")
	c, _ := gonanoid.ID(255)
	_, _ = f.WriteString(c)
	defer f.Close()
	// Generate BigLetters
	s, _ := pterm.DefaultBigText.WithLetters(pterm.NewLettersFromString("AZ-OPS")).Srender()
	pterm.DefaultCenter.Println(s) // Print BigLetters with the default CenterPrinter
	pterm.DefaultCenter.WithCenterEachLineSeparately().Println("Congratulations!\nWe has been successfully\ninitialized and configured\nyour operating system.\nHave a nice day!")
	os.Exit(0)
}
