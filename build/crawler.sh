#!/usr/bin/bash

scripts=("biden" "conspiracy" "facepalm" "news" "space" "todayilearned")

cd bin

for scriptname in "${scripts[@]}"
do
    sh $scriptname".sh"
done
