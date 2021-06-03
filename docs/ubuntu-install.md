#  Building from source / Ubuntu Installation Guide 

Although personally I feel that using the docker container is the best way of using and enjoying something like Hammond, a lot of people in the community are still not comfortable with using Docker and wanted to host it natively on their Linux servers.

This guide has been written with Ubuntu in mind. If you are using any other flavour of Linux and are decently competent with using command line tools, it should be easy to figure out the steps for your specific distro. 

## Install Go and Node

Hammond is built using Go and VueJS which means GO and Node would be needed to compile and build the source code. Hammond is written with Go 1.15/ Node v14 so any version equal to or above this should be good to go. 

If you already have Go and Node installed on your machine, you can skip to the next step.

Get precise Go installation process at the official link here - https://golang.org/doc/install

Get precise Node installation process at the official link here - https://nodejs.org/en/


Following steps will only work if Go and Node are installed and configured properly.

## Install dependencies

``` bash
 sudo apt-get install -y git ca-certificates ufw gcc
```

## Clone from Git

``` bash
git clone --depth 1 https://github.com/akhilrex/hammond
```

## Build and Copy dependencies

``` bash
cd hammond/server
mkdir -p ./dist
cp .env ./dist
go build -o ./dist/hammond ./main.go
```

## Create final destination and copy executable
``` bash
sudo mkdir -p /usr/local/bin/hammond
mv -v dist/* /usr/local/bin/hammond
mv -v dist/.* /usr/local/bin/hammond
```


## Building the UI

Go back to the root of the hammond folder.

``` bash
cd ui
npm install
npm run build
mv dist /usr/local/bin/hammond 
```

At this point theoretically the installation is complete. You can make the relevant changes in the ```.env``` file present at ```/usr/local/bin/hammond``` path and run the following command 

``` bash
cd /usr/local/bin/hammond && ./hammond
```

Point your browser to http://localhost:3000 (if trying on the same machine) or http://server-ip:3000 from other machines.

If you are using ufw or some other firewall, you might have to make an exception for this port on that.

## Setup as service (Optional)

If you want to run Hammond in the background as a service or auto-start whenever the server starts, follow the next steps.

Create new file named ```hammond.service``` at ```/etc/systemd/system``` and add the following content. You will have to modify the content accordingly if you changed the installation path in the previous steps.


``` unit
[Unit]
Description=Hammond

[Service]
ExecStart=/usr/local/bin/hammond/hammond
WorkingDirectory=/usr/local/bin/hammond/
[Install]
WantedBy=multi-user.target
```

Run the following commands 
``` bash
sudo systemctl daemon-reload
sudo systemctl enable hammond.service
sudo systemctl start hammond.service
```

Run the following command to check the service status.

``` bash
sudo systemctl status hammond.service
```

# Update Hammond

In case you have installed Hammond and want to update the latest version (another area where Docker really shines) you need to repeat the steps from cloning to building and copying.

Stop the running service (if using)
``` bash
sudo systemctl stop hammond.service
```

## Clone from Git

``` bash
git clone --depth 1 https://github.com/akhilrex/hammond
```

## Build and Copy dependencies

``` bash
cd hammond
mkdir -p ./dist
cp .env ./dist
go build -o ./dist/hammond ./main.go
```

## Create final destination and copy executable
``` bash
sudo mkdir -p /usr/local/bin/hammond
mv -v dist/* /usr/local/bin/hammond
```

Go back to the root of the hammond folder.

``` bash
cd ui
npm install
npm run build
mv dist /usr/local/bin/hammond 
```

Restart the service (if using)
``` bash
sudo systemctl start hammond.service
```
