# <img src="https://avatars.githubusercontent.com/u/93947921?s=200&" width="32"/> Netmore Connect Go sample application.

# Setup
Download you certificate and stor it in a directory with the same name as you customerId.
```
mkdir certs/<customerId>
cd certs/<customerId>/
unzip download.zip
> client.crt
> client.key
> ca.crt
```
# Build
```
go build
```
# Run
```
./netmore-mqtt-go-sample <customerId>
```
