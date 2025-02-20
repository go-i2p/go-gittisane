#! /bin/sh
nohup /startapp.sh 2> err.log 1> log.log &
export countdown=30
while [ countdown -gt 0 ]; do
    echo "Waiting for app to start... $countdown seconds left"
    countdown=$((countdown - 1))
    sleep 1
done
ls /usr/local/bin
/usr/local/bin/gittisane web
