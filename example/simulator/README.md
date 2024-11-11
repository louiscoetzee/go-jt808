# Performance Testing

macOS Parameters Temporary Modification
``` shell
# Increase system default max connection limit
sudo sysctl -w kern.maxfiles=8880000
# Increase default max connection limit per process
sudo sysctl -w kern.maxfilesperproc=8990000
# Set the maximum number of files the current shell can open
ulimit -n 1100000
# Set the maximum number of user threads the current shell can create
ulimit -u 26600
# Adjust the range of available ports
sysctl -w net.inet.ip.portrange.first=5000
```

Linux System Permanent Modification
``` shell
vi /etc/sysctl.conf
# System-level parameter for maximum open files
fs.file-max=1100000
# Process-level parameter for maximum open files
fs.nr_open=1100000
# Increase available process ports
net.ipv4.ip_local_port_range = 5000 65000

vi /etc/security/limits.conf
# Set max open files limit for user processes - soft (soft limit) and hard (hard limit)
*  soft  nofile  1000000
*  hard  nofile  1010000

# Apply changes
sysctl -p
# View settings
sysctl -a
```

---
<h2 id="online"> 1. Connection Test </h2>

- Simulator sends 1 registration, 1 authentication, and cycles heartbeat (20-second interval)
- Using 2 cloud servers [2 cores, 4 GB] since one IP has a valid port range of 65535; testing 100k+ requires 2 IPs

### 1.1 Steps
Server
``` shell
cd ../quick_start && go build
nohup ./start >./start.log &
```

Simulator
``` shell
cd ./client && go build
# Simulator on Server A
nohup ./client -ip=127.0.0.1 -max=50005 -blc=0 -lc=0 >./test1.log &
# Simulator on Server B
nohup ./client -ip=10.0.16.14 -addr=10.0.16.5:8080 -max=50005 -blc=0 -lc=0 >./test2.log  &
```

### 1.2 Test Results
| Server Version | Scenario  | Concurrency | Server Config | Server Resource Usage | Description |
| :---:   | :-------: | :--: | :------: | :-------------- | :----------------------------: |
| v0.3.0 | Connection Test | 100k+ | 2 cores, 4 GB | 120%+ CPU, 1.7 GB memory | Service and simulator running on 10.0.16.5, <br/> simulator on 10.0.16.14 |

Tencent Cloud statistics *2 on local machine shows 150017, actual 100012
![Tencent Cloud Monitoring](./testdata/tx.png)

<h2 id="save"> 2. Latitude and Longitude Storage Simulation Test </h2>

- Simulator sends 1 registration, 1 authentication, and cycles heartbeat and location reports
- Default intervals are 20 seconds and 5 seconds
- Message queue uses NATS, database uses TDengine

### 2.1 Required Services: NATS and TDengine
| Description         | Link                         |
|---------------------|------------------------------|
| NATS                | https://github.com/nats-io/nats-server |
| Quick Install TDengine | https://docs.taosdata.com/get-started/package/ |

### 2.2 Steps

Simulate storage of latitude and longitude by retrieving from NATS and saving to TDengine
``` shell
cd ./save && GOOS=linux GOARCH=amd64 go build
# Receive data, each terminal has a table named (T+phone number)
./save -nats=127.0.0.1:4222 -dsn='root:taosdata@ws(127.0.0.1:6041)/information_schema' >./save.log
```

Server sends latitude and longitude packets to NATS
``` shell
cd ./server && GOOS=linux GOARCH=amd64 go build
./server -nats=127.0.0.1:4222 >./server.log
```

Simulator simulates device sending latitude and longitude
``` shell
cd ./client && GOOS=linux GOARCH=amd64 go build
# Open up to 10,000 clients, each sends 10,000 0x0200 latitude and longitude messages, totaling 1 billion locations (for ease, no 0x0704 sent)
./client -ip=127.0.0.1 -blc=0 -limit=10000 -max=10000 -lc=2 >./client.log
```

Resource consumption statistics
``` shell
# Save status every 1800 seconds
sudo atop -w ./atop.log 1800
# View each service’s resource usage
# t-advance, T-backtrack, m-memory, g-cpu, c-details
atop -r ./atop.log
```

### 2.3 Test Results
``` sql
# On TDengine server, enter taos to access the command line
select count(*) from power.meters;
# View saved data per simulator terminal
select tbname, count(*) from power.meters group by tbname order by count(*) >> /home/test2/td.log;
```
![Data Storage Overview](./testdata/db.png)

10,000 clients, each sends 100 0x0200 messages
| Sent Count | Success Rate | Frequency | Description |
| :---: | :-----: | :------: | :------: |
| 1 million |  98.5%+ | 10,000 per second | Three tests: 985430, 985299, 989668 |
| 1 million | 100% | 5,000 per second | Three tests: 1 million, 1 million, 1 million |

| Server Version | Scenario | Clients | Server Config | Service Resource Usage | Description |
| :---: | :-------: | :--: | :------: | :-------------- | :----------------------------: |
| v0.3.0 | Real Scenario Simulation | 10,000 | 2 cores, 4 GB | 35% CPU, 180.4 MB memory | 5,000 per second, 1 billion locations saved <br/> Actual saved: 99,999,174, success rate: 99.999% |

- Data loss due to save process channel queue overflow <br/>
![Data Loss](./testdata/save.png)

Resource usage of each service
- [atop sample details](./testdata/atop.log)
- [CPU usage by atop](./testdata/atop_cpu.png)
- [Memory usage by atop](./testdata/atop_cpu.png)

| Service        | CPU   | Memory    | Description      |
| :-------------:| :---: | :-------: | :--------------: |
| server         | 35%   | 180.4 MB  | 808 server       |
| client         | 23%   | 196 MB    | Simulation client |
| save           | 18%   | 68.8 MB   | Data storage service |
| nats-server    | 20%   | 14.8 MB   | Message queue    |
| taosadapter    | 37%   | 124.3 MB  | TDengine adapter |
| taosd          | 15%   | 124.7 MB  | TDengine database |
