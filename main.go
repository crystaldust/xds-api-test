package main

import (
	"context"
	"fmt"
	"os"
	"time"

	apiv2 "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	"github.com/envoyproxy/go-control-plane/envoy/service/discovery/v2"
	"google.golang.org/grpc"
)

func main() {
	// mux := http.NewServeMux()
	// mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	//     w.WriteHeader(http.StatusOK)
	//     w.Write([]byte("Hello"))
	// })
	//
	// http.ListenAndServe(":8881", mux)
	getServices()

	// stop forever
	// for {
	//     time.Sleep(time.Second * 10)
	// }
}

func getServices() {
	// server := &v2.DiscoveryServer{}
	// fmt.Println(server)

	fmt.Println(os.Args)
	pilotAddr := "192.168.43.70:8680"
	if len(os.Args) >= 2 && os.Args[1] != "" {
		pilotAddr = os.Args[1]
	}
	fmt.Println("[JUZHEN DEBUG]: ", "connecting to ", pilotAddr)

	conn, err := grpc.Dial(pilotAddr, grpc.WithInsecure())
	if err != nil {
		fmt.Println("failed to connect to port")
		fmt.Println(err)
		return
	}

	getAds(conn)
	// go getEds(conn)

	// client := v2.NewDiscoveryServer(model.Environment.ServiceDiscovery)
	// fmt.Println(conn)

	// client := v2.NewClusterDiscoveryServiceClient(conn)
	// req := &v2.DiscoveryRequest{}
	// req.ResourceNames = []string{}
	// fmt.Println("req: ", req)
	// resp, err := client.FetchClusters(ctx, req)
	// if err != nil {
	//     fmt.Println("Failed to fetch clusters")
	//     fmt.Println(err)
	// }
	// fmt.Println(resp)

}

func getAds(conn *grpc.ClientConn) {
	ctx := context.Background()

	// ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
	// defer cancel()

	adsClient := v2.NewAggregatedDiscoveryServiceClient(conn)
	adsResClient, err := adsClient.StreamAggregatedResources(ctx)
	if err != nil {
		fmt.Println("failed to get stream adsResClient")
		fmt.Println(err)
	}

	// go func(adsResClient v2.AggregatedDiscoveryService_StreamAggregatedResourcesClient) {
	//     req := &apiv2.DiscoveryRequest{}
	//
	//     fmt.Println("Send: ", adsResClient.Send(req))
	// }(adsResClient)

	fmt.Println(adsResClient != nil)
	for {
		req := &apiv2.DiscoveryRequest{}
		fmt.Println("send: ", adsResClient.Send(req))

		if resp, err := adsResClient.Recv(); err != nil {
			fmt.Println(err)
		} else if bs, e := resp.Marshal(); e != nil {
			fmt.Println(e)
		} else {
			fmt.Println(string(bs))
		}

		time.Sleep(time.Second * 5)
	}
}

func getEds(conn *grpc.ClientConn) {
	ctx := context.Background()

	client := apiv2.NewEndpointDiscoveryServiceClient(conn)
	req := &apiv2.DiscoveryRequest{}
	for {
		if resp, err := client.FetchEndpoints(ctx, req); err != nil {
			fmt.Println("Failed to fetch endpoints: ", err)
		} else if bs, e := resp.Marshal(); e != nil {
			fmt.Println("Failed to marshal EDS resp: ", e)
		} else {
			fmt.Println("EDS resp: ", string(bs))
		}
		time.Sleep(time.Second * 5)
	}
}
