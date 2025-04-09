package main

import (
	"encoding/hex"
	"fmt"
	"os"

	sync "github.com/utxorpc/go-codegen/utxorpc/v1alpha/sync"
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

	// Run them all
	// fetchBlock(
	// 	client,
	// 	"f3417447ae82a8e0b22fbbc06f720bee844a5e4f8839b0daf303bdc9fcdeab8c",
	// 	2127764,
	// )
	// followTip(
	// 	client,
	// 	"6d6ed8490f568af3170eded8f0d55c363020a1e374668879e94ae8f782b485f7",
	// 	375460,
	// )

	// mainnet
	followTip(
		client,
		"c287feed566fb59b7959cd4f6dcefce46d184edb1495ced707465ee29890116d",
		15998,
	)
}

func fetchBlock(
	client *utxorpc.UtxorpcClient,
	blockHash string,
	blockIndex int64,
) {
	fmt.Println("connecting to utxorpc host:", client.URL())
	resp, err := client.FetchBlock(blockHash, blockIndex)
	if err != nil {
		utxorpc.HandleError(err)
	}
	fmt.Printf("Response: %+v\n", resp)
	for i, blockRef := range resp.Msg.GetBlock() {
		fmt.Printf("Block[%d]:\n", i)
		fmt.Printf("Index: %d\n", blockRef.GetCardano().GetHeader().GetSlot())
		fmt.Printf("Hash: %x\n", blockRef.GetCardano().GetHeader().GetHash())
	}
}

func followTip(
	client *utxorpc.UtxorpcClient,
	blockHash string,
	blockIndex int64,
) {
	fmt.Println("connecting to utxorpc host:", client.URL())
	stream, err := client.FollowTip(blockHash, blockIndex)
	if err != nil {
		utxorpc.HandleError(err)
		return
	}
	fmt.Println("Connected to utxorpc host, following tip...")

	for stream.Receive() {
		resp := stream.Msg()
		// action := resp.GetAction()
		fmt.Println(resp.GetReset_().Index)
		// switch a := action.(type) {
		// case *sync.FollowTipResponse_Apply:
		// 	fmt.Println("Action: Apply")
		// 	printAnyChainBlock(a.Apply)
		// case *sync.FollowTipResponse_Undo:
		// 	fmt.Println("Action: Undo")
		// 	printAnyChainBlock(a.Undo)
		// case *sync.FollowTipResponse_Reset_:
		// 	fmt.Println("Action: Reset")
		// 	printBlockRef(a.Reset_)
		// default:
		// 	fmt.Println("Unknown action type")
		// }
	}

	if err := stream.Err(); err != nil {
		fmt.Println("Stream ended with error:", err)
	} else {
		fmt.Println("Stream ended normally.")
	}
}

func printAnyChainBlock(block *sync.AnyChainBlock) {
	if block == nil {
		return
	}
	if cardanoBlock := block.GetCardano(); cardanoBlock != nil {
		hash := hex.EncodeToString(cardanoBlock.GetHeader().GetHash())
		slot := cardanoBlock.GetHeader().GetSlot()
		fmt.Printf("Block Slot: %d, Block Hash: %s\n", slot, hash)
	}
}

func printBlockRef(blockRef *sync.BlockRef) {
	if blockRef == nil {
		return
	}
	hash := hex.EncodeToString(blockRef.GetHash())
	slot := blockRef.GetIndex()
	fmt.Printf("Block Slot: %d, Block Hash: %s\n", slot, hash)
}
