package blockchain

import (
	context "context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sonm-io/core/blockchain/market"
	marketAPI "github.com/sonm-io/core/blockchain/market/api"
	pb "github.com/sonm-io/core/proto"
)

type BlockchainAPI struct {
	client         *ethclient.Client
	gasPrice       int64
	marketContract *marketAPI.Market
}

func NewBlockchainAPI(ethEndpoint *string, gasPrice *int64) (*BlockchainAPI, error) {
	client, err := initEthClient(ethEndpoint)
	if err != nil {
		return nil, err
	}

	if gasPrice == nil {
		*gasPrice = defaultGasPrice
	}

	marketContract, err := marketAPI.NewMarket(common.HexToAddress(market.MarketAddress), client)
	if err != nil {
		return nil, err
	}

	bch := &BlockchainAPI{
		client:         client,
		gasPrice:       *gasPrice,
		marketContract: marketContract,
	}

	return bch, nil
}

func (api *BlockchainAPI) GetDealInfo(ctx context.Context, dealID *big.Int) (*pb.Deal, error) {
	deal, err := api.marketContract.GetDealInfo(getCallOptions(ctx), dealID)
	if err != nil {
		return nil, err
	}

	var benchmarks = make([]uint64, len(deal.Benchmarks))
	for idx, benchmark := range deal.Benchmarks {
		benchmarks[idx] = benchmark.Uint64()
	}

	return &pb.Deal{
		Id:         dealID.String(),
		BuyerID:    deal.Consumer.String(),
		SupplierID: deal.Supplier.String(),
		Price:      pb.NewBigInt(deal.Price),
		Status:     pb.DealStatus(deal.Status),
		StartTime:  &pb.Timestamp{Seconds: deal.StartTime.Int64()},
		WorkTime:   deal.Duration.Uint64(),
		EndTime:    &pb.Timestamp{Seconds: deal.EndTime.Int64()},
		AskID:      deal.AskId.String(),
		BidID:      deal.BidId.String(),
		Benchmarks: benchmarks,
	}, nil
}

func (api *BlockchainAPI) GetDealsAmount(ctx context.Context) (*big.Int, error) {
	return api.marketContract.GetDealsAmount(getCallOptions(ctx))
}

func (api *BlockchainAPI) GetOrderInfo(ctx context.Context, orderID *big.Int) (*pb.Order, error) {
	order, err := api.marketContract.GetOrderInfo(getCallOptions(ctx), orderID)
	if err != nil {
		return nil, err
	}

	var benchmarks = make([]uint64, len(order.Benchmarks))
	for idx, benchmark := range order.Benchmarks {
		benchmarks[idx] = benchmark.Uint64()
	}

	return &pb.Order{
		Id:             orderID.String(),
		BuyerID:        order.Counteragent.String(),
		SupplierID:     order.Author.String(),
		PricePerSecond: pb.NewBigInt(order.Price),
		OrderType:      pb.OrderType(pb.DealStatus(order.OrderType)),
		Slot: &pb.Slot{
			Geo:        &pb.Geo{},
			Duration:   order.Duration.Uint64(),
			Benchmarks: benchmarks,
			Resources: &pb.Resources{
				Properties: make(map[string]float64),
			},
		},
	}, nil
}

func (api *BlockchainAPI) GetOrdersAmount(ctx context.Context) (*big.Int, error) {
	return api.marketContract.GetOrdersAmount(getCallOptions(ctx))
}
