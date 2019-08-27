devproxy
--------



Running on the cluster:

$ devproxy --listen-addr=:9000 --services blockmeta:9000,search-archive:9000,search-liverouter:9000,merger:9000,relayer:9000


Using it
--------

[k8s:dfuseio-global:eth-mainnet] ~$ kubectl port-forward svc/devproxy 9000
Listening and forwarding :9000 to target environment...


$ grpcurl -plaintext localhost:9000 list
$ grpcui -plaintext -port 8002 localhost:9000
Go to http://localhost:8002
