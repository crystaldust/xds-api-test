package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	apiv2 "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	apiv2core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	apiv2route "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	"github.com/envoyproxy/go-control-plane/envoy/service/discovery/v2"
	"github.com/gogo/protobuf/proto"
	"google.golang.org/grpc"
)

var (
	pilotAddr string
	err       error
	conn      *grpc.ClientConn
)

var (
	cdsVersionInfo   string = time.Now().String()
	cdsResponseNonce string = time.Now().String()

	edsVersionInfo   string = time.Now().String()
	edsResponseNonce string = time.Now().String()

	rdsVersionInfo   string = time.Now().String()
	rdsResponseNonce string = time.Now().String()

	ldsVersionInfo   string = time.Now().String()
	ldsResponseNonce string = time.Now().String()
)

func init() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: xds-api-test {pilot-address}")
		os.Exit(1)
	}

	if os.Args[1] == "" {
		fmt.Println("Invalid pilot-address")
		os.Exit(1)
	}

	pilotAddr = os.Args[1]
	fmt.Println("connecting to ", pilotAddr)

	conn, err = grpc.Dial(pilotAddr, grpc.WithInsecure())
	if err != nil {
		fmt.Println("failed to connect to port")
		fmt.Println(err)
		return
	}
}

func main() {
	clusters, err := cds()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%d clusters found\n", len(clusters))
	jsonPrint(clusters)

	listeners, err := lds()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("listeners:")
	jsonPrint(listeners)

	for _, cluster := range clusters {
		// if strings.Index(cluster.Name, "pilotv2server") != -1 {
		if strings.Index(cluster.Name, "reviews") != -1 {
			// if endpoints, err := eds(cluster.Name); err != nil {
			//     fmt.Println(err)
			// } else {
			//     fmt.Println("endpoints of ", cluster.Name)
			//     jsonPrint(endpoints)
			// }

			if routes, err := rds(cluster.Name); err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("routes of", cluster.Name)
				jsonPrint(routes)
			}

		}
	}
}

func cds() ([]apiv2.Cluster, error) {
	ctx := context.Background()

	adsClient := v2.NewAggregatedDiscoveryServiceClient(conn)
	adsResClient, err := adsClient.StreamAggregatedResources(ctx)
	if err != nil {
		fmt.Println("failed to get stream adsResClient")
		return nil, err
	}

	req := &apiv2.DiscoveryRequest{
		TypeUrl:       "type.googleapis.com/envoy.api.v2.Cluster",
		VersionInfo:   cdsVersionInfo,
		ResponseNonce: cdsResponseNonce,
	}
	fmt.Printf("request clusters with versioninfo[%s] and responsenonce[%s]\n", cdsVersionInfo, cdsResponseNonce)
	req.Node = &apiv2core.Node{
		// Sample taken from istio: router~172.30.77.6~istio-egressgateway-84b4d947cd-rqt45.istio-system~istio-system.svc.cluster.local-2
		// The Node.Id should be in format {nodeType}~{ipAddr}~{serviceId~{domain}, splitted by '~'
		// The format is required by pilot
		Id:      "sidecar~192.168.43.100~xds-api-test~localhost",
		Cluster: "my-powerful-machine-ouya",
	}
	jsonPrint(req)
	if err := adsResClient.Send(req); err != nil {
		return nil, err
	}

	resp, err := adsResClient.Recv()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	cdsResponseNonce = resp.GetNonce()
	resources := resp.GetResources()
	cdsVersionInfo = resp.GetVersionInfo()

	var cluster apiv2.Cluster
	clusters := []apiv2.Cluster{}
	for _, res := range resources {
		if err := proto.Unmarshal(res.GetValue(), &cluster); err != nil {
			fmt.Println("Failed to unmarshal resource: ", err)
		} else {
			clusters = append(clusters, cluster)
		}
	}
	return clusters, nil
}

func eds(clusterName string) ([]apiv2.ClusterLoadAssignment, error) {
	ctx := context.Background()

	adsClient := v2.NewAggregatedDiscoveryServiceClient(conn)
	adsResClient, err := adsClient.StreamAggregatedResources(ctx)
	if err != nil {
		fmt.Println("failed to get stream adsResClient")
		return nil, err
	}

	req := &apiv2.DiscoveryRequest{
		TypeUrl:       "type.googleapis.com/envoy.api.v2.ClusterLoadAssignment",
		VersionInfo:   edsVersionInfo,
		ResponseNonce: edsResponseNonce,
	}

	fmt.Printf("eds with versioninfo[%s] and nonce[%s]\n", edsVersionInfo, edsResponseNonce)
	req.Node = &apiv2core.Node{
		Id:      "sidecar~192.168.43.100~xds-api-test~localhost",
		Cluster: "juzhen-x79",
	}
	req.ResourceNames = []string{clusterName}
	if err := adsResClient.Send(req); err != nil {
		return nil, err
	}

	resp, err := adsResClient.Recv()
	if err != nil {
		return nil, err
	}

	resources := resp.GetResources()
	edsResponseNonce = resp.GetNonce()
	edsVersionInfo = resp.GetVersionInfo()

	var endpoint apiv2.ClusterLoadAssignment
	endpoints := []apiv2.ClusterLoadAssignment{}

	for _, res := range resources {
		if err := proto.Unmarshal(res.GetValue(), &endpoint); err != nil {
			fmt.Println("Failed to unmarshal resource: ", err)
		} else {
			endpoints = append(endpoints, endpoint)
		}
	}
	return endpoints, nil
}

func rds(clusterName string) ([]apiv2route.VirtualHost, error) {
	// Cluster name in format:
	//
	parts := strings.Split(clusterName, "|")
	port := parts[1]
	serviceName := parts[3]

	ctx := context.Background()

	adsClient := v2.NewAggregatedDiscoveryServiceClient(conn)
	adsResClient, err := adsClient.StreamAggregatedResources(ctx)
	if err != nil {
		fmt.Println("failed to get stream adsResClient")
		return nil, err
	}

	req := &apiv2.DiscoveryRequest{
		TypeUrl:       "type.googleapis.com/envoy.api.v2.RouteConfiguration",
		VersionInfo:   rdsVersionInfo,
		ResponseNonce: rdsResponseNonce,
	}

	fmt.Printf("rds with versioninfo[%s] and nonce[%s]\n", edsVersionInfo, edsResponseNonce)
	req.Node = &apiv2core.Node{
		Id:      "sidecar~192.168.43.100~xds-api-test~localhost",
		Cluster: "juzhen-x79",
	}
	req.ResourceNames = []string{port}
	if err := adsResClient.Send(req); err != nil {
		return nil, err
	}

	resp, err := adsResClient.Recv()
	if err != nil {
		return nil, err
	}

	resources := resp.GetResources()
	rdsResponseNonce = resp.GetNonce()
	rdsVersionInfo = resp.GetVersionInfo()

	var route apiv2.RouteConfiguration
	// routes := []apiv2.RouteConfiguration{}
	virtualHosts := []apiv2route.VirtualHost{}

	for _, res := range resources {
		if err := proto.Unmarshal(res.GetValue(), &route); err != nil {
			fmt.Println("Failed to unmarshal resource: ", err)
		} else {
			// Filter the virtual hosts
			// routes = append(routes, route)
			vhosts := route.GetVirtualHosts()
			for _, vhost := range vhosts {
				if vhost.Name == fmt.Sprintf("%s:%s", serviceName, port) {
					virtualHosts = append(virtualHosts, vhost)
					for _, r := range vhost.Routes {
						routerClusterName := r.GetRoute().GetCluster()
						fmt.Println("[DEBUG] routerClusterName: ", routerClusterName)
					}
				}
			}
		}
	}
	return virtualHosts, nil
}

func lds() ([]apiv2.Listener, error) {
	ctx := context.Background()

	adsClient := v2.NewAggregatedDiscoveryServiceClient(conn)
	adsResClient, err := adsClient.StreamAggregatedResources(ctx)
	if err != nil {
		fmt.Println("failed to get stream adsResClient")
		return nil, err
	}

	req := &apiv2.DiscoveryRequest{
		TypeUrl:       "type.googleapis.com/envoy.api.v2.ClusterLoadAssignment",
		VersionInfo:   ldsVersionInfo,
		ResponseNonce: ldsResponseNonce,
	}

	fmt.Printf("lds with versioninfo[%s] and nonce[%s]\n", edsVersionInfo, edsResponseNonce)
	req.Node = &apiv2core.Node{
		Id:      "sidecar~192.168.43.100~xds-api-test~localhost",
		Cluster: "juzhen-x79",
	}
	if err := adsResClient.Send(req); err != nil {
		return nil, err
	}

	resp, err := adsResClient.Recv()
	if err != nil {
		return nil, err
	}

	resources := resp.GetResources()
	ldsResponseNonce = resp.GetNonce()
	ldsVersionInfo = resp.GetVersionInfo()

	var listener apiv2.Listener
	listeners := []apiv2.Listener{}

	for _, res := range resources {
		if err := proto.Unmarshal(res.GetValue(), &listener); err != nil {
			fmt.Println("Failed to unmarshal resource: ", err)
		} else {
			listeners = append(listeners, listener)
		}
	}
	return listeners, nil
}

// TODO Extract the common part of xds calls
func xds(urlType string) error {

	return nil
}

func jsonPrint(content interface{}) {
	bs, _ := json.MarshalIndent(content, "", "  ")
	fmt.Println(string(bs))
}
