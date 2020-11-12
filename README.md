# kafka-best-practices

Read more: https://medium.com/swlh/how-to-consume-kafka-efficiently-in-golang-264f7fe2155b

Run
```shell script
# produce
go run main.go -m produce -h bootstrap.192.168.55.44.nip.io:443 -u useradmin -p qjFTBoCt6n8O -cert /home/nv4re/cluster-1-ca/ca.crt -ca /home/nv4re/cluster-1-ca/ca.p12 -key /home/nv4re/cluster-1-ca/ca.key -t topic-1

# consume sync
go run main.go -m batch -h bootstrap.192.168.55.44.nip.io:443 -u useradmin -p qjFTBoCt6n8O -cert /home/nv4re/cluster-1-ca/ca.crt -ca /home/nv4re/cluster-1-ca/ca.p12 -key /home/nv4re/cluster-1-ca/ca.key -t topic-1

# consume batch
go run main.go -m batch -h bootstrap.192.168.55.44.nip.io:443 -u useradmin -p qjFTBoCt6n8O -cert /home/nv4re/cluster-1-ca/ca.crt -ca /home/nv4re/cluster-1-ca/ca.p12 -key /home/nv4re/cluster-1-ca/ca.key -t topic-1
```

This repo discusses some techniques of consuming kafka (Sync, Batch, MultiAsync and MultiBatch) and try to demonstrate some best practices which I think would be generally useful to consume data efficiently.

* Sync

consume messages one by one

* Batch

consume messages batch by batch

* MultiAsync

the "Fan In / Fan Out" pattern

* MultiBatch

the "Fan In / Fan Out" pattern batch by batch 

