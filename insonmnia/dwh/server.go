package dwh

import (
	"crypto/ecdsa"
	"crypto/tls"
	"database/sql"
	"fmt"
	"math/big"
	"net"
	"strings"
	"sync"

	"time"

	"github.com/grpc-ecosystem/go-grpc-prometheus"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/noxiouz/zapctx/ctxlog"
	"github.com/pkg/errors"
	pb "github.com/sonm-io/core/proto"
	"github.com/sonm-io/core/util"
	"github.com/sonm-io/core/util/xgrpc"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type DWH struct {
	mu          sync.RWMutex
	ctx         context.Context
	cfg         *Config
	cancel      context.CancelFunc
	grpc        *grpc.Server
	logger      *zap.Logger
	db          *sql.DB
	creds       credentials.TransportCredentials
	certRotator util.HitlessCertRotator
}

func NewDWH(ctx context.Context, cfg *Config, key *ecdsa.PrivateKey) (w *DWH, err error) {
	ctx, cancel := context.WithCancel(ctx)
	defer func() {
		if err != nil {
			cancel()
		}
	}()

	w = &DWH{
		ctx:    ctx,
		cfg:    cfg,
		logger: log.GetLogger(ctx),
	}

	var TLSConfig *tls.Config
	w.certRotator, TLSConfig, err = util.NewHitlessCertRotator(ctx, key)
	if err != nil {
		return nil, err
	}

	w.creds = util.NewTLS(TLSConfig)
	server := xgrpc.NewServer(w.logger,
		xgrpc.Credentials(w.creds),
		xgrpc.DefaultTraceInterceptor(),
	)
	w.grpc = server
	pb.RegisterDWHServer(w.grpc, w)
	grpc_prometheus.Register(w.grpc)

	return
}

func (w *DWH) Serve() error {
	if err := w.setupDB(); err != nil {
		return err
	}

	lis, err := net.Listen("tcp", w.cfg.ListenAddr)
	if err != nil {
		return err
	}

	return w.grpc.Serve(lis)
}

func (w *DWH) GetOrdersList(ctx context.Context, request *pb.OrdersListRequest) (*pb.OrdersListReply, error) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	var (
		query      = "SELECT * FROM orders"
		conditions = []string{" WHERE"}
		values     []interface{}
	)
	// Prepare WHERE clause.
	if request.Type != pb.OrderType_ANY {
		conditions = append(conditions, "type=?")
		values = append(values, request.Type)
	}
	if len(request.AuthorID) > 0 {
		conditions = append(conditions, "author=?")
		values = append(values, request.AuthorID)
	}
	if len(request.CounterAgentID) > 0 {
		conditions = append(conditions, "counter_agent=?")
		values = append(values, request.CounterAgentID)
	}
	if request.DurationSeconds > 0 {
		conditions = append(conditions, fmt.Sprintf("duration%s?", getOperator(request.DurationOperator)))
		values = append(values, request.DurationSeconds)
	}
	if request.Price > 0 {
		conditions = append(conditions, fmt.Sprintf("price%s?", getOperator(request.PriceOperator)))
		values = append(values, request.Price)
	}
	query += " ORDER BY price ASC"
	if len(conditions) > 1 {
		query += strings.Join(conditions, " AND ")
	}
	if request.Limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", request.Limit)
	}
	if request.Offset > 0 {
		query += fmt.Sprintf(" OFFSET %d", request.Offset)
	}
	query += ";"

	rows, err := w.db.Query(query, values...)
	if err != nil {
		return nil, errors.Wrapf(err, "query `%s` failed", query)
	}
	defer rows.Close()

	var orders []*pb.Order
	for rows.Next() {
		var (
			id           string
			orderType    uint64
			author       string
			counterAgent string
			duration     uint64
			price        string
		)
		if err := rows.Scan(&id, &orderType, &author, &counterAgent, &duration, &price); err != nil {
			w.logger.Error("failed to scan order row", zap.Error(err))
			return nil, err
		}

		bigPrice := new(big.Int)
		bigPrice.SetString(price, 10)
		orders = append(orders, &pb.Order{
			Id:             id,
			ByuerID:        author,
			SupplierID:     counterAgent,
			PricePerSecond: pb.NewBigInt(bigPrice),
			Slot: &pb.Slot{
				Duration: duration,
			},
		})
	}

	return &pb.OrdersListReply{Orders: orders}, nil
}

func (w *DWH) GetOrderDetails(ctx context.Context, request *pb.ID) (*pb.Order, error) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	rows, err := w.db.Query("SELECT * FROM orders WHERE id=?", request.Id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if ok := rows.Next(); !ok {
		return nil, errors.Errorf("order `%s` not found", request.Id)
	}

	var (
		id           string
		orderType    uint64
		author       string
		counterAgent string
		duration     uint64
		price        string
	)
	if err := rows.Scan(&id, &orderType, &author, &counterAgent, &duration, &price); err != nil {
		return nil, err
	}

	bigPrice := new(big.Int)
	bigPrice.SetString(price, 10)
	return &pb.Order{
		Id:             id,
		ByuerID:        author,
		SupplierID:     counterAgent,
		PricePerSecond: pb.NewBigInt(bigPrice),
		Slot: &pb.Slot{
			Duration: duration,
		},
	}, nil
}

func (w *DWH) GetDealsList(ctx context.Context, request *pb.DealsListRequest) (*pb.DealsListReply, error) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	var (
		query      = "SELECT * FROM deals"
		conditions = []string{" WHERE"}
		values     []interface{}
	)
	// Prepare WHERE clause.
	if request.Status != pb.DealStatus_ANY_STATUS {
		conditions = append(conditions, "status=?")
		values = append(values, request.Status)
	}
	if len(request.SupplierID) > 0 {
		conditions = append(conditions, "supplier=?")
		values = append(values, request.SupplierID)
	}
	if len(request.ConsumerID) > 0 {
		conditions = append(conditions, "consumer=?")
		values = append(values, request.ConsumerID)
	}
	if request.DurationSeconds > 0 {
		conditions = append(conditions, fmt.Sprintf("duration%s?", getOperator(request.DurationOperator)))
		values = append(values, request.DurationSeconds)
	}
	if request.Price > 0 {
		conditions = append(conditions, fmt.Sprintf("price%s?", getOperator(request.PriceOperator)))
		values = append(values, request.Price)
	}
	if len(conditions) > 1 {
		query += strings.Join(conditions, " AND ")
	}
	query += " ORDER BY price ASC"
	if request.Limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", request.Limit)
	}
	if request.Offset > 0 {
		query += fmt.Sprintf(" OFFSET %d", request.Offset)
	}
	query += ";"

	rows, err := w.db.Query(query, values...)
	if err != nil {
		return nil, errors.Wrapf(err, "query `%s` failed", query)
	}
	defer rows.Close()

	var deals []*pb.Deal
	for rows.Next() {
		var (
			id        string
			status    uint64
			supplier  string
			consumer  string
			duration  uint64
			price     string
			startTime int64
		)
		if err := rows.Scan(&id, &status, &supplier, &consumer, &duration, &price, &startTime); err != nil {
			return nil, err
		}

		bigPrice := new(big.Int)
		bigPrice.SetString(price, 10)
		deals = append(deals, &pb.Deal{
			Id:         id,
			BuyerID:    consumer,
			SupplierID: supplier,
			Status:     pb.DealStatus(status),
			Price:      pb.NewBigInt(bigPrice),
			StartTime:  &pb.Timestamp{Seconds: startTime},
			WorkTime:   duration,
		})
	}

	return &pb.DealsListReply{Deals: deals}, nil
}

func (w *DWH) GetDealDetails(ctx context.Context, request *pb.ID) (*pb.Deal, error) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	rows, err := w.db.Query("SELECT * FROM deals WHERE id=?", request.Id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if ok := rows.Next(); !ok {
		return nil, errors.Errorf("deal `%s` not found", request.Id)
	}

	var (
		id        string
		status    uint64
		supplier  string
		consumer  string
		duration  uint64
		price     string
		startTime int64
	)
	if err := rows.Scan(&id, &status, &supplier, &consumer, &duration, &price, &startTime); err != nil {
		return nil, err
	}

	bigPrice := new(big.Int)
	bigPrice.SetString(price, 10)
	return &pb.Deal{
		Id:         id,
		BuyerID:    consumer,
		SupplierID: supplier,
		Status:     pb.DealStatus(status),
		Price:      pb.NewBigInt(bigPrice),
		StartTime:  &pb.Timestamp{Seconds: startTime},
		WorkTime:   duration,
	}, nil
}

func (w *DWH) GetDealChangeRequests(ctx context.Context, request *pb.ID) (*pb.DealChangeRequestsReply, error) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	rows, err := w.db.Query("SELECT * FROM change_requests WHERE deal=?", request.Id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var changeRequests []*pb.DealChangeRequest
	for rows.Next() {
		var (
			duration uint64
			price    string
			dealID   string
		)
		if err := rows.Scan(&duration, &price, &dealID); err != nil {
			return nil, err
		}

		bigPrice := new(big.Int)
		bigPrice.SetString(price, 10)
		changeRequests = append(changeRequests, &pb.DealChangeRequest{
			DurationSeconds: duration,
			Price:           pb.NewBigInt(bigPrice),
		})
	}

	return &pb.DealChangeRequestsReply{ChangeRequests: changeRequests}, nil
}

func (w *DWH) setupDB() (err error) {
	db, err := w.setupSQLite()
	if err != nil {
		return err
	}

	w.db = db

	return nil
}

func (w *DWH) setupSQLite() (*sql.DB, error) {
	db, err := sql.Open(w.cfg.Storage.Backend, w.cfg.Storage.Endpoint)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			db.Close()
		}
	}()

	_, err = db.Exec(createTableDealsSQLite)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create deals table (%s)", w.cfg.Storage.Backend)
	}

	_, err = db.Exec(createTableOrdersSQLite)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create orders table (%s)", w.cfg.Storage.Backend)
	}

	_, err = db.Exec(createTableChangesSQLite)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create orders table (%s)", w.cfg.Storage.Backend)
	}

	return db, nil
}

func (w *DWH) monitor() error {
	w.logger.Info("starting initial sync")
	if err := w.sync(); err != nil {
		return err
	}

	tk := time.NewTicker(time.Second)
	for {
		select {
		case <-tk.C:
			if err := w.sync(); err != nil {
				w.logger.Error("failed to sync", zap.Error(err))
			}
		case <-w.ctx.Done():
			w.logger.Info("stop blockchain monitoring")
			return nil
		}
	}
}

func (w *DWH) sync() error {
	if err := w.syncOrdersTS(); err != nil {
		return err
	}

	return w.syncDealsTS()
}

func (w *DWH) syncOrdersTS() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	id, err := w.getLastOrderID()
	if err != nil {
		return err
	}

	return nil
}

func (w *DWH) getLastOrderID() (string, error) {
	rows, err := w.db.Query("SELECT * FROM orders ORDER BY id DESC LIMIT 1")
	if err != nil {
		return "", err
	}
	defer rows.Close()

	if ok := rows.Next(); !ok {
		return "", errors.New("no entries found")
	}

	var (
		id           string
		orderType    uint64
		author       string
		counterAgent string
		duration     uint64
		price        string
	)
	if err := rows.Scan(&id, &orderType, &author, &counterAgent, &duration, &price); err != nil {
		return "", err
	}

	return id, nil
}

func (w *DWH) syncDealsTS() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	id, err := w.getLastOrderID()
	if err != nil {
		return err
	}

	return nil
}

func (w *DWH) getLastDealID() (string, error) {
	rows, err := w.db.Query("SELECT * FROM deals ORDER BY id DESC LIMIT 1")
	if err != nil {
		return "", err
	}
	defer rows.Close()

	if ok := rows.Next(); !ok {
		return "", errors.New("no entries found")
	}

	var (
		id        string
		status    uint64
		supplier  string
		consumer  string
		duration  uint64
		price     string
		startTime int64
	)
	if err := rows.Scan(&id, &status, &supplier, &consumer, &duration, &price, &startTime); err != nil {
		return nil, err
	}

	return id, nil
}
