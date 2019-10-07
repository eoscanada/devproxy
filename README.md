## devproxy

Running on the cluster:

```
devproxy --listen-addr=:9000 --services blockmeta:9000,search-archive:9000,search-liverouter:9000,merger:9000,relayer:9000
```

### Development

Development of the `devproxy` project requires you to make a series of port forward.
To make it easier, the project has a script `port-forward.sh` that does all the
port forwarding for you automatically.

The `devproxy` project is configured in development to reach the services defined
in the `port-forward.sh` script.

So to develop the `devproxy` directly, you will need to have two terminal open.
In the first terminal, launch the port forwarding script:

```
./port-forward.sh
```

Then in the second terminal run the `devproxy` project:

```
go build -o devproxy ./ && ./devproxy --listen-addr=:9000 | zap-pretty
```

**Note** Ensure that `port-forward` script `services` variable is correctly
set to the list of services the `devproxy` should do and the `--services` flag
in the `main.go` fits with proxied services.

### Using it

```
kubectl port-forward svc/devproxy 9000
```

Listening and forwarding :9000 to target environment...

```
grpcurl -plaintext localhost:9000 list
grpcui -plaintext -port 8002 localhost:9000
```

Go to http://localhost:8002
