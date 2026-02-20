package main

import (
	"context"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"connectrpc.com/connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	econotel "Olympus2/90000-Enablement-Labs/P0000-pkg/000-econotel"
	whisper "Olympus2/90000-Enablement-Labs/P0000-pkg/000-whisper"
	finv1 "OlympusGCP-FinOps/40000-Communication-Contracts/430-Protocol-Definitions/000-gen/finops/v1"
	"OlympusGCP-FinOps/40000-Communication-Contracts/430-Protocol-Definitions/000-gen/finops/v1/finopsv1connect"
)

type FinOpsServer struct {
	logger *whisper.WhisperLog
	mu     sync.Mutex
	usage  map[string]float64
	budget float64
}

func (s *FinOpsServer) EstimateCost(ctx context.Context, req *connect.Request[finv1.EstimateCostRequest]) (*connect.Response[finv1.EstimateCostResponse], error) {
	start := time.Now()

	// Workstation Emulation: Emulate local price lookup
	price := 0.0
	switch req.Msg.Service {
	case "ComputeEngine":
		price = 0.05 // per hour
	case "CloudStorage":
		price = 0.02 // per GB
	case "BigQuery":
		price = 5.00 // per TB
	default:
		price = 0.01
	}

	s.logger.Log("ESTIMATE_COST", "SUCCESS", req.Msg.Service, "FORGED_PRICING", time.Since(start))
	return connect.NewResponse(&finv1.EstimateCostResponse{
		EstimatedUsd: price,
		Confidence:   "HIGH_LOCAL",
	}), nil
}

func (s *FinOpsServer) ValidateBudget(ctx context.Context, req *connect.Request[finv1.ValidateBudgetRequest]) (*connect.Response[finv1.ValidateBudgetResponse], error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	used := 0.0
	for _, v := range s.usage {
		used += v
	}

	approved := (used + req.Msg.RequestedAmount) <= s.budget
	msg := "Within Workstation Local Quota"
	if !approved {
		msg = "Workstation Budget Exceeded"
	}

	return connect.NewResponse(&finv1.ValidateBudgetResponse{
		Approved: approved,
		Message:  msg,
	}), nil
}

func (s *FinOpsServer) TrackUsage(ctx context.Context, req *connect.Request[finv1.TrackUsageRequest]) (*connect.Response[finv1.TrackUsageResponse], error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Track usage units for the requested service
	s.usage[req.Msg.Service] += req.Msg.ConsumptionUnits

	used := 0.0
	for _, v := range s.usage {
		used += v
	}

	return connect.NewResponse(&finv1.TrackUsageResponse{
		CurrentMtdUsd: used,
	}), nil
}

func main() {
	w := whisper.New("FinOpsManager", "gcp_finops.lpsv")
	server := &FinOpsServer{
		logger: w,
		usage:  make(map[string]float64),
		budget: 100.0, // Default workstation budget pool
	}

	mux := http.NewServeMux()
	path, handler := finopsv1connect.NewFinOpsServiceHandler(server)
	mux.Handle(path, handler)

	port := "8098"
	slog.Info("FinOpsManager: Booting Economic Substrate...", "port", port)

	http.ListenAndServe(
		"localhost:"+port,
		h2c.NewHandler(mux, &http2.Server{}),
	)
}
func init() {
	_ = econotel.RecordEconomicEvent // Dependency check
}
