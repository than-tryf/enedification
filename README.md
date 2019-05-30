# Description

This program installs the main ENEDI monitoring infrastructure on the premises of a Datacenters. The monitoring infrastructure consists of a netdata agent installed on the machine, a consul server, the Prometheus monitoring system, InfluxDB for persistency, Prometheus-Influx adapter that connects Prometheus with the Influx, and Telegraph that visualizes data from influx db. The program is also responsible for the proper configuration of the system.

## Prerequisites

### Server Requirements

The operating system of the server should be Ubuntu 16.04 64-bits. The server should have installed the following:

* Go Programming language 1.11.5
* Docker 
* Docker-Compose
* Public Access to Internet 
* Git is installed

### Install Go

To install GO please follow the following steps:

```
sudo curl -O https://dl.google.com/go/go1.11.5.src.tar.gz

sudo tar -C /usr/local -xzf go$VERSION.$OS-$ARCH.tar.gz

export PATH=$PATH:/usr/local/go/bin

sudo mkdir $HOME/goimports

export GOPATH=$HOME/goimports

export GOROOT=/usr/local/go

```

To verify that go is installed please run the following command:
```
    go -v
```


### Install Docker 

```
sudo apt-get update

sudo apt-get install \
    apt-transport-https \
    ca-certificates \
    curl \
    gnupg-agent \
    software-properties-common
    
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -

sudo apt-key fingerprint 0EBFCD88

sudo add-apt-repository \
   "deb [arch=amd64] https://download.docker.com/linux/ubuntu \
   $(lsb_release -cs) \
   stable"
   
sudo apt-get update

sudo apt-get install docker-ce docker-ce-cli containerd.io

sudo groupadd docker

sudo usermod -aG docker $USER
```

Then you will have to logout and login again


### Install Docker Compose
```
sudo curl -L "https://github.com/docker/compose/releases/download/1.23.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose

sudo chmod +x /usr/local/bin/docker-compose

```

Confirm that docker-compose is installed by typing
```
docker-compose
```

## How to run

1. Change the working directory to point at the newly created $HOME/goimports ```cd $HOME/goimports```
2. Check if there exists a directory named src. If not create it and cd into it

```
mkdir src
cd src
```
3. Download and install Consul API dependency
```
go get github.com/hashicorp/consul/api
``` 
4. Clone the repo 
```
git clone https://github.com/than-tryf/enedification.git
```   

5. Change directory to point at $HOME/goimports/src/enedification/datacenter-agent and run the script build-with-docker.sh
```bash
cd $HOME/goimports/src/enedification/datacenter-agent
./build-with-docker.sh 

```

6. Run the produced excutable with root privilleges

```bash
sudo ./datacenter-agent
```

Wait the agent to finish installation. After the installation finishes you can check the installation.

| Service        | IP           | 
| ------------- |:-------------:| 
| Prometheus    | <Public_IP>:9090 | 
| Netdata     | <Public_IP>:19999  |
| Consul | <Public_IP>:8500      |  
| InfluDB Telegraph | <Public_IP>:8888      |  


