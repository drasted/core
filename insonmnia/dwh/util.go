package dwh

import pb "github.com/sonm-io/core/proto"

func getOperator(operator pb.ComparisonOperator) string {
	switch operator {
	case pb.ComparisonOperator_GTE:
		return ">="
	case pb.ComparisonOperator_LTE:
		return "<="
	default:
		return "="
	}
}
