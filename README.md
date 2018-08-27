## Testing istio-pilot's xDS API in a kubernetes cluster

This demo project is an illustration of how to call pilot to get service/config info in a kubernetes cluster.



### Prerequisites
- [golang](https://golang.org/)
- [glide](https://github.com/Masterminds/glide) or [vgo](https://github.com/golang/vgo) to manage the dependencies
- [git](https://git-scm.com/)
- a [kubernetes](https://kubernetes.io/) cluster

### How to run

First manage the dependencies with glide or vgo:
```bash
$ glide install

# or if you're using vgo:
$ vgo install
```

Then make sure `kubectl` and `make` are already in $PATH, then just run `make all`, the scripts will:

- Build the go binary
- Build the docker image `go-chassis/xds-api-test:v1`
- Dispatch the image to the nodes in k8s cluster(make sure the host machine could SSH into the nodes)
- Apply the deployment defined in `xds-api-test.yaml` to the kubernetes cluster

Here we go!



### See the result

If everything goes well with `make`, then we should have a `xds-api-test` pod in `istio-system` namespace, find it and get the logs:

```bash
$ kubectl logs -f -n istio-system $(kubectl get pods -n istio-system | grep xds-api-test | awk '{print $1}')

eds with versioninfo[2018-08-25T09:34:04Z] and nonce[2018-08-25 09:37:07.993629455 +0000 UTC m=+22137.693872102]
endpoints of  outbound|9093||istio-pilot.istio-system.svc.cluster.local
[
  {
    "cluster_name": "outbound|9093||istio-pilot.istio-system.svc.cluster.local",
    "endpoints": [
      {
        "locality": {},
        "lb_endpoints": [
          {
            "endpoint": {
              "address": {
                "Address": {
                  "SocketAddress": {
                    "address": "172.30.1.3",
                    "PortSpecifier": {
                      "PortValue": 9093
                    }
                  }
                }
              }
            }
            // ......
          }
        ]
      }
    ]
  }
]
```

The program first get all the clusters, then get the endpoints of istio-pilot service. Feel free to explore the code and get the info you are interested!
