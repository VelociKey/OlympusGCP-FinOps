package main

import (
	"context"
	"testing"

	finopsv1 "OlympusGCP-FinOps/gen/v1/finops"
	"OlympusGCP-FinOps/10000-Autonomous-Actors/10700-Processing-Engines/10710-Reasoning-Inference/inference"
	"connectrpc.com/connect"
)

func TestFinOpsServer(t *testing.T) {
	server := &inference.FinOpsServer{}
	ctx := context.Background()

	// Test ValidateBudget
	budgetRes, err := server.ValidateBudget(ctx, connect.NewRequest(&finopsv1.ValidateBudgetRequest{
		ProjectId:       "test-project",
		RequestedAmount: 100.0,
	}))
	if err != nil {
		t.Fatalf("ValidateBudget failed: %v", err)
	}
	if !budgetRes.Msg.Approved {
		t.Error("Expected budget to be approved")
	}

	// Test EstimateCost
	costRes, err := server.EstimateCost(ctx, connect.NewRequest(&finopsv1.EstimateCostRequest{
		Service: "storage",
		Action:  "upload",
	}))
	if err != nil {
		t.Fatalf("EstimateCost failed: %v", err)
	}
	if costRes.Msg.EstimatedUsd <= 0 {
		t.Errorf("Expected positive cost, got %f", costRes.Msg.EstimatedUsd)
	}

	// Test TrackUsage
	usageRes, err := server.TrackUsage(ctx, connect.NewRequest(&finopsv1.TrackUsageRequest{
		Service:          "compute",
		ConsumptionUnits: 10,
	}))
	if err != nil {
		t.Fatalf("TrackUsage failed: %v", err)
	}
	if usageRes.Msg.CurrentMtdUsd < 0 {
		t.Error("Expected non-negative MTD USD")
	}
}
