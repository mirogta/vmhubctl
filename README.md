# vmhubctl

CLI for VirginMedia Hub 3 using SMNP

## How to Use

Create a `.vmhubctl.yml` config file based on the `.vmhubctl.yml-template` file.

* Show help: `vmhubctl --help`
* List some settings: `vmhubctl list`
* Disable WiFi: `vmhubctl disable wifi`
* Enable WiFi: `vmhubctl enable wifi`

## How to build

```cmd
mkdir -p bin
GOBIN=$PWD/bin go get golang.org/x/tools/cmd/stringer
./bin/stringer -type=Name ./oids
go build -o vmhubctl
```

NOTE: The `stringer` command needs to be re-ran after changing `./oids/name.go`.

## Compile and Install

The idea is to run this on Raspberri Pi connected via a LAN cable to the Wifi router, so that it could turn on and off Wifi at certain times.

For Raspberry Pi:

```bash
env GOOS=linux GOARCH=arm GOARM=5 go build -o vmhubctl
scp .vmhubctl.yml pi@pihole:~/
scp vmhubctl pi@pihole:~/
```

* Install CRON via `crontab -e`:

```crontab
# switch off wifi at 2300
0 23 * * * /home/pi/vmhubctl disable wifi --config /home/pi/.vmhubctl.yml
# switch on wifi at 0500
0 5 * * * /home/pi/vmhubctl enable wifi --config /home/pi/.vmhubctl.yml
```

* Check if it was running: `grep CRON /var/log/syslog`
* Check logs: `cat ~/vmhubctl.log`

## Useful links

* https://ninet.org/2018/06/fixing-the-virgin-media-superhub-3/
* https://www.netscylla.com/blog/2019/02/04/Arris-CableModem-SNMP.html
* https://github.com/alexmartinio/vmsuperhub-smnp
* https://mibs.observium.org/mib/ARRIS-CM-DEVICE-MIB/
* https://mibs.observium.org/mib/ARRIS-ROUTER-DEVICE-MIB/
