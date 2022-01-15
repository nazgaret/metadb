Metadb Administrator Guide
==========================

##### Contents  
1\. [Overview](#1-overview)  
2\. [System requirements](#2-system-requirements)  
3\. [Installing Metadb](#3-installing-metadb)  
4\. [Running the server](#4-running-the-server)  
5\. [Running the client](#5-running-the-client)  
6\. [Adding an analytic database](#6-adding-an-analytic-database)  
7\. [Adding a data source](#7-adding-a-data-source)  
8\. [Resynchronizing a data stream](#8-resynchronizing-a-data-stream)


1\. Overview
------------

The software consists of a server (`metadb`) and a command-line client
(`mdb`).

The `metadb` server is used to create a new instance, e.g. with a data
directory `metadb_data`:

```bash
$ metadb init -D metadb_data
```

Once this is done, the server can be started:

```bash
$ metadb start -D metadb_data
```

The client can then be used to configure the instance, for example by 
adding a database connector:

```bash
$ mdb config db.main.type postgresql
$ mdb config db.main.host dbserver
$ mdb config db.main.dbname metadb
$ mdb config db.main.adminuser metadbadmin
$ mdb config --pwprompt db.main.adminpassword
$ mdb config db.main.sslmode require
$ mdb enable db.main
```


2\. System requirements
-----------------------

* Architecture:  x86-64 (AMD64)
* Operating system:  Ubuntu Linux 22.04 LTS
* Database systems supported:
  * [PostgreSQL](https://www.postgresql.org/) 14.1 or later
* Other software dependencies:
  * [Go](https://golang.org/) 1.18 or later
  * [GCC C compiler](https://gcc.gnu.org/) 9.3.0 or later


3\. Installing Metadb
---------------------

### Branches

There are two primary types of branches:

* The main branch (`main`).  This is a development branch where new 
  features are first merged.  It is relatively unstable.  Note that it 
  is also the default view when browsing the repository in GitHub.

* Release branches (`release-*`).  These are releases made from 
  `main`.  They are managed as stable branches; i.e. they may receive 
  bug fixes but generally no new features.  Most users should pull 
  from a recent release branch.

### Building the software

First set the `GOPATH` environment variable to specify a path that can
serve as the build workspace for Go, e.g.:

```bash
$ export GOPATH=$HOME/go
```

Then:

```bash
$ ./build.sh
```

The `build.sh` script creates a `bin/` subdirectory and builds the `metadb` and
`mdb` executables there:

```bash
$ ./bin/metadb help
```
```bash
$ ./bin/mdb help
```


4\. Running the server
----------------------

### Metadb user account

It is recommended to create a user account that will run the `metadb` server.
This guide assumes that a user "metadb" has been created for this purpose.

### Creating a data directory

Metadb stores an instance's state and metadata in a "data directory,"
which is created using the `init` command of `metadb`.  The data
directory is required to be on a local file system; it may not be on a
network file system.

As root:
```bash
# chown metadb /usr/local/metadb
```
As metadb:
```
$ metadb init -D /usr/local/metadb/data
metadb: initializing new instance in /usr/local/metadb/data
```

If the directory already exists, `metadb` will exit with an error.

Creating the data directory is generally a one-time operation.

Note that the data directory contains important data and there is
currently no function to regenerate it.  It must be kept in sync with
the analytic database(s) and they should be backed up together.  See
the section on backups below for details.

### Starting the server

To start the server:

```bash
$ nohup metadb start -D /usr/local/metadb/data -l metadb.log &
```

The server log is by default written to standard error.  The `-l` or 
`--log` option specifies a log file.  The `--csvlog` option writes a 
log in CSV format.

The server listens by default on the loopback address.  This allows 
the command-line client, `mdb`, to connect to the server when running 
locally.

The server uses two ports:

* The "admin port" defaults to 8440.  This provides administrative 
  services, e.g. server configuration.

* The "client port" defaults to 8441.  It is currently disabled but is 
  planned to support metadata services.

The `--listen` option allows listening on a specified address.  When
`--listen` is used, the `--cert` and `--key` options also must be
included to provide a server certificate (including the CA's
certificate and intermediates) and matching private key.  As an
alternative, `--notls` may be used in both server and client to
disable TLS entirely; however, this is insecure and for testing
purposes only.

The `--debug` option enables detailed logging.

### Stopping the server

To stop the server:

```bash
$ metadb stop -D /usr/local/metadb/data
```

It is recommended to stop the server before making a backup of the
data directory.

<!--

### Upgrading to a new version

Please note that this section documents an upgrade feature that is
planned for Metadb 1.2 and not yet available.

When installing a new version of Metadb, the instance should be
"upgraded" before starting the new server:

1. Stop the old version of the server.

2. Make a backup of the data directory and database(s).  See the 
   section below on backups to make sure this is done correctly.

3. Use the `upgrade` command in the new version of Metadb to perform 
   the upgrade, e.g.:

```bash
$ metadb upgrade -D /usr/local/metadb/data
```

4. Start the new version of the server.

In automated deployments, the `upgrade-database` command can be run
after `git pull`, whether or not any new changes were pulled.  If no
upgrade is needed, it will exit normally:

```bash
$ metadb upgrade -D /usr/local/metadb/data ; echo $?
metadb: instance is up to date
0
```

-->


5\. Running the client
----------------------

By default the `mdb` client tries to connect to the server on the
loopback address.  To specify a different address, the `-h` or
`--host` option may be used, which also will enable TLS unless
`--notls` is included.

The `-v` option enables verbose output.


6\. Adding an analytic database
-------------------------------

(To be written)

Example:

```bash
$ mdb config db.main.type postgresql
$ mdb config db.main.host dbserver
$ mdb config db.main.dbname metadb
$ mdb config db.main.adminuser metadbadmin
$ mdb config db.main.adminpassword @metadb_creds.txt
$ mdb config db.main.sslmode require
$ mdb enable db.main
```


7\. Adding a data source
------------------------

(To be written)

Example:

```bash
$ mdb config src.example.brokers kafka:29092
$ mdb config src.example.topics '^metadb_example[.].*'
$ mdb config src.example.group metadb_example
$ mdb config src.example.schemapassfilter 'example_.+'
$ mdb config src.example.schemaprefix 'example_'
$ mdb config src.example.dbs main
$ mdb enable src.example
```


8\. Resynchronizing a data stream
---------------------------------

If a Kafka data stream fails and cannot be resumed, it may be
necessary to re-stream data to Metadb.  For example, a source database
may become unsynchronized with the analytic database, requiring a new
snapshot of the source database to be streamed.  Metadb can accept
re-streamed data in order to resynchronize with the source, using this
procedure:

1. Disable the source connector, for example:

```bash
$ mdb disable src.example
```

2. Update the source connector's `topics` and `group` configuration
   settings for the new data stream, or temporarily delete them or set
   them to empty strings.

```bash
$ mdb config src.example.topics ''
$ mdb config src.example.group ''
```

3. Stop the Metadb server.

4. "Reset" the analytic database to mark current data as old.  This
   may take some time to run.

```bash
$ metadb reset -D /usr/local/metadb/data --origin 's1,s2,s3' db.main
```

Note that `--origin` should include all origins associated with the
source, or empty string ( `''` ) if there are no origins.

5. Start the Metadb server, configure and enable the source connector
   for the new stream, and begin streaming the data.

6. Once the new data have finished or nearly finished streaming, stop
   the Metadb server, and "clean" the analytic database to remove old
   data.

```bash
$ metadb clean -D /usr/local/metadb/data --origin 's1,s2,s3' db.main
```

The `--origin` should be the same as used for the reset.

7. Start the server.

If a failed stream is re-streamed without following the process above,
then the analytic database may become unsynchronized with the source.
