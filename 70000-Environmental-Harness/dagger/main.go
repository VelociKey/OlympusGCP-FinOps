package main

import "context"
import "olympus.fleet/00SDLC/OlympusForge/70000-Environmental-Harness/dagger/olympusgcp-finops/internal/dagger"

type OlympusGCPFinOps struct{}

func (m *OlympusGCPFinOps) HelloWorld(ctx context.Context) string { return "Hello from OlympusGCP-FinOps!" }

func main() { dagger.Serve() }
