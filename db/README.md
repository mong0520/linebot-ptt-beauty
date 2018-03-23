# Fetch raw data

- Clone repo: https://github.com/mong0520/ptt-web-crawler (fork from https://github.com/jwlin/ptt-web-crawler)
- execute `run.sh ${PAGE_OFFSET}` to generate `Beauty.json`, 


#Set Unique Index ID

```
mongo
use ptt
db.beauty.createIndex( { "article_id": 1 }, { unique: true } )
```


# Import raw data to MongoDB (minimum verion 3.2 is required)

> mongoimport --db ptt --collection beauty --type json --file Beauty.json --jsonArray --mode merge --upsertFields  article_id
