package main

import (
	"fmt"

	"os"

	"os/exec"

	"path/filepath"

"runtime"
)

func main() {

	cluster := "FinOps"

	fmt.Printf("Starting %s Provisioner...\n", cluster)
	targets := []struct{ name, path string }{

		{"FinOpsManager", "10000-Autonomous-Actors/900-FinOpsManager"},

		{"FinOpsBridge", "20000-Context-Bridges/900-FinOpsBridge"},
	}
	for _, t := range targets {

		fmt.Printf("Building %s...\n", t.name)
		binaryName := filepath.Base(t.path)

		if runtime.GOOS == "windows" {
			binaryName += ".exe"
		}

		cmd := exec.Command("go", "build", "-o", binaryName, "main.go")
		cmd.Dir = filepath.Join(cluster, t.path)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			os.Exit(1)
		}
	}

	fmt.Printf("%s built successfully.\n", cluster)
}
