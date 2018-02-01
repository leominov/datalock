#!/bin/bash

SRCDOMAIN='http://1seasonvar.ru'
DSTDOMAIN='http://127.0.0.1:7000'

wget -q $SRCDOMAIN/sitemap.xml --no-cache -O - | egrep -o "$SRCDOMAIN[^<]+" | while read subsite;
do
	subsite="$DSTDOMAIN"$(echo "$subsite" | egrep -o "/serial-.*")
	echo --- Reading serial info: $subsite: ---
	curl -sI -X GET "$subsite"
	echo --- FINISHED reading serial info: $subsite: ---
done
