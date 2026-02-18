package main

import (
	"context"
	"fmt"
"log"
	"net/http"

	"connectrpc.com/connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"OlympusGCP-FinOps/gen/finops/v1x/finopsv1connect"
	finopsv1 "OlympusGCP-FinOps/gen/finops/v1x"
)

type FinOpsServer struct{}

func (s *FinOpsServer) EstimateCost(
	ctx context.Context,
	req *connect.Request[finopsv1.EstimateCostRequest],
) (*connect.Response[finopsv1.EstimateCostResponse], error) {
	log.Printf("FinOpsManager: Estimating cost for %s:%s", req.Msg.Service, req.Msg.Action)

	// Dynamic Cost Engine: Heuristics based on size or complexity
	estimate := 0.01 // Minimum

	sizeStr := req.Msg.Parameters["size_bytes"]
	var size float64
	fmt.Sscanf(sizeStr, "%f", &size)

	switch req.Msg.Service {
	case "bigquery":
		// $5 per Terabyte (simplified to bytes)
		estimate = (size / 1e12) * 5.0
		if estimate < 0.01 {
			estimate = 0.01
		}
	case "storage":
		// $0.02 per GB per month (simplified to immediate operation cost)
		estimate = (size / 1e9) * 0.02
	case "cloudrun":
		estimate = 0.00001 * size // Mock: cost per byte processed
	}

	return connect.NewResponse(&finopsv1.EstimateCostResponse{
		EstimatedUsd: estimate,
		Confidence:   "CALCULATED (SIZE_AWARE)",
	}), nil
}

func (s *FinOpsServer) ValidateBudget(
	ctx context.Context,
	req *connect.Request[finopsv1.ValidateBudgetRequest],
) (*connect.Response[finopsv1.ValidateBudgetResponse], error) {
	log.Printf("FinOpsManager: Validating budget for project %s", req.Msg.ProjectId)

	approved := true
	message := "Budget validated."
	if req.Msg.RequestedAmount > 100.0 {
		approved = false
		message = "Requested amount exceeds local developer safety threshold ($100)."
	}

	return connect.NewResponse(&finopsv1.ValidateBudgetResponse{
		Approved: approved,
		Message:  message,
	}), nil
}

func main() {
	manager := &FinOpsServer{}
	mux := http.NewServeMux()
	path, handler := finopsv1connect.NewFinOpsServiceHandler(manager)
	mux.Handle(path, handler)

	fmt.Println("FinOpsManager (ConnectRPC) starting on :8093")
	http.ListenAndServe(
		"localhost:8093",
		h2c.NewHandler(mux, &http2.Server{}),
	)
}
