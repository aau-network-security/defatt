#!/bin/bash


# todo: will be modified with more dynamic way of cleaning stuff

#request sudo....
#if [[ $UID != 0 ]]; then
#    echo "Please run this script with sudo:"
#    echo "sudo $0 $*"
#    exit 1
#fi

sudo ovs-vsctl del-br game
sudo ip tuntap del tap0 mode tap
sudo ip tuntap del tap1 mode tap
sudo ip tuntap del tap2 mode tap
sudo ip tuntap del tap3 mode tap
sudo ip tuntap del tap4 mode tap
sudo ip tuntap del vlan10 mode tap
sudo ip tuntap del vlan20 mode tap
sudo ip tuntap del vlan30 mode tap
sudo ip tuntap del mon10 mode tap
sudo ip tuntap del ALLblue mode tap


VBoxManage list runningvms | awk '/nap/ {print $1}' | xargs -I vmid VBoxManage controlvm vmid poweroff
VBoxManage list vms | awk '/nap/ {print $2}' | xargs -I vmid VBoxManage unregistervm --delete vmid


rm -rf ~/VirtualBox\ VMs/nap

#while read -r line; do
 #   vm=$(echo $line | cut -d ' ' -f 2)
  #  echo $vm
  #  vboxmanage controlvm $vm poweroff
   # vboxmanage unregistervm $vm --delete
#done <<< "$VMS"



# Remove all docker containers that have a UUID as name
#docker ps -a --format '{{.Names}}' | grep -E '[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}' | xargs docker rm -f

docker kill $(docker ps -q -a -f "label=nap")

docker rm $(docker ps -q -a -f "label=nap")




# Remove all macvlan networks
docker network rm $(docker network ls -q -f "label=defatt")

# Prune entire docker
docker system prune --filter "label=nap"

# Prune volumes
docker volume prune --filter "label=nap"
