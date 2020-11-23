#!/usr/bin/env bash

for vm in `vboxmanage list vms | awk -F'"' '$0=$2'`
do
    echo "Rules for VM $vm"
    VBoxManage showvminfo $vm --machinereadable | awk -F '[",]' '/^Forwarding/ { printf ("Rule %s host port %-5d forwards to guest port %-5d\n", $2, $5, $7); }'
    printf '\n'
done
