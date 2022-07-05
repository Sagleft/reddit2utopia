#!/usr/bin/bash

scripts=("conspiracy" "facepalm" "news" "space" "todayilearned")

cd bin

for scriptname in "${scripts[@]}"
do
    sh $scriptname".sh"
done
