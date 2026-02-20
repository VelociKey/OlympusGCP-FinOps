package main

import "context"
import "dagger/olympusgcp-finops/internal/dagger"

type OlympusGCPFinOps struct{}

func (m *OlympusGCPFinOps) HelloWorld(ctx context.Context) string { return "Hello from OlympusGCP-FinOps!" }

func main() { dagger.Serve() }
