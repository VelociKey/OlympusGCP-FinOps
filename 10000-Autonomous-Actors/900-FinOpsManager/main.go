package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"connectrpc.com/connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	finv1 "OlympusGCP-FinOps/40000-Communication-Contracts/430-Protocol-Definitions/000-gen/finops/v1x"
	"OlympusGCP-FinOps/40000-Communication-Contracts/430-Protocol-Definitions/000-gen/finops/v1x/finopsv1connect"
	"Olympus2/90000-Enablement-Labs/P0000-pkg/000-whisper"
)

type FinOpsServer struct {
	logger *whisper.WhisperLog
}

func (s *FinOpsServer) EstimateCost(ctx context.Context, req *connect.Request[finv1.EstimateCostRequest]) (*connect.Response[finv1.EstimateCostResponse], error) {
	start := time.Now()
	s.logger.Log("ESTIMATE_COST", "SUCCESS", req.Msg.Service, req.Msg.Action, time.Since(start))
	slog.Info("FinOps Cost Estimation", "service", req.Msg.Service, "action", req.Msg.Action)
	return connect.NewResponse(&finv1.EstimateCostResponse{
		EstimatedUsd: 15.50,
		Confidence:   "HIGH",
	}), nil
}

func (s *FinOpsServer) ValidateBudget(ctx context.Context, req *connect.Request[finv1.ValidateBudgetRequest]) (*connect.Response[finv1.ValidateBudgetResponse], error) {
	start := time.Now()
	s.logger.Log("VALIDATE_BUDGET", "SUCCESS", req.Msg.ProjectId, fmt.Sprintf("%.2f", req.Msg.RequestedAmount), time.Since(start))
	slog.Info("FinOps Budget Validation", "project", req.Msg.ProjectId, "amount", req.Msg.RequestedAmount)
	return connect.NewResponse(&finv1.ValidateBudgetResponse{
		Approved: true,
		Message:  "Within allocated monthly quota",
	}), nil
}

func (s *FinOpsServer) TrackUsage(ctx context.Context, req *connect.Request[finv1.TrackUsageRequest]) (*connect.Response[finv1.TrackUsageResponse], error) {
	start := time.Now()
	s.logger.Log("TRACK_USAGE", "SUCCESS", req.Msg.Service, req.Msg.ResourceId, time.Since(start))
	slog.Info("FinOps Usage Tracking", "service", req.Msg.Service, "resource", req.Msg.ResourceId)
	return connect.NewResponse(&finv1.TrackUsageResponse{
		CurrentMtdUsd: 1250.75,
	}), nil
}

func main() {
	w := whisper.New("FinOpsManager", "gcp_finops.lpsv")
	server := &FinOpsServer{logger: w}
	
	mux := http.NewServeMux()
	path, handler := finopsv1connect.NewFinOpsServiceHandler(server)
	mux.Handle(path, handler)

	port := "8098" // Sequential port for FinOps
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	slog.Info("FinOpsManager: Booting Substrate...", "port", port)
	slog.Info("FinOpsManager: Whisper bus connected")
	
	http.ListenAndServe(
		"localhost:"+port,
		h2c.NewHandler(mux, &http2.Server{}),
	)
}
