# Estuary CID to Blockstore

A standalone job that queries Estuary CIDs and pulls them to a specified local blockstore.

## Installation
### Create the DB connection .env file

```
DB_NAME=
DB_HOST=
DB_USER=
DB_PASS=
DB_PORT=
```

### Install run
```
// Install
go mod tidy
go mod download
go build -tags netgo -ldflags '-s -w' -o cid-to-bs

// Run
./cid-to-bs
```


### Output

Number of workers: 10
Each worker is assigned with a group of CID(s) to pull
```
worker 10 started  job 1 cid {{bafkreic4ykm2yqwayjphykisv2vgugsf4oqxqw25mn6eonmi43zsvt3bui}}
worker 5 started  job 2 cid {{bafkreib5jsrjvkalimteyrs5fsoi4e6s6r4tdbea33hg27vjdjozusee3q}}
worker 7 started  job 3 cid {{bafkreigqomn42blxrsf67sd3gxenck46sfm5smo2m7h6cttdf7e225e7nq}}
worker 6 started  job 4 cid {{bafkreibsbq7rtux5bbzertabqzl56iiv2kxnf6vjy6yfxjfeb4oyi2c2jm}}
worker 1 started  job 5 cid {{bafybeihn6vn4ppkpck3ffiv2xsdvjzn3kampwizcqqupdf662fuixrc6ny}}
worker 8 started  job 6 cid {{bafybeibcayare5lx4y5mv4npupzyaodkmmvex3kjzlil72qap7z3m65oge}}
worker 2 started  job 7 cid {{bafkreifmaz57r7qgtnxxbtyzqjyjh2kb5lgvukge7jczfualoi3hxute4e}}
worker 3 started  job 8 cid {{bafy2bzaceb4iu6mhxrzvnw6rvlvnpggkr4n7cp5e6fmvo4zd25gnn3uagsw4o}}
worker 4 started  job 9 cid {{bafy2bzacec6hpnxarhoj2uqa2ssx677pq5ip3k3jm5nrglzkfii4vtvujqkco}}
worker 9 started  job 10 cid {{bafkreifr2mjxhtj7b4izefndzokmhzj7mi4yvhaiidk6xv5fl46vx2mtby}}
```

