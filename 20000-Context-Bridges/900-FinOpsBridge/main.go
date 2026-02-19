package main

import (
	"context"
	"fmt"
// 	"log"
	"net/http"

	"connectrpc.com/connect"
	"github.com/mark3labs/mcp-go/mcp"

	"OlympusGCP-FinOps/40000-Communication-Contracts/430-Protocol-Definitions/000-gen/finops/v1x/finopsv1connect"
	finopsv1 "OlympusGCP-FinOps/40000-Communication-Contracts/430-Protocol-Definitions/000-gen/finops/v1x"
	"Olympus2/90000-Enablement-Labs/P0000-pkg/000-mcp-bridge"
)

func main() {
	s := mcpbridge.NewBridgeServer("OlympusFinOpsBridge", "1.0.0")

	client := finopsv1connect.NewFinOpsServiceClient(
		http.DefaultClient,
		"http://localhost:8093",
	)

	s.AddTool(mcp.NewTool("finops_estimate_cost",
		mcp.WithDescription("Estimate the GCP cost of a request. Args: {service: string, action: string}"),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		m, err := mcpbridge.ExtractMap(request)
		if err != nil {
			return mcpbridge.HandleError(err)
		}

		service, _ := m["service"].(string)
		action, _ := m["action"].(string)

		resp, err := client.EstimateCost(ctx, connect.NewRequest(&finopsv1.EstimateCostRequest{
			Service: service,
			Action:  action,
		}))
		if err != nil {
			return mcpbridge.HandleError(err)
		}

		return mcp.NewToolResultText(fmt.Sprintf("Estimated Cost: $%.2f (Confidence: %s)", resp.Msg.EstimatedUsd, resp.Msg.Confidence)), nil
	})

	s.AddTool(mcp.NewTool("finops_validate_budget",
		mcp.WithDescription("Check if an operation fits within local safety budget. Args: {amount: number}"),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		m, err := mcpbridge.ExtractMap(request)
		if err != nil {
			return mcpbridge.HandleError(err)
		}

		amount, _ := m["amount"].(float64)

		resp, err := client.ValidateBudget(ctx, connect.NewRequest(&finopsv1.ValidateBudgetRequest{
			ProjectId:       "local-dev",
			RequestedAmount: amount,
		}))
		if err != nil {
			return mcpbridge.HandleError(err)
		}

		return mcp.NewToolResultText(fmt.Sprintf("Budget Approval: %t. Message: %s", resp.Msg.Approved, resp.Msg.Message)), nil
	})

	s.Run()
}
