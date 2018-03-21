package dwh

import (
	"context"
	"fmt"
	"os"
	"testing"

	log "github.com/noxiouz/zapctx/ctxlog"
	pb "github.com/sonm-io/core/proto"
	"github.com/stretchr/testify/assert"
)

const (
	testDBPath = "/tmp/sonm/dwh.db"
)

var (
	w *DWH
)

func TestMain(m *testing.M) {
	var err error
	w, err = getTestDWH()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	retCode := m.Run()
	w.db.Close()
	os.Remove(testDBPath)
	os.Exit(retCode)
}

func TestDWH_GetOrdersList(t *testing.T) {
	reply, err := w.GetOrdersList(context.Background(), &pb.OrdersListRequest{
		Type:  pb.OrderType_ANY,
		Limit: 1,
	})

	if err != nil {
		t.Error(err)
		return
	}

	if len(reply.Orders) != 1 {
		t.Errorf("Expected 1 oreder in reply, got %d", len(reply.Orders))
		return
	}

	assert.Equal(t, uint64(100), reply.Orders[0].Slot.Duration)

	reply, err = w.GetOrdersList(context.Background(), &pb.OrdersListRequest{
		Type:   pb.OrderType_ANY,
		Limit:  1,
		Offset: 1,
	})

	if err != nil {
		t.Error(err)
		return
	}

	if len(reply.Orders) != 1 {
		t.Errorf("Expected 1 order in reply, got %d", len(reply.Orders))
		return
	}

	assert.Equal(t, uint64(200), reply.Orders[0].Slot.Duration)
}

func TestDWH_GetOrderDetails(t *testing.T) {
	reply, err := w.GetOrderDetails(context.Background(), &pb.ID{Id: "2222"})

	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, uint64(200), reply.Slot.Duration)
}

func TestDWH_GetDealDetails(t *testing.T) {
	reply, err := w.GetDealDetails(context.Background(), &pb.ID{Id: "4444"})

	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, int64(41), reply.StartTime.Seconds)
}

func TestDWH_GetDealChangeRequests(t *testing.T) {
	reply, err := w.GetDealChangeRequests(context.Background(), &pb.ID{Id: "4444"})

	if err != nil {
		t.Error(err)
		return
	}

	if len(reply.ChangeRequests) != 2 {
		t.Errorf("Expected 2 change requests in reply, got %d", len(reply.ChangeRequests))
		return
	}

	assert.Equal(t, uint64(1234), reply.ChangeRequests[0].DurationSeconds)
}

func getTestDWH() (*DWH, error) {
	var (
		ctx = context.Background()
		cfg = &Config{
			Storage: &storageConfig{
				Backend:  "sqlite3",
				Endpoint: testDBPath,
			},
		}
		w = &DWH{
			ctx:    ctx,
			cfg:    cfg,
			logger: log.GetLogger(ctx),
		}
	)

	w.mu.Lock()
	defer w.mu.Unlock()

	if err := w.setupDB(); err != nil {
		return nil, err
	}

	_, err := w.db.Exec(`INSERT INTO orders VALUES (1111, 1, "author_1", "counter_agent_1", 100, 200);`)
	if err != nil {
		return nil, err
	}

	_, err = w.db.Exec(`INSERT INTO orders VALUES (2222, 1, "author_2", "counter_agent_2", 200, 400);`)
	if err != nil {
		return nil, err
	}

	_, err = w.db.Exec(`INSERT INTO deals VALUES (3333, 1, "supplier_1", "consumer_1", 100, 200, 40);`)
	if err != nil {
		return nil, err
	}

	_, err = w.db.Exec(`INSERT INTO deals VALUES (4444, 1, "supplier_2", "consumer_2", 200, 400, 41);`)
	if err != nil {
		return nil, err
	}

	_, err = w.db.Exec(`INSERT INTO change_requests VALUES (1234, 5678, 4444);`)
	if err != nil {
		return nil, err
	}

	_, err = w.db.Exec(`INSERT INTO change_requests VALUES (2345, 6789, 4444);`)
	if err != nil {
		return nil, err
	}

	return w, nil
}
