# balboa-backends

`balboa-backends` is the central repository for database backends leveraged by
[balboa](https://www.github.colm/DCSO/balboa)

## Building and Installation

### RocksDB Backend

First, install the required depencencies. Aside from `make` and a working `gcc`
installation `rocksdb` development files are necessary.

On Debian (testing and stretch-backports), one can satisfy these dependencies
with:

```
apt install librocksdb-dev
```

On Void Linux

```
xbps-install rocksdb-devel
```

Building the RocksDB backend:

```
cd backend/balboa/balboa-rocksdb
make
```

This yields the backend binary `build/linux/balboa-rocksdb`.

### Usage

Show available parameters to the RocksDB backend:

```
balboa-rocksdb -h
```

Now start *balboa* and the backend to feed pDNS observations into it:

```
balboa-rocksdb -p 4242 -h 127.0.0.1 -d /tmp/balboa-rocksdb -D
balboa serve -l '' -host 127.0.0.1:4242 -f my-feeders.yaml
```

### Other tools

#### balboa-backend-console

`balboa-backend-console` is a small utility managing balboa backends. It speaks
the *backend protocol* directly. You can `backup`, `dump` and `replay`
databases. Building is as easy as:

```
cd backend/balboa-backend-console
make
```

assuming `make` and `gcc` are working on your system. It drops the self
contained binary `build/linux/balboa-backend-console`.

```
balboa-backend-console -h
```

#### balboa-rocksdb-v1-dump

`balboa-rocksdb-v1-dump` is a small utility dealing with migration from
*balboa* v1 databases to the new and internal format of v2.

```
cd backend/balboa-rocksdb-v1-dump
make
```

if successfull the binary is located under
`build/linux/balboa-rocksdb-v1-dump`.

## Migrating from balboa/rocksdb v1

`balboa` v1 uses a different format to store data in the `rocksdb` backend. In
order not to loose all of our precicous pDNS observations we need a migration
strategy. The suggested approach looks a bit cumbersome but is in part due to
the reason `rocksdb` does not yet support hot backup.

Let's assume you have a moderatly large pDNS observation database stored in
`/data/balboa-rocksdb-v1`.

```
du -s /data/balboa-rocksdb-v1
42GB
```

Now, we want to migrate this database to balboa v2. The feeder configuration
remains the same. Start the new balboa backend and frontend.

```
balboa-rocksdb -d /data/balboa-rocksdb-v2 -h 127.0.0.1 -p 4242
balboa serve -l '' -f my-feeders.yaml -host 127.0.0.1:4242
```

Stop the old balboa v1 service and dump the database:

```
balboa-rocksdb-v1-dump dump /data/balboa-rocksdb | lz4 > /data/pdns-backup.dmp.lz4
lz4cat /data/pdns-backup.dmp.lz4 | balboa-backend-console replay -h 127.0.0.1 -p 4242
```

`babloa-rocksdb-v1-dump` will dump the observation entries to stdout which in
turn gets redirected to the `balboa-backend-console` tool in `replay` mode
(reading from stdin in this case, use -d <path> to read from dump file)

Wait some time. Done.

## Author/Contact

- Sascha Steinbiss
- Alexander Wamack

## License

BSD-3-clause