# nexus-carp

CAS Authentication Reverse Proxy (CARP) for Sonatype Nexus.

## Requirements

* [Go](https://golang.org/) >= 1.12.x

## Testing and Development

* start Cloudogu EcoSystem
* install at least CAS-Dogu
* enable development mode and restart CAS in your Cloudogu EcoSystem
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
* sign in with following credentials:
  * Username: admin
  * Password: read out from Docker container via

    ```docker exec -it nexus-carp_nexus_1 cat /nexus-data/admin.password```
* finish the initialization wizard (remember the password)  
* enable "Rut Auth Realm" (settings at Security -> Realms)
* Add "Rut Auth" Capability with `X-CARP-Authentication` as Header (settings at System -> Capabilities -> Create Capability)
* Add property that allows to add scripts to Nexus and restart container
```
docker exec -it nexus-carp_nexus_1 bash
echo "nexus.scripts.allowCreation=true" >> /nexus-data/etc/nexus.properties
exit
docker-compose restart
```

* Build
```bash
export GO111MODULE=on
make
```

* Set required environment variables (use the password you set in the wizard at the first start)
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

## What is the Cloudogu EcoSystem?
The Cloudogu EcoSystem is an open platform, which lets you choose how and where your team creates great software. Each service or tool is delivered as a Dogu, a Docker container. Each Dogu can easily be integrated in your environment just by pulling it from our registry.

We have a growing number of ready-to-use Dogus, e.g. SCM-Manager, Jenkins, Nexus Repository, SonarQube, Redmine and many more. Every Dogu can be tailored to your specific needs. Take advantage of a central authentication service, a dynamic navigation, that lets you easily switch between the web UIs and a smart configuration magic, which automatically detects and responds to dependencies between Dogus.

The Cloudogu EcoSystem is open source and it runs either on-premises or in the cloud. The Cloudogu EcoSystem is developed by Cloudogu GmbH under [AGPL-3.0-only](https://spdx.org/licenses/AGPL-3.0-only.html).

## License
Copyright © 2020 - present Cloudogu GmbH
This program is free software: you can redistribute it and/or modify it under the terms of the GNU Affero General Public License as published by the Free Software Foundation, version 3.
This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more details.
You should have received a copy of the GNU Affero General Public License along with this program. If not, see https://www.gnu.org/licenses/.
See [LICENSE](LICENSE) for details.


---
MADE WITH :heart:&nbsp;FOR DEV ADDICTS. [Legal notice / Imprint](https://cloudogu.com/en/imprint/?mtm_campaign=ecosystem&mtm_kwd=imprint&mtm_source=github&mtm_medium=link)
