# AZ-OPS

### Overview
Automated operation and maintenance script collection, can do things automatically without out typing command by yourself. Only support Centos7/Centos8. This program is designed for me to operate my servers (VMs) more efficiently. This project is also my first open-source project written with Golang (I'm still studying it!).

### Features
- [x] Perform an OS initialization operation.
    - Change Yum mirror to the fastest one.
    - Change Timezone settings based on your selection.
    - Configure NTP sync settings. (if the software doesn't exist, then install it).
    - Optimize Kernel parameter settings (not so sure what's that called in English).
    - Configure BBR TCP acceleration (will install the latest kernel for Centos7).
    - Turn off SELinux if needed.
    - Turn off and disable the firewall if needed.
    - Change the DNS server based on a benchmark (you can choose not to change, not perform a benchmark or enter DNS manually).
    - Install some really useful software (epel-release, curl, wget, telnet, vim, screen, make, net-tools).

- [ ] Network configuration
    - Create a bridged network
    
- [ ] One-click software installation and configuration (using Shell-Script or Ansible-Playbook)
    - App Store (you can pick what to install in a list)

Waiting for you to contribute more thoughts...

### Usage

- download the latest release (you can put it in `/bin`)

- set permission `chmod +x az-ops`

- run to see commands `./az-ops` if you put it under the `/bin` folder, you can simply run `az-ops` command wherever you want.

### Extra Information
There is nothing important here, just want say feel free to help me improve my code and let make this little program better!

### Some Screenshots

![image-20201117171410509](README.assets/image-20201117171410509.png)

![image-20201117171441239](README.assets/image-20201117171441239.png)

