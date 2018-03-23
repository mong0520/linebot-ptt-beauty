#!/bin/bash

python PttWebCrawler/crawler.py -b beauty -o $1

# Set Unique Index ID
# db.beauty.createIndex( { "article_id": 1 }, { unique: true } )

# Import
# mongoimport --db ptt --collection beauty --type json --file /data/db/Beauty-2425-2426.json --jsonArray

