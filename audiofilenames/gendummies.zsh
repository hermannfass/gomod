#!/bin/zsh
dir=~/audiofiles
rm $dir/*
for i in 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16
# for i in   1 2 3   5 6 7 8 9 10 11    13 14 15 16
do
	fn=`printf "TRACK%02d.WAV" $i`
	touch $dir/$fn
done

