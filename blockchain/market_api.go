package blockchain

import (
	"context"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
	marketAPI "github.com/sonm-io/core/blockchain/market/api"
	pb "github.com/sonm-io/core/proto"
	"go.uber.org/zap"
)

const MarketAddress string = "0x6dbc7e4bc031895cab4ef593fb431279ec9b7dcd"

type BlockchainAPI struct {
	client         *ethclient.Client
	gasPrice       int64
	marketContract *marketAPI.Market
	logger         *zap.Logger
}

func NewBlockchainAPI(ethEndpoint *string, gasPrice *int64) (*BlockchainAPI, error) {
	client, err := initEthClient(ethEndpoint)
	if err != nil {
		return nil, err
	}

	if gasPrice == nil {
		*gasPrice = defaultGasPrice
	}

	marketContract, err := marketAPI.NewMarket(common.HexToAddress(MarketAddress), client)
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

func (api *BlockchainAPI) GetDealEvents(ctx context.Context, fromBlockInitial *big.Int) (chan interface{}, error) {
	var (
		topics     [][]common.Hash
		eventTopic = []common.Hash{DealOpenedTopic, DealClosedTopic, RequestAcceptedTopic}
		out        = make(chan interface{}, 128)
	)
	topics = append(topics, eventTopic)

	marketABI, err := abi.JSON(strings.NewReader(string(marketAPI.MarketABI)))
	if err != nil {
		close(out)
		return nil, err
	}

	go func() {
		var (
			fromBlock = fromBlockInitial
			tk        = time.NewTicker(time.Minute)
		)
		for {
			select {
			case <-ctx.Done():
				return
			case <-tk.C:
				logs, err := api.client.FilterLogs(ctx, ethereum.FilterQuery{
					Topics:    topics,
					FromBlock: fromBlock,
				})
				if err != nil {
					api.logger.Error("failed to filter logs", zap.Error(err),
						zap.String("from_block", fromBlock.String()))
				}

				fromBlock = big.NewInt(int64(logs[len(logs)-1].BlockNumber))

				for _, log := range logs {
					if err := api.processLog(log, marketABI, out); err != nil {
						api.logger.Error("failed to process log", zap.Error(err),
							zap.String("topic", log.Topics[0].String()))
					}
				}
			}
		}
	}()

	return out, nil
}

func (api *BlockchainAPI) processLog(log types.Log, marketABI abi.ABI, out chan interface{}) error {
	switch log.Topics[0] {
	case DealOpenedTopic:
		var dealOpenedEvent = &DealOpenedEvent{}
		if err := marketABI.Unpack(&dealOpenedEvent, "DealOpened", log.Data); err != nil {
			return err
		}
		out <- dealOpenedEvent
	case DealClosedTopic:
		var dealClosedEvent = &DealClosedEvent{}
		if err := marketABI.Unpack(&dealClosedEvent, "DealClosed", log.Data); err != nil {
			return err
		}
		out <- dealClosedEvent
	case NewConditionsRequestedTopic:
		var newConditionsRequestedEvent = &NewConditionsRequestedEvent{}
		if err := marketABI.Unpack(&newConditionsRequestedEvent, "NewConditionsRequested", log.Data); err != nil {
			return err
		}
		out <- newConditionsRequestedEvent
	case RequestAcceptedTopic:
		var requestAcceptedEvent = &RequestAcceptedEvent{}
		if err := marketABI.Unpack(&requestAcceptedEvent, "RequestAccepted", log.Data); err != nil {
			return err
		}
		out <- requestAcceptedEvent
	case RequestDeclinedTopic:
		var requestDeclinedEvent = &RequestDeclinedEvent{}
		if err := marketABI.Unpack(&requestDeclinedEvent, "RequestDeclined", log.Data); err != nil {
			return err
		}
		out <- requestDeclinedEvent
	default:
		return errors.Errorf("unknown topic: %s", log.Topics[0].String())
	}

	return nil
}

type DealOpenedEvent struct {
	Supplier common.Address
	Consumer common.Address
	Id       *big.Int
}

type DealClosedEvent struct {
	Supplier common.Address
	Consumer common.Address
	Id       *big.Int
}

type NewConditionsRequestedEvent struct {
	Id          *big.Int
	NewPrice    *big.Int
	NewDuration *big.Int
}

type RequestDeclinedEvent struct {
	Id          *big.Int
	NewPrice    *big.Int
	NewDuration *big.Int
}

type RequestAcceptedEvent struct {
	Id          *big.Int
	NewPrice    *big.Int
	NewDuration *big.Int
}
