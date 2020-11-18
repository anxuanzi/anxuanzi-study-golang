#!/bin/bash

echo net.core.default_qdisc=fq >>/etc/sysctl.conf

echo net.ipv4.tcp_congestion_control=bbr >>/etc/sysctl.conf

sysctl -p

echo "Verify configuration status:"

sysctl net.ipv4.tcp_available_congestion_control

echo "Double check configuration status:"

lsmod | grep bbr
