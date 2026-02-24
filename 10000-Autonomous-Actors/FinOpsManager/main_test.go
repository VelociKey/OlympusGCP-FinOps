package main

import (
	"context"
	"testing"

	finopsv1 "OlympusGCP-FinOps/40000-Communication-Contracts/430-Protocol-Definitions/000-gen/finops/v1"
	"connectrpc.com/connect"
)

func TestFinOpsServer(t *testing.T) {
	server := &FinOpsServer{}
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
	if costRes.Msg.EstimatedCost <= 0 {
		t.Errorf("Expected positive cost, got %f", costRes.Msg.EstimatedCost)
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
