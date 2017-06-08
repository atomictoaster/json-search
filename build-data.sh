#!/bin/bash

filename=$1; shift
max=$1; shift

function word() {
    # /usr/share/dict/words is part of the 'words' package
    dict=/usr/share/dict/words;
    rand=$[($RANDOM % `wc -l $dict|sed "s/[^0-9]//g"`)+1];
    sed $rand"q;d" $dict
}
echo "[" > $filename

for record in $(seq 1 $max); do
cat<<EOF>>$filename
  {
    "_id": ${record},
    "url": "http://$(word).$(word).com/api/v2/organizations/${record}.json",
    "external_id": "$(cat /proc/sys/kernel/random/uuid)",
    "name": "$(word)",
    "domain_names": [
      "$(word).com",
      "$(word).com",
      "$(word).com",
      "$(word).com"
    ],
    "created_at": "$(date -Ins)",
    "details": "$(word)",
    "shared_tickets": false,
    "tags": [
      "$(word)",
      "$(word)",
      "$(word)",
      "$(word)"
    ]
  }
EOF

if [ $record != $max ]; then
    echo "    ," >> $filename
fi

done

echo "]" >> $filename
