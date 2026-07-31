package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
)

// Metrics holds the model evaluation results exported by R.
type Metrics struct {
	ROCAUC   float64 `json:"roc_auc"`
	Accuracy float64 `json:"accuracy"`
	FMeasure float64 `json:"f_measure"`
}

// Thresholds defines the minimum acceptable values for each metric.
var Thresholds = Metrics{
	ROCAUC:   0.80,
	Accuracy: 0.75,
	FMeasure: 0.70,
}

func runRScript(script string) error {
	fmt.Printf("\n>>> Running %s...\n", script)
	cmd := exec.Command("Rscript", script)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func validateMetrics(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("could not read %s: %w", path, err)
	}

	var m Metrics
	if err := json.Unmarshal(data, &m); err != nil {
		return fmt.Errorf("could not parse metrics: %w", err)
	}

	fmt.Printf("\n>>> Metrics\n")
	fmt.Printf("  ROC-AUC:   %.4f (threshold: %.2f)\n", m.ROCAUC, Thresholds.ROCAUC)
	fmt.Printf("  Accuracy:  %.4f (threshold: %.2f)\n", m.Accuracy, Thresholds.Accuracy)
	fmt.Printf("  F-Measure: %.4f (threshold: %.2f)\n", m.FMeasure, Thresholds.FMeasure)

	var failures []string
	if m.ROCAUC < Thresholds.ROCAUC {
		failures = append(failures, fmt.Sprintf("ROC-AUC %.4f < %.2f", m.ROCAUC, Thresholds.ROCAUC))
	}
	if m.Accuracy < Thresholds.Accuracy {
		failures = append(failures, fmt.Sprintf("Accuracy %.4f < %.2f", m.Accuracy, Thresholds.Accuracy))
	}
	if m.FMeasure < Thresholds.FMeasure {
		failures = append(failures, fmt.Sprintf("F-Measure %.4f < %.2f", m.FMeasure, Thresholds.FMeasure))
	}

	if len(failures) > 0 {
		fmt.Println("\n>>> Quality gate FAILED:")
		for _, f := range failures {
			fmt.Printf("  - %s\n", f)
		}
		return fmt.Errorf("model did not pass quality gates")
	}

	fmt.Println("\n>>> Quality gate PASSED")
	return nil
}

func main() {
	scripts := []string{
		"r/train.R",
		"r/evaluate.R",
		"r/export_metrics.R",
	}

	for _, script := range scripts {
		if err := runRScript(script); err != nil {
			log.Fatalf("Error running %s: %v", script, err)
		}
	}

	if err := validateMetrics("metrics.json"); err != nil {
		log.Fatal(err)
	}
}
