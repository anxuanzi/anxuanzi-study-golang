#!/bin/bash

ulimit -SHn 1024000
echo "ulimit -SHn 1024000" >>/etc/rc.d/rc.local
source /etc/rc.d/rc.local

FILE=/etc/sysctl.conf
if [ -f "$FILE" ]; then
  rm -rf "$FILE"
fi

touch /etc/sysctl.conf

#https://www.yisu.com/zixun/153661.html
echo 'net.ipv4.tcp_syncookies = 1' >>/etc/sysctl.conf
echo 'net.ipv4.tcp_abort_on_overflow = 1' >>/etc/sysctl.conf
echo 'net.ipv4.tcp_max_tw_buckets = 6000' >>/etc/sysctl.conf
echo 'net.ipv4.tcp_sack = 1' >>/etc/sysctl.conf
echo 'net.ipv4.tcp_window_scaling = 1' >>/etc/sysctl.conf
echo 'net.ipv4.tcp_rmem = 4096    87380  4194304' >>/etc/sysctl.conf
echo 'net.ipv4.tcp_wmem = 4096    66384  4194304' >>/etc/sysctl.conf
echo 'net.ipv4.tcp_mem = 94500000 915000000 927000000' >>/etc/sysctl.conf
echo 'net.core.optmem_max = 81920' >>/etc/sysctl.conf
echo 'net.core.wmem_default = 8388608' >>/etc/sysctl.conf
echo 'net.core.wmem_max = 16777216' >>/etc/sysctl.conf
echo 'net.core.rmem_default = 8388608' >>/etc/sysctl.conf
echo 'net.core.rmem_max = 16777216' >>/etc/sysctl.conf
echo 'net.ipv4.tcp_max_syn_backlog = 1020000' >>/etc/sysctl.conf
echo 'net.core.netdev_max_backlog = 862144' >>/etc/sysctl.conf
echo 'net.core.somaxconn = 262144' >>/etc/sysctl.conf
echo 'net.ipv4.tcp_max_orphans = 327680' >>/etc/sysctl.conf
echo 'net.ipv4.tcp_timestamps = 0' >>/etc/sysctl.conf
echo 'net.ipv4.tcp_synack_retries = 1' >>/etc/sysctl.conf
echo 'net.ipv4.tcp_syn_retries = 1' >>/etc/sysctl.conf
echo 'net.ipv4.tcp_tw_reuse = 1' >>/etc/sysctl.conf
echo 'net.ipv4.tcp_fin_timeout = 15' >>/etc/sysctl.conf
echo 'net.ipv4.tcp_keepalive_time = 30' >>/etc/sysctl.conf
echo 'net.ipv4.ip_local_port_range = 1024  65000' >>/etc/sysctl.conf
echo 'net.netfilter.nf_conntrack_tcp_timeout_established = 180' >>/etc/sysctl.conf
echo 'net.netfilter.nf_conntrack_max = 1048576' >>/etc/sysctl.conf
echo 'net.nf_conntrack_max = 1048576' >>/etc/sysctl.conf

sysctl -p

# - 优化文件描述符
echo '  *  -  nofile   100000  ' >>/etc/security/limits.conf
