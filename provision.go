package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	cluster := "FinOps"

	fmt.Printf("Starting %s Provisioner...\n", cluster)
	targets := []struct{ name, path string }{
		{"FinOpsManager", "10000-Autonomous-Actors/900-FinOpsManager"},
		{"FinOpsBridge", "20000-Context-Bridges/900-FinOpsBridge"},
	}
	for _, t := range targets {
		fmt.Printf("Building %s via Dagger...\n", t.name)

		// Execute Dagger workstation-native build
		cmd := exec.Command("dagger", "call", "build")
		cmd.Dir = "." // Cluster root
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("Dagger build failed for %s: %v\n", t.name, err)
			os.Exit(1)
		}
	}

	fmt.Printf("%s built successfully.\n", cluster)
}
