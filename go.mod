module OlympusGCP-FinOps

go 1.25.7

// Local Resolution
replace Olympus2 => ../Olympus2

replace OlympusActors-Cognition => ../OlympusActors-Cognition

replace OlympusActors-Delegation => ../OlympusActors-Delegation

replace OlympusAscent => ../OlympusAscent

replace OlympusAssurance => ../OlympusAssurance

replace OlympusAtelier => ../OlympusAtelier

replace OlympusFabric => ../OlympusFabric

replace OlympusForge => ../OlympusForge

replace OlympusGCP-Compute => ../OlympusGCP-Compute

replace OlympusGCP-Data => ../OlympusGCP-Data

replace OlympusGCP-Events => ../OlympusGCP-Events

replace OlympusGCP-Firebase => ../OlympusGCP-Firebase

replace OlympusGCP-Intelligence => ../OlympusGCP-Intelligence

replace OlympusGCP-Messaging => ../OlympusGCP-Messaging

replace OlympusGCP-Observability => ../OlympusGCP-Observability

replace OlympusGCP-Storage => ../OlympusGCP-Storage

replace OlympusGCP-Vault => ../OlympusGCP-Vault

replace OlympusGrammar => ../OlympusGrammar

replace OlympusInfrastructure => ../OlympusInfrastructure

replace OlympusVision => ../OlympusVision

replace text => ../Olympus2/00000-Identity-Foundations/P0000-pkg/text

replace pretty => ../Olympus2/00000-Identity-Foundations/P0000-pkg/pretty

replace go-internal => ../Olympus2/00000-Identity-Foundations/P0000-pkg/go-internal

replace check.v1 => ../Olympus2/00000-Identity-Foundations/P0000-pkg/check.v1

replace gopkg.in/check.v1 => ../Olympus2/00000-Identity-Foundations/P0000-pkg/check.v1

require (
	Olympus2 v0.0.0-00010101000000-000000000000
	connectrpc.com/connect v1.19.1
	github.com/mark3labs/mcp-go v0.44.0
	golang.org/x/net v0.50.0
	google.golang.org/protobuf v1.36.11
)

require (
	github.com/bahlo/generic-list-go v0.2.0 // indirect
	github.com/buger/jsonparser v1.1.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/invopop/jsonschema v0.13.0 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/rogpeppe/go-internal v1.13.1 // indirect
	github.com/spf13/cast v1.7.1 // indirect
	github.com/stretchr/testify v1.11.1 // indirect
	github.com/wk8/go-ordered-map/v2 v2.1.8 // indirect
	github.com/yosida95/uritemplate/v3 v3.0.2 // indirect
	golang.org/x/text v0.34.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
