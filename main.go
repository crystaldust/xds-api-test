package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/profile"
)

const TARGET_SVC_NAME = "reviews"

func main() {
	defer profile.Start(profile.MemProfile).Stop()
	n := 0
	for n < 1000 {
		fmt.Println(n)
		doQuery()
		n++
	}
}

func doQuery() {
	clusters, err := cds()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%d clusters found\n", len(clusters))
	// jsonPrint(clusters)
	// for _, c := range clusters {
	//     if strings.Index(c.Name, TARGET_SVC_NAME) != -1 {
	// jsonPrint(c)
	//     }
	// }

	// listeners, err := lds()
	// if err != nil {
	//     fmt.Println(err)
	//     return
	// }
	// fmt.Println("listeners:")
	// jsonPrint(listeners)

	for _, cluster := range clusters {
		// if strings.Index(cluster.Name, "pilotv2server") != -1 {
		if strings.Index(cluster.Name, TARGET_SVC_NAME) != -1 {
			if endpoints, err := eds(cluster.Name); err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("endpoints of ", cluster.Name, len(endpoints))
				// jsonPrint(endpoints)
			}

			// if routes, err := rds(cluster.Name); err != nil {
			//     fmt.Println(err)
			// } else {
			//     fmt.Println("routes of", cluster.Name)
			//     jsonPrint(routes)
			// }

		}
	}
}

// TODO Extract the common part of xds calls
func xds(urlType string) error {

	return nil
}

func jsonPrint(content interface{}) {
	bs, _ := json.MarshalIndent(content, "", "  ")
	fmt.Println(string(bs))
}
