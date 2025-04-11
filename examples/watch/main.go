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

	// watchTx(
	// 	client,
	// 	"18b1810ba4f0ce794dbbd5d7aee10a832328b1ea4f0f8ebb5a7a8d41a9e2fba1",
	// 	86746488,
	// )

	// watchTx(
	// 	client,
	// 	"357e5d4e0451af3b7ea5db0a2b11efbe490f92db446424b381726e829b3f2cfb",
	// 	82450164,
	// )

	watchTxWithAddress(
		client,
		"357e5d4e0451af3b7ea5db0a2b11efbe490f92db446424b381726e829b3f2cfb",
		82450164,
		"addr_test1xr7xs02kjwr7v3frqrx4exearkd5nmx5ashhzsj5l3nja7yke8x9mpjf7aerjt3n3nfd5tnzkfhlprp09mpf4sdy8dzq6ptcdp",
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
	fmt.Println("Connected to utxorpc host, watching tx ...")

	address := "addr_test1wr9gquc23wc7h8k4chyaad268mjft7t0c08wqertwms70sc0fvx8w"
	addr, err := common.NewAddress(address)
	if err != nil {
		log.Fatalf("failed to create address: %v", err)
	}
	addrCbor, err := addr.MarshalCBOR()
	if err != nil {
		log.Fatalf("failed to marshal address to CBOR: %v", err)
	}

	for stream.Receive() {
		resp := stream.Msg()

		fmt.Printf("Hash %x\n", resp.GetApply().GetCardano().GetHash())
		fmt.Printf("Fee %v\n", resp.GetApply().GetCardano().GetFee())
		// fmt.Printf("Inputs %v\n", resp.GetApply().GetCardano().GetInputs())
		for _, input := range resp.GetApply().GetCardano().GetInputs() {
			fmt.Printf("Input Hash %x\n", input.GetTxHash())
			fmt.Printf("Input Index %v\n", input.GetOutputIndex())
			fmt.Printf("Input Redeemer %v\n", input.GetRedeemer())
			fmt.Printf("Input As Output %v\n", input.GetAsOutput())
		}
		// fmt.Printf("Outputs %v\n", resp.GetApply().GetCardano().GetOutputs())
		for _, output := range resp.GetApply().GetCardano().GetOutputs() {
			outputAddrHex := fmt.Sprintf("%x", output.GetAddress())
			fmt.Printf("Output Address %s\n", outputAddrHex)
			addrHex := fmt.Sprintf("%x", addrCbor)
			fmt.Printf("Compare Address %s\n", addrHex)
			if outputAddrHex == addrHex {
				return
			}
			fmt.Printf("Output Coin %v\n", output.GetCoin())
			fmt.Printf("Output Assets %v\n", output.GetAssets())
			fmt.Printf("Output Script %v\n", output.GetScript())
			fmt.Printf("Output Datum %x\n", output.GetDatum())
		}
		fmt.Printf("Mint %v\n", resp.GetApply().GetCardano().GetMint())
		fmt.Printf("Successful %v\n", resp.GetApply().GetCardano().GetSuccessful())
		fmt.Printf("Validity %v\n", resp.GetApply().GetCardano().GetValidity())
		fmt.Printf("Withdrawals %v\n", resp.GetApply().GetCardano().GetWithdrawals())
		fmt.Printf("Witnesses %v\n", resp.GetApply().GetCardano().GetWitnesses())
		fmt.Printf("Collateral %v\n", resp.GetApply().GetCardano().GetCollateral())
		fmt.Printf("Certificates %v\n", resp.GetApply().GetCardano().GetCertificates())
		fmt.Printf("Auxiliary %v\n", resp.GetApply().GetCardano().GetAuxiliary())

	}

	if err := stream.Err(); err != nil {
		fmt.Println("Stream ended with error:", err)
	} else {
		fmt.Println("Stream ended normally.")
	}
}

func watchTxWithAddress(
	client *utxorpc.UtxorpcClient,
	blockHash string,
	blockIndex int64,
	address string,
) {
	fmt.Println("connecting to utxorpc host:", client.URL())
	ctx := context.Background()
	addr, err := common.NewAddress(address)
	if err != nil {
		log.Fatalf("failed to create address: %v", err)
	}
	addrCbor, err := addr.MarshalCBOR()
	if err != nil {
		log.Fatalf("failed to marshal address to CBOR: %v", err)
	}
	req := &watch.WatchTxRequest{
		Predicate: &watch.TxPredicate{
			// AnyOf: []*watch.TxPredicate{
			// 	{
			Match: &watch.AnyChainTxPattern{
				Chain: &watch.AnyChainTxPattern_Cardano{
					Cardano: &cardano.TxPattern{
						HasAddress: &cardano.AddressPattern{
							ExactAddress: addrCbor,
						},
					},
				},
			},
			// 	},
			// },
		},
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

		fmt.Printf("resp %x\n", resp.GetApply().GetCardano().GetHash())

	}

	if err := stream.Err(); err != nil {
		fmt.Println("Stream ended with error:", err)
	} else {
		fmt.Println("Stream ended normally.")
	}
}
