#! /usr/bin/env bash

echo 'Welcome! Asuka'
git reset HEAD --hard --quiet && git pull --rebase --quiet && go get asuka
if [ $? -ne 0 ];then
    echo 'Update failed!'
    exit 1;
fi

killall asuka >/dev/null 2>/dev/null
echo 'Killed Asuka!'


echo 'Rebuild Asuka!'
go build asuka

echo 'Now, Asuka is alive .'
nohup ./asuka env.json > asuka.out 2>&1 &
