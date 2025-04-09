package main

import (
	"context"
	"fmt"

	"log"
	"os"

	"github.com/blinklabs-io/gouroboros/ledger/common"
	"github.com/utxorpc/go-codegen/utxorpc/v1alpha/cardano"
	"github.com/utxorpc/go-codegen/utxorpc/v1alpha/watch"
	utxorpc "github.com/utxorpc/go-sdk"
)

func main() {
	fmt.Println("Hello, Go!")

	baseUrl := os.Getenv("UTXORPC_URL")
	if baseUrl == "" {
		baseUrl = "https://preview.utxorpc-v0.demeter.run"
	}
	client := utxorpc.NewClient(utxorpc.WithBaseUrl(baseUrl))
	dmtrApiKey := os.Getenv("DMTR_API_KEY")
	// set API key for demeter
	if dmtrApiKey != "" {
		client.SetHeader("dmtr-api-key", dmtrApiKey)
	}

	// watchTxWithContext(
	// 	client,
	// 	"b826466e04a612077c887cc842a6ca2143b08534e5ec18b77c07084f0259a8a6",
	// 	82450281,
	// )

	watchTxWithContext(
		client,
		"357e5d4e0451af3b7ea5db0a2b11efbe490f92db446424b381726e829b3f2cfb",
		82450164,
	)
}

func watchTx(
	client *utxorpc.UtxorpcClient,
	blockHash string,
	blockIndex int64,
) {
	fmt.Println("connecting to utxorpc host:", client.URL())
	stream, err := client.WatchTx(blockHash, blockIndex)
	if err != nil {
		utxorpc.HandleError(err)
		return
	}
	fmt.Println("Connected to utxorpc host, watching tx with context ...")

	for stream.Receive() {
		resp := stream.Msg()

		fmt.Println("resp", resp)

	}

	if err := stream.Err(); err != nil {
		fmt.Println("Stream ended with error:", err)
	} else {
		fmt.Println("Stream ended normally.")
	}
}

func watchTxWithContext(
	client *utxorpc.UtxorpcClient,
	blockHash string,
	blockIndex int64,
) {
	fmt.Println("connecting to utxorpc host:", client.URL())
	ctx := context.Background()

	addr, err := common.NewAddress("addr_test1vz09v9yfxguvlp0zsnrpa3tdtm7el8xufp3m5lsm7qxzclgmzkket")
	if err != nil {
		log.Fatalf("failed to create address: %v", err)
	}
	addrCbor, err := addr.MarshalCBOR()
	if err != nil {
		log.Fatalf("failed to marshal address to CBOR: %v", err)
	}

	txPattern := &watch.AnyChainTxPattern{
		Chain: &watch.AnyChainTxPattern_Cardano{
			Cardano: &cardano.TxPattern{
				HasAddress: &cardano.AddressPattern{
					ExactAddress: addrCbor,
				},
			},
		},
	}

	predicate1 := &watch.TxPredicate{
		Match: txPattern,
	}

	predicate := &watch.TxPredicate{
		AnyOf: []*watch.TxPredicate{predicate1},
	}

	req := &watch.WatchTxRequest{
		Predicate: predicate,
		Intersect: utxorpc.WatchIntersect(blockHash, blockIndex),
	}

	stream, err := client.WatchTxWithContext(ctx, req)
	if err != nil {
		utxorpc.HandleError(err)
		return
	}
	fmt.Println("Connected to utxorpc host, watching tx with context ...")

	for stream.Receive() {
		resp := stream.Msg()

		fmt.Printf("resp: %x\n", resp.GetApply().GetCardano().Inputs[0].TxHash)

	}

	if err := stream.Err(); err != nil {
		fmt.Println("Stream ended with error:", err)
	} else {
		fmt.Println("Stream ended normally.")
	}
}
