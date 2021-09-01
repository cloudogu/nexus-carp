# nexus-carp

CAS Authentication Reverse Proxy (CARP) for Sonatype Nexus.

## Requirements

* [Go](https://golang.org/) >= 1.12.x

## Testing and Development

* start Cloudogu EcoSystem
* install at least CAS-Dogu
* enable development mode and restart CAS
```bash
etcdctl set /config/_global/stage development
cesapp stop cas
cesapp start cas
```

* Checkout nexus-carp
```bash
git clone git@github.com:cloudogu/nexus-carp.git
cd nexus-carp
```

* start nexus on your host system
```bash
docker-compose up -d
```
* open Nexus at http://localhost:8081
* sign in with admin and admin123
* finish the initialization wizard  
* enable "Rut Auth Realm" (settings at Security -> Realms)
* Add "Rut Auth" Capability with `X-CARP-Authentication` as Header (settings at System -> Capabilities -> Create Capability)


* Build
```bash
export GO111MODULE=on
make
```

* Set required environment variables
```bash
export NEXUS_URL="http://localhost:8081"
export NEXUS_USER="admin"
export NEXUS_PASSWORD="admin123"
export CES_ADMIN_GROUP="cesAdmins"
```

* Run
```bash
./target/nexus-carp
```

* Test Nexus with Browser and Maven at http://192.168.56.1:9090
