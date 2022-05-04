#!/bin/bash

rsync -Pav --exclude 'config.yml' --exclude 'scenario.yml' --exclude "ansible" --exclude "packer" --exclude "/scripts" --exclude "sync3.sh" /Users/rvm/Downloads/AAUJOB/NAP/nap22/defatt/* vlad@sec03.lab.es.aau.dk:/home/vlad/2022/test/defatt
