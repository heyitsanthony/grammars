#!/bin/bash

# Run over the entire stdlib
# find /usr/lib/go/src/ -name \*.go | ./bench.sh
while read a; do
	# realtime (1/0)
	dat=`((time ../grammars -grammar go -ast=false 2>/dev/null <$a && echo OK) | wc -l) 2>&1 | xargs echo | awk ' { print $2" "$7 } '`
	echo $a $dat
done
