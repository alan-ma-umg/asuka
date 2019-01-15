#! /usr/bin/env bash

echo 'Killed Asuka!'
killall asuka >/dev/null 2>/dev/null

echo 'Rebuild Asuka!'
go build asuka

nohup ./asuka env.json &
echo 'Now, Asuka is alive .'

