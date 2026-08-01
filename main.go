package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"

	"gopkg.in/yaml.v3"
)

// Config holds the contents of config.yaml.
type Config struct {
	Thresholds struct {
		ROCAUC   float64 `yaml:"roc_auc"`
		Accuracy float64 `yaml:"accuracy"`
		FMeasure float64 `yaml:"f_measure"`
	} `yaml:"thresholds"`
}

// Metrics holds the model evaluation results exported by R.
type Metrics struct {
	ROCAUC   float64 `json:"roc_auc"`
	Accuracy float64 `json:"accuracy"`
	FMeasure float64 `json:"f_measure"`
}

func loadConfig(path string) (Config, error) {
	var cfg Config
	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, fmt.Errorf("could not read %s: %w", path, err)
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, fmt.Errorf("could not parse %s: %w", path, err)
	}
	return cfg, nil
}

func runRScript(script string) error {
	fmt.Printf("\n>>> Running %s...\n", script)
	cmd := exec.Command("Rscript", script)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func validateMetrics(path string, thresholds Metrics) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("could not read %s: %w", path, err)
	}

	var m Metrics
	if err := json.Unmarshal(data, &m); err != nil {
		return fmt.Errorf("could not parse metrics: %w", err)
	}

	fmt.Printf("\n>>> Metrics\n")
	fmt.Printf("  ROC-AUC:   %.4f (threshold: %.2f)\n", m.ROCAUC, thresholds.ROCAUC)
	fmt.Printf("  Accuracy:  %.4f (threshold: %.2f)\n", m.Accuracy, thresholds.Accuracy)
	fmt.Printf("  F-Measure: %.4f (threshold: %.2f)\n", m.FMeasure, thresholds.FMeasure)

	var failures []string
	if m.ROCAUC < thresholds.ROCAUC {
		failures = append(failures, fmt.Sprintf("ROC-AUC %.4f < %.2f", m.ROCAUC, thresholds.ROCAUC))
	}
	if m.Accuracy < thresholds.Accuracy {
		failures = append(failures, fmt.Sprintf("Accuracy %.4f < %.2f", m.Accuracy, thresholds.Accuracy))
	}
	if m.FMeasure < thresholds.FMeasure {
		failures = append(failures, fmt.Sprintf("F-Measure %.4f < %.2f", m.FMeasure, thresholds.FMeasure))
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
	flagROCAUC   := flag.Float64("roc-auc", -1, "override ROC-AUC threshold from config.yaml")
	flagAccuracy := flag.Float64("accuracy", -1, "override Accuracy threshold from config.yaml")
	flagFMeasure := flag.Float64("f-measure", -1, "override F-Measure threshold from config.yaml")
	flag.Parse()

	cfg, err := loadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	thresholds := Metrics{
		ROCAUC:   cfg.Thresholds.ROCAUC,
		Accuracy: cfg.Thresholds.Accuracy,
		FMeasure: cfg.Thresholds.FMeasure,
	}

	if *flagROCAUC >= 0 {
		thresholds.ROCAUC = *flagROCAUC
	}
	if *flagAccuracy >= 0 {
		thresholds.Accuracy = *flagAccuracy
	}
	if *flagFMeasure >= 0 {
		thresholds.FMeasure = *flagFMeasure
	}

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

	if err := validateMetrics("metrics.json", thresholds); err != nil {
		log.Fatal(err)
	}
}
