package blockchain

import (
	context "context"

	bind "github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	DealOpenedTopic             = common.HexToHash("0x873cb35202fef184c9f8ee23c04e36dc38f3e26fb285224ca574a837be976848")
	DealClosedTopic             = common.HexToHash("0x72615f99a62a6cc2f8452d5c0c9cbc5683995297e1d988f09bb1471d4eefb890")
	OrderPlacedTopic            = common.HexToHash("0x3")
	NewConditionsRequestedTopic = common.HexToHash("0x4")
	RequestAcceptedTopic        = common.HexToHash("0x5")
	RequestDeclinedTopic        = common.HexToHash("0x6")
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
