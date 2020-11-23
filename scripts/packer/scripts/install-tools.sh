#!/usr/bin/env sh

apk update
echo "export CONFIG_PATH=/root/config.yml" >> /etc/profile
apk add --no-cache sudo
apk add -U wireguard-tools
apk add --no-cache openrc-doc
apk add --no-cache iptables
apk add --no-cache virtualbox-guest-additions virtualbox-guest-modules-virt
apk add --no-cache wget
apk add --no-cache unzip
apk add --no-cache openssh


wget https://github.com/aau-network-security/gwireguard/releases/download/v1.0.3/gwireguard_1.0.3_linux_64-bit.zip
unzip gwireguard_1.0.3_linux_64-bit.zip && mv gwireguard_1.0.3_linux_64-bit/wgsservice /etc/init.d/
chmod +x /etc/init.d/wgsservice
rm -rf gwireguard*
wget -P /root/ https://raw.githubusercontent.com/aau-network-security/gwireguard/master/config/config.yml

echo '#!/sbin/openrc-run
command="/etc/init.d/wgsservice start"
command_background="yes"
output_log="/var/log/wg-service"
' > /run/wgsservice
export CONFIG_PATH=/root/config.yml
rc-update add wgsservice
rc-update add sshd
/etc/init.d/sshd start