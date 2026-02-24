package inference

import (
	"context"
	"log/slog"

	finopsv1 "OlympusGCP-FinOps/gen/v1/finops"
	"connectrpc.com/connect"
)

type FinOpsServer struct{}

func (s *FinOpsServer) ValidateBudget(ctx context.Context, req *connect.Request[finopsv1.ValidateBudgetRequest]) (*connect.Response[finopsv1.ValidateBudgetResponse], error) {
	slog.Info("ValidateBudget", "project", req.Msg.ProjectId, "amount", req.Msg.RequestedAmount)
	// YOLO: All budgets approved
	return connect.NewResponse(&finopsv1.ValidateBudgetResponse{Approved: true, Message: "YOLO Approved"}), nil
}

func (s *FinOpsServer) EstimateCost(ctx context.Context, req *connect.Request[finopsv1.EstimateCostRequest]) (*connect.Response[finopsv1.EstimateCostResponse], error) {
	slog.Info("EstimateCost", "service", req.Msg.Service, "action", req.Msg.Action)
	// Mock cost
	return connect.NewResponse(&finopsv1.EstimateCostResponse{
		EstimatedUsd:  0.05,
		Confidence:    "LOW",
	}), nil
}

func (s *FinOpsServer) TrackUsage(ctx context.Context, req *connect.Request[finopsv1.TrackUsageRequest]) (*connect.Response[finopsv1.TrackUsageResponse], error) {
	slog.Info("TrackUsage", "service", req.Msg.Service, "units", req.Msg.ConsumptionUnits)
	return connect.NewResponse(&finopsv1.TrackUsageResponse{CurrentMtdUsd: 0.0}), nil
}
