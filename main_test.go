package main

import (
	"os"
	"testing"
)

var defaultThresholds = Metrics{
	ROCAUC:   0.80,
	Accuracy: 0.75,
	FMeasure: 0.70,
}

func writeTempMetrics(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "metrics-*.json")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Remove(f.Name()) })
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestValidateMetrics_AllPass(t *testing.T) {
	path := writeTempMetrics(t, `{"roc_auc":0.90,"accuracy":0.85,"f_measure":0.80}`)
	if err := validateMetrics(path, defaultThresholds); err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestValidateMetrics_AtThreshold(t *testing.T) {
	path := writeTempMetrics(t, `{"roc_auc":0.80,"accuracy":0.75,"f_measure":0.70}`)
	if err := validateMetrics(path, defaultThresholds); err != nil {
		t.Errorf("expected no error at exact threshold, got: %v", err)
	}
}

func TestValidateMetrics_ROCAUCFails(t *testing.T) {
	path := writeTempMetrics(t, `{"roc_auc":0.75,"accuracy":0.85,"f_measure":0.80}`)
	if err := validateMetrics(path, defaultThresholds); err == nil {
		t.Error("expected error for ROC-AUC below threshold")
	}
}

func TestValidateMetrics_AccuracyFails(t *testing.T) {
	path := writeTempMetrics(t, `{"roc_auc":0.90,"accuracy":0.70,"f_measure":0.80}`)
	if err := validateMetrics(path, defaultThresholds); err == nil {
		t.Error("expected error for Accuracy below threshold")
	}
}

func TestValidateMetrics_FMeasureFails(t *testing.T) {
	path := writeTempMetrics(t, `{"roc_auc":0.90,"accuracy":0.85,"f_measure":0.65}`)
	if err := validateMetrics(path, defaultThresholds); err == nil {
		t.Error("expected error for F-Measure below threshold")
	}
}

func TestValidateMetrics_AllFail(t *testing.T) {
	path := writeTempMetrics(t, `{"roc_auc":0.70,"accuracy":0.60,"f_measure":0.50}`)
	if err := validateMetrics(path, defaultThresholds); err == nil {
		t.Error("expected error when all metrics fail")
	}
}

func TestValidateMetrics_MalformedJSON(t *testing.T) {
	path := writeTempMetrics(t, `not valid json`)
	if err := validateMetrics(path, defaultThresholds); err == nil {
		t.Error("expected error for malformed JSON")
	}
}

func TestValidateMetrics_FileNotFound(t *testing.T) {
	if err := validateMetrics("nonexistent.json", defaultThresholds); err == nil {
		t.Error("expected error for missing file")
	}
}
