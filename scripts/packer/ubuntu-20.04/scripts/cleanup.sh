#!/bin/bash -eux

# Apt cleanup.
apt autoremove
apt update

# Delete unneeded files.
rm -f /home/vagrant/*.sh
rm -rf /home/vagrant/netman/netman_1.0.2_linux_64-bit.zip /home/vagrant/netman/netman_1.0.2_linux_64-bit
rm -rf /home/vagrant/gip_*.zip

# Zero out the rest of the free space using dd, then delete the written file.
dd if=/dev/zero of=/EMPTY bs=1M
rm -f /EMPTY

# Add `sync` so Packer doesn't quit too early, before the large file is deleted.
sync
