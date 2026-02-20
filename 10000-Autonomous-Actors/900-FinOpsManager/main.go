package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
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

type SkuConfig struct {
	Skus map[string]float64 `json:"skus"`
}

type FinOpsServer struct {
	logger *whisper.WhisperLog
	mu     sync.Mutex
	usage  map[string]float64
	budget float64
	skus   map[string]float64
}

func (s *FinOpsServer) EstimateCost(ctx context.Context, req *connect.Request[finv1.EstimateCostRequest]) (*connect.Response[finv1.EstimateCostResponse], error) {
	start := time.Now()

	price, ok := s.skus[req.Msg.Service]
	if !ok {
		price = 0.01 // Default safety price
	}

	s.logger.Log("ESTIMATE_COST", "SUCCESS", req.Msg.Service, "SKU_LOOKUP", time.Since(start))
	return connect.NewResponse(&finv1.EstimateCostResponse{
		EstimatedUsd: price,
		Confidence:   "HIGH_FIDELITY_LOCAL",
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

	price, ok := s.skus[req.Msg.Service]
	if !ok {
		price = 0.01
	}

	totalCost := req.Msg.ConsumptionUnits * price
	s.usage[req.Msg.Service] += totalCost

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

	// Load SKU Database
	skuData, err := os.ReadFile("C0100-Configuration-Registry/settings/sku_pricing.json")
	var skus SkuConfig
	if err == nil {
		json.Unmarshal(skuData, &skus)
	} else {
		slog.Warn("FinOpsManager: Failed to load SKU database, using empty map", "error", err)
		skus.Skus = make(map[string]float64)
	}

	server := &FinOpsServer{
		logger: w,
		usage:  make(map[string]float64),
		budget: 500.0, // Expanded workstation budget
		skus:   skus.Skus,
	}

	mux := http.NewServeMux()
	path, handler := finopsv1connect.NewFinOpsServiceHandler(server)
	mux.Handle(path, handler)

	port := "8098"
	slog.Info("FinOpsManager: Booting High-Fidelity Economic Substrate...", "port", port)

	http.ListenAndServe(
		"localhost:"+port,
		h2c.NewHandler(mux, &http2.Server{}),
	)
}
func init() {
	_ = econotel.RecordEconomicEvent
}
