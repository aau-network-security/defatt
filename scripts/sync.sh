#!/bin/bash

rsync -Pav -e "ssh -p 5555" --exclude 'packer' --exclude 'go.mod' /Users/rvm/Downloads/AAUJOB/NAP/2021/defatt/* ubuntu@130.226.98.173:/home/ubuntu/vlad/sec03/defatt