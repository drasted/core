package blockchain

import (
	context "context"

	bind "github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
)

func initEthClient(ethEndpoint *string) (*ethclient.Client, error) {
	var endpoint string
	if ethEndpoint == nil {
		endpoint = defaultEthEndpoint
	} else {
		endpoint = *ethEndpoint
	}

	ethClient, err := ethclient.Dial(endpoint)
	if err != nil {
		return nil, err
	}

	return ethClient, nil
}

func getCallOptions(ctx context.Context) *bind.CallOpts {
	return &bind.CallOpts{
		Pending: true,
		Context: ctx,
	}
}
