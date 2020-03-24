#!/usr/bin/env bash


declare -i RC=0
declare -i rc

cd "$(dirname $0)"
for D in $(ls)
do
    [ ! -f "$D/main.go" ] && continue
    echo $D
    pushd $D > /dev/null
    GOOS=linux GOARCH=arm GOARM=7 go build
    rc=$?
    popd > /dev/null
    RC=$(( RC + rc ))
done

exit $RC