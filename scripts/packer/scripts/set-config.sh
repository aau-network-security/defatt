#!/usr/bin/env sh

echo '%wheel ALL=(ALL) NOPASSWD:ALL' > /etc/sudoers.d/wheel
user={{user `ssh_username`}}
echo "Add user $user with NOPASSWD sudo"
adduser $user --disabled-password
echo '{{user `ssh_username`}}:{{user `ssh_password`}}' | chpasswd
adduser $user wheel
echo add ssh key
cd ~{{user `ssh_username`}}
mkdir .ssh
chmod 700 .ssh
echo {{user `ssh_key`}} > .ssh/authorized_keys
chown -R $user .ssh
echo disable ssh root login
sed '/PermitRootLogin yes/d' -i /etc/ssh/sshd_config
