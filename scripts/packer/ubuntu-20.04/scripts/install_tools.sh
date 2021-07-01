#!/bin/bash -eux

## install wireguard
apt-get update -y
apt-get install wireguard -y
apt-get install wget -y

## configure port forwarding ...
modprobe wireguard
systemctl enable wg-quick@wg0
echo "net.ipv4.ip_forward=1" >> /etc/sysctl.conf
sysctl -p



# install net tools like ifconfig
apt-get install net-tools -y
apt-get install ifupdown -y

## install zip and unzip
apt-get install zip  -y
apt-get install unzip -y

## install netman service to manage down network interfaces
## pop up version if required
mkdir /home/vagrant/netman && cd /home/vagrant/netman
wget https://github.com/mrturkmenhub/netman/releases/download/1.0.4/netman_1.0.4_linux_64-bit.zip
unzip netman_1.0.4_linux_64-bit.zip && mv netman_1.0.4_linux_64-bit/* /home/vagrant/netman/
chmod +x /home/vagrant/netman/netman
cp /home/vagrant/uploads/netman.service /etc/systemd/system/
systemctl daemon-reload
systemctl enable netman.service

## install git
apt-get install git-all -y

## install wireguard gRPC service
cd /home/vagrant
wget https://github.com/aau-network-security/gwireguard/releases/download/v1.0.5/gwireguard_1.0.5_linux_64-bit.zip
unzip gwireguard_1.0.5_linux_64-bit.zip && mv gwireguard_1.0.5_linux_64-bit/wgsservice /home/vagrant/wg-service
chmod +x /home/vagrant/wg-service
rm -rf gwireguard*
wget -P /home/vagrant/ https://raw.githubusercontent.com/aau-network-security/gwireguard/master/config/config.yml

## install gip gRPC service
cd /home/vagrant
wget https://github.com/aau-network-security/gip/releases/download/1.0.0/gip_1.0.0_linux_64-bit.zip
unzip gip_1.0.0_linux_64-bit.zip  && mv gip_1.0.0_linux_64-bit gip
chmod +x /home/vagrant/gip/gip
cp /home/vagrant/uploads/gip.service /etc/systemd/system/
systemctl daemon-reload
systemctl enable gip.service


## enable wg-service in system daemon
cp /home/vagrant/uploads/wg-service.service /etc/systemd/system/wg-service.service
sudo chmod 644  /etc/systemd/system/wg-service.service
sudo systemctl start wg-service
sudo systemctl enable wg-service

# install docker and docker compose
sudo apt-get update -y
sudo apt-get install -y apt-transport-https ca-certificates curl  gnupg-agent software-properties-common
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
sudo apt-key fingerprint 0EBFCD88
sudo add-apt-repository \
   "deb [arch=amd64] https://download.docker.com/linux/ubuntu \
   $(lsb_release -cs) \
   stable"
sudo apt-get update -y
sudo apt-get install docker-ce docker-ce-cli containerd.io -y

sudo curl -L https://github.com/docker/compose/releases/download/1.27.4/docker-compose-`uname -s`-`uname -m` -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
sudo usermod -aG docker vagrant

# Will be managed later
#git clone https://github.com/aau-network-security/nap-monitoring.git
#cd /home/vagrant/nap-monitoring/
#cp /home/vagrant/uploads/monitoring.service /etc/systemd/system/
#systemctl daemon-reload
#systemctl enable monitoring.service

