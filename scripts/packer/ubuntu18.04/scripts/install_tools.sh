#!/bin/bash -eux

## install wireguard
add-apt-repository ppa:wireguard/wireguard
apt-get update -y
apt-get install wireguard -y
apt-get install wget -y
systemctl enable wg-quick@wg0


## install zip and unzip
apt-get install zip  -y
apt-get install unzip -y

## install wireguard gRPC service
wget https://github.com/aau-network-security/gwireguard/releases/download/v1.0.3/gwireguard_1.0.3_linux_64-bit.zip
unzip gwireguard_1.0.3_linux_64-bit.zip && mv gwireguard_1.0.3_linux_64-bit/wgsservice /home/vagrant/wg-service
chmod +x /home/vagrant/wg-service
rm -rf gwireguard*
wget -P /home/vagrant/ https://raw.githubusercontent.com/aau-network-security/gwireguard/master/config/config.yml


