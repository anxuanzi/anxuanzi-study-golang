#!/bin/bash

rpm --import https://www.elrepo.org/RPM-GPG-KEY-elrepo.org

yum install https://www.elrepo.org/elrepo-release-7.el7.elrepo.noarch.rpm

yum --enablerepo=elrepo-kernel install kernel-ml -y

awk -F\' '$1=="menuentry " {print i++ " : " $2}' /etc/grub2.cfg

grub2-set-default 0

echo 'Kernel had been updated, please restart your OS after this installation!'
