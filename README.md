# public-clouds-ip-ranges-mmdb
#### TL;DR:

Contains daily compilation of AWS public IP ranges in MaxMind MMDB format.

#### Intro

MaxMind prodives geo ip lookup data in custom [MMDB](https://support.maxmind.com/hc/en-us/articles/4408216157723-Database-Formats) format which is exteremly quick. MMDB format specification is [open](https://maxmind.github.io/MaxMind-DB/) and everyone is encouraged to create its own MMDB files or enrich existing with custom data. MaxMind provides a [Go module](https://pkg.go.dev/github.com/maxmind/mmdbwriter) to write MMDB files and a couple of [examples](https://github.com/maxmind/mmdb-from-go-blogpost).

#### AWS

Amazon (AWS) publishes (and updates) their own IP ranges in a JSON file here: https://ip-ranges.amazonaws.com/ip-ranges.json. 
I needed to have this file in a MMDB format to be able to perform quick lookups against AWS IPs and get info about AWS regions and services.

Publishing a Go code to download JSON, parse it and convert into MMDB.

On top of it, I'm running a Lambda that does the same convertion and uploads resulting MMDB in a public S3 bucket. Lambda runs **every 3 hours** and saves MMDB file at: https://public-clouds-ip-ranges-mmdb.s3.amazonaws.com/aws.mmdb. Feel free to grab the file directly from the bucket if you don't want to complile and regularly run Go code. You can check last update timestamp by listing the bucket content (https://public-clouds-ip-ranges-mmdb.s3.amazonaws.com/)
