package main

import (
	"context"
	"fmt"
	"time"

	"github.com/xtomars2048/jupiter-go/jupiter"
	"github.com/xtomars2048/jupiter-go/solana"
)

func main() {
	jupClient, err := jupiter.NewClientWithResponses(jupiter.DefaultAPIURL)
	if err != nil {
		panic(err)
	}

	ctx := context.TODO()

	slippageBps := float32(250.0)

	// Get the current quote for a swap.
	// Ensure that the input and output mints are valid.
	// The amount is the smallest unit of the input token.
	quoteResponse, err := jupClient.GetQuoteWithResponse(ctx, &jupiter.GetQuoteParams{
		InputMint:   "So11111111111111111111111111111111111111112",
		OutputMint:  "JUPyiwrYJFskUPiHa7hkeR8VUtAeFoSYbKedZNsDvCN",
		Amount:      100000,
		SlippageBps: &slippageBps,
	})
	if err != nil {
		panic(err)
	}

	if quoteResponse.JSON200 == nil {
		panic("invalid GetQuoteWithResponse response")
	}

	quote := quoteResponse.JSON200

	dynamicComputeUnitLimit := true

	// Define the prioritization fee in lamports.
	prioritizationFeeLamports := &struct {
		JitoTipLamports              *int `json:"jitoTipLamports,omitempty"`
		PriorityLevelWithMaxLamports *struct {
			MaxLamports   *int    `json:"maxLamports,omitempty"`
			PriorityLevel *string `json:"priorityLevel,omitempty"`
		} `json:"priorityLevelWithMaxLamports,omitempty"`
	}{
		JitoTipLamports: new(int),
	}

	*prioritizationFeeLamports.JitoTipLamports = 1000

	// Get instructions for a swap.
	// Ensure your public key is valid.
	swapResponse, err := jupClient.PostSwapWithResponse(ctx, jupiter.PostSwapJSONRequestBody{
		PrioritizationFeeLamports: prioritizationFeeLamports,
		QuoteResponse:             *quote,
		UserPublicKey:             "{YOUR_PUBLIC_KEY}",
		DynamicComputeUnitLimit:   &dynamicComputeUnitLimit,
	})
	if err != nil {
		panic(err)
	}

	if swapResponse.JSON200 == nil {
		panic("invalid PostSwapWithResponse response")
	}

	swap := swapResponse.JSON200
	fmt.Printf("%+v", swap)

	// Create a wallet from private key.
	walletPrivateKey := "{YOUR_PRIVATE_KEY}"
	wallet, err := solana.NewWalletFromPrivateKeyBase58(walletPrivateKey)
	if err != nil {
		panic(err)
	}

	// Create a Solana client. Change the URL to the desired Solana node.
	solanaClient, err := solana.NewClient(wallet, "https://api.mainnet-beta.solana.com")
	if err != nil {
		panic(err)
	}

	// Sign and send the transaction.
	signedTx, err := solanaClient.SendTransactionOnChain(ctx, swap.SwapTransaction)
	if err != nil {
		panic(err)
	}

	// Wait a bit to let the transaction propagate to the network.
	// This is just an example and not a best practice.
	// You could use a ticker or wait until we implement the WebSocket monitoring ;)
	time.Sleep(20 * time.Second)

	// Get the status of the transaction (pull the status from the blockchain at intervals
	// until the transaction is confirmed)
	_, err = solanaClient.CheckSignature(ctx, signedTx)
	if err != nil {
		panic(err)
	}
}
