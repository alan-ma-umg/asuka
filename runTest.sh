#! /usr/bin/env bash

echo 'Welcome! Asuka'
git pull --rebase  --quiet && killall asuka >/dev/null 2>/dev/null
echo 'Killed Asuka!'


echo 'Rebuild Asuka!'
go build asuka

nohup ./asuka env.json & >/dev/null 2>/dev/null
echo 'Now, Asuka is alive .'

sleep 1
tail -f nohup.out
