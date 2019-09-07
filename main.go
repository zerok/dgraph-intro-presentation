package main

import (
	"context"
	"fmt"
	"log"

	"github.com/dgraph-io/dgo/v2"
	"github.com/dgraph-io/dgo/v2/protos/api"
	"google.golang.org/grpc"
)

func main() {
	ctx := context.Background()
	d, err := grpc.Dial("localhost:9080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to open a gRPC connection: %s", err.Error())
	}
	defer d.Close()

	dc := api.NewDgraphClient(d)
	client := dgo.NewDgraphClient(dc)

	if err := queryCities(ctx, client); err != nil {
		log.Fatalf("Failed to query cities: %s", err.Error())
	}

	if err := querySpecificCity(ctx, client, "Graz"); err != nil {
		log.Fatalf("Failed to query Graz: %s", err.Error())
	}
}

func queryCities(ctx context.Context, client *dgo.Dgraph) error {
	txn := client.NewReadOnlyTxn()
	defer txn.Discard(ctx)
	res, err := txn.Query(ctx, `{
		cities(func:type(City)) {
			uid
			name@.
		}
	}`)
	if err != nil {
		return err
	}
	fmt.Println(res)
	return nil
}

func querySpecificCity(ctx context.Context, client *dgo.Dgraph, name string) error {
	log.Printf("Querying for city %s", name)
	txn := client.NewReadOnlyTxn()
	defer txn.Discard(ctx)
	res, err := txn.QueryWithVars(ctx, `query City($name: string) {
		cities(func:eq(name@., $name)) @filter(type(City)) {
			uid
			name@.
		}
	}`, map[string]string{
		"$name": name,
	})
	if err != nil {
		return err
	}
	fmt.Println(res)
	return nil
}
