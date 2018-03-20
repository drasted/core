package dealer

import (
	"errors"

	"github.com/sonm-io/core/proto"
)

// Matcher interface describes method to select most profitable order from search results.
type Matcher interface {
	Match(orders []*sonm.Order) (*sonm.Order, error)
}

// bidMatcher matches the cheapest order to deal
type matcher struct {
	less func(lhs, rhs *sonm.Order) bool
}

// NewBidMatcher returns Matcher implementation which
// matches given BID order with most cheapest ASK.
func NewBidMatcher() Matcher {
	return &matcher{
		less: func(lhs, rhs *sonm.Order) bool {
			return lhs.PricePerSecond.Cmp(rhs.PricePerSecond) < 0
		},
	}
}

// NewAskMatcher returns Matcher implementation which
// matches given ASK order with most expensive BID.
func NewAskMatcher() Matcher {
	return &matcher{
		less: func(lhs, rhs *sonm.Order) bool {
			return lhs.PricePerSecond.Cmp(rhs.PricePerSecond) > 0
		},
	}
}

func (m *matcher) Match(orders []*sonm.Order) (*sonm.Order, error) {
	if len(orders) == 0 {
		return nil, errors.New("no orders provided")
	}

	var best = orders[0]
	for _, o := range orders {
		if m.less(o, best) {
			best = o
		}
	}

	return best, nil
}
