# ImmuDB web client
This package is a rest proxy that convert json rest calls in grpc for immudb service.
Is implemented with [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway).

The docker-compose file launch 4 containers:
* the **immudb service** which store data in immudb-data folder on 8083
* the **grpc-gateway rest proxy** on 8081
* a **swagger client** on port 8084
* an **immu web client** on 8085

To verify immudb functionalities set all values you like with demo client.

## Tamperproof
To tamper the database use [nimmu](https://github.com/codenotary/immudb/tree/master/tools/nimmu)
```bash
sudo ./nimmu rawset {key} {val} -d /path/to/immudb-data/demo/data
```
Verification is done inside wrapped safe set get methods.

## Improvements
Add tests
