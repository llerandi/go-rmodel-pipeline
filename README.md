# go-rmodel-pipeline

[![CI](https://img.shields.io/github/actions/workflow/status/llerandi/go-rmodel-pipeline/ci.yaml?label=CI&logo=github)](https://github.com/llerandi/go-rmodel-pipeline/actions/workflows/ci.yaml) [![License](https://img.shields.io/github/license/llerandi/go-rmodel-pipeline)](https://github.com/llerandi/go-rmodel-pipeline/blob/main/LICENSE) [![Stars](https://img.shields.io/github/stars/llerandi/go-rmodel-pipeline?style=social)](https://github.com/llerandi/go-rmodel-pipeline/stargazers) [![Last commit](https://img.shields.io/github/last-commit/llerandi/go-rmodel-pipeline)](https://github.com/llerandi/go-rmodel-pipeline/commits/main) [![Go](https://img.shields.io/badge/go-1.22%2B-blue?logo=go)](https://go.dev/) [![R](https://img.shields.io/badge/R-4.4-276DC3?logo=r)](https://www.r-project.org/)

Go-orchestrated pipeline for training, evaluating, and validating R machine learning models with threshold-based quality gates.

## How it works

```
go run main.go
    │
    ├── Rscript r/train.R           → trains model, saves r/model.rds
    ├── Rscript r/evaluate.R        → evaluates on test set, saves r/eval_results.rds
    ├── Rscript r/export_metrics.R  → writes metrics.json
    │
    └── Go reads metrics.json → checks thresholds → exit 0 / exit 1
```

## Quick start

```bash
# 1. Clone and enter the repo
git clone https://github.com/llerandi/go-rmodel-pipeline.git
cd go-rmodel-pipeline

# 2. Download Go dependencies
go mod tidy

# 3. Install R packages (once)
Rscript -e 'install.packages(c("tidymodels", "ranger", "modeldata", "jsonlite"))'

# 4. Add your data - edit r/train.R to load your own dataset

# 5. Adjust thresholds if needed - edit config.yaml

# 6. Run
go run main.go
```

Override any threshold at runtime without editing `config.yaml`:

```bash
go run main.go --roc-auc 0.85 --accuracy 0.80 --f-measure 0.75
```

## Requirements

- Go 1.22+
- R 4.4 with packages: `tidymodels`, `ranger`, `modeldata`, `jsonlite`

## Quality gates

| Metric    | Threshold |
|-----------|-----------|
| ROC-AUC   | ≥ 0.80    |
| Accuracy  | ≥ 0.75    |
| F-Measure | ≥ 0.70    |

If any metric falls below its threshold, the pipeline exits with code 1 - making it suitable as a CI/CD gate.

## Sample output

After a successful run, `metrics.json` looks like:

```json
{
  "roc_auc": 0.9247,
  "accuracy": 0.8611,
  "f_measure": 0.8034
}
```

## Project structure

```
.
├── main.go                 # Go orchestrator and validator
├── go.mod
├── r/
│   ├── train.R             # Model training
│   ├── evaluate.R          # Model evaluation
│   └── export_metrics.R    # Export results to metrics.json
└── metrics.json            # Generated - gitignored
```

## Customisation

- **Swap the dataset**: edit `r/train.R` to load your own data.
- **Change thresholds**: edit `config.yaml`, or use `--roc-auc`, `--accuracy`, `--f-measure` flags at runtime.
- **Add metrics**: extend the `Metrics` struct in `main.go` and the export in `r/export_metrics.R`.

## Roadmap

**Pipeline**
- [x] Go orchestrator that runs R scripts in sequence
- [x] R scripts for training, evaluation, and metrics export
- [x] Threshold-based quality gates in Go
- [x] Configurable thresholds via `config.yaml` with CLI flag overrides
- [x] R example dataset (`modeldata::cells` - binary classification)
- [ ] Structured logging in Go (JSON output, log levels)
- [ ] Timeout handling per R script
- [ ] Support for regression models, not just classification
- [ ] Model versioning: persist metrics history per run

**Testing**
- [x] GitHub Actions CI - Go build and vet
- [x] GitHub Actions CI - R syntax check
- [x] Go tests for `validateMetrics`
- [ ] Integration test that runs the full pipeline on the example dataset
- [ ] R unit tests with `testthat`

**Documentation**
- [x] Shields in README (CI, License, Stars, Last commit, Go, R)
- [x] `.gitignore` for R artifacts and `metrics.json`
- [x] Sample `metrics.json` in README

**Nice to have**
- [ ] Docker image with Go + R + dependencies
- [ ] GitHub Actions workflow that runs the full pipeline on a sample dataset
- [ ] Compare metrics against previous run (model regression detection)
