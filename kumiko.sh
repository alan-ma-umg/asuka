#! /usr/bin/env bash

killall kumiko >/dev/null 2>/dev/null

echo 'Welcome! kumiko tcp filter server'
git reset HEAD --hard --quiet && git pull --rebase --quiet
if [ $? -ne 0 ];then
    echo 'Update failed!'
    exit 1;
fi

echo 'Rebuild Kumiko!'
go build -ldflags "-s -w" -o kumiko

sleep 5

nohup ./kumiko  -bloomFilterServer ":17654" -bloomFilterClient "" &