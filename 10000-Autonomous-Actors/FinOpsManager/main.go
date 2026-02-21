package main

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	finopsv1 "OlympusGCP-FinOps/40000-Communication-Contracts/430-Protocol-Definitions/000-gen/finops/v1"
	"OlympusGCP-FinOps/40000-Communication-Contracts/430-Protocol-Definitions/000-gen/finops/v1/finopsv1connect"

	"connectrpc.com/connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
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
		EstimatedCost: 0.05,
		EstimatedUsd:  0.05,
		Confidence:    "LOW",
	}), nil
}

func (s *FinOpsServer) TrackUsage(ctx context.Context, req *connect.Request[finopsv1.TrackUsageRequest]) (*connect.Response[finopsv1.TrackUsageResponse], error) {
	slog.Info("TrackUsage", "service", req.Msg.Service, "units", req.Msg.ConsumptionUnits)
	return connect.NewResponse(&finopsv1.TrackUsageResponse{CurrentMtdUsd: 0.0}), nil
}

func main() {
	server := &FinOpsServer{}
	mux := http.NewServeMux()
	path, handler := finopsv1connect.NewFinOpsServiceHandler(server)
	mux.Handle(path, handler)

	port := "8098" // From genesis.json
	slog.Info("FinOpsManager starting", "port", port)

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           h2c.NewHandler(mux, &http2.Server{}),
		ReadHeaderTimeout: 3 * time.Second,
	}
	err := srv.ListenAndServe()
	if err != nil {
		slog.Error("Server failed", "error", err)
	}
}
