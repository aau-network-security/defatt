#!/usr/bin/env sh

apk update
apk add --no-cache sudo
apk add -U wireguard-tools
apk add --no-cache openrc-doc
apk add --no-cache iptables
apk add --no-cache virtualbox-guest-additions virtualbox-guest-modules-virt
apk add --no-cache wget
apk add --no-cache unzip

wget https://github.com/aau-network-security/gwireguard/releases/download/v1.0.3/gwireguard_1.0.3_linux_64-bit.zip
unzip gwireguard_1.0.3_linux_64-bit.zip && mv gwireguard_1.0.3_linux_64-bit/wgsservice /etc/init.d/
chmod +x /etc/init.d/wgsservice
rm -rf gwireguard*
wget -P /root/ https://raw.githubusercontent.com/aau-network-security/gwireguard/master/config/config.yml

echo "#!/sbin/openrc-run
command=/usr/bin/wgsservice
" > /run/wgsservice
echo "export CONFIG_PATH=/root/config.yml" >> /etc/profile
