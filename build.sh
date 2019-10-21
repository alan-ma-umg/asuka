#! /usr/bin/env bash

killall asuka >/dev/null 2>/dev/null

echo 'Welcome! Asuka'
#git reset HEAD --hard --quiet && git pull --rebase --quiet
#if [ $? -ne 0 ];then
#    echo 'Update failed!'
#    exit 1;
#fi

echo 'Rebuild Asuka!'
go build

sleep 2
echo 'Now, Asuka is alive .'
nohup ./asuka  -listen :666 &
