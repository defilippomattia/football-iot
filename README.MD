# Introduction

For learning purposes. Technology used: go, nats, grafana, mongo, postgres, hashicorp vault, docker

# Architecture

**Player 1...N** - csv files that are simulating sensor data  
**publisher.go** - reads csv files and publishes data to football-iot.players NATS subject  
**NATS** - messaging system  
**subscriber_raw.go** - subscribes to football-iot.players subject and stores unmodified (raw) data to postgres db  
**subscriber_modifed.go** - subscribes to football-iot.players subject and stores modified data to mongo db  
**mongo** - nosql db  
**postgres** - relatioanl db (used because postgres data source in grafana is free)  
**grafana** - visualisation tool  
**hashicorp vault** - storing secrets (not implemented yet)  

![alt text](https://github.com/defilippomattia/football-iot/blob/main/resources/other/architecture_diagram.png)

## How to run

```
git clone https://github.com/defilippomattia/football-iot.git  
docker-compose up -d 
go run subscriber_raw.go
go run subscriber_modified.go
go run publisher.go
```
### Some notes
grafana at localhost:3000   
(credentials:admin & admin)  
pgadmin at localhost:5050  
(credentials: hostname for pg server: host.docker.internal, u:admin@admin.com p:root)  
mongoexpres at localhost:8082

