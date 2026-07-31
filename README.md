# go-rmodel-pipeline

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

## Requirements

- Go 1.22+
- R 4.x with the following packages: `tidymodels`, `ranger`, `jsonlite`

Install R packages:
```r
install.packages(c("tidymodels", "ranger", "jsonlite"))
```

## Usage

```bash
go run main.go
```

## Quality gates

| Metric    | Threshold |
|-----------|-----------|
| ROC-AUC   | ≥ 0.80    |
| Accuracy  | ≥ 0.75    |
| F-Measure | ≥ 0.70    |

If any metric falls below its threshold, the pipeline exits with code 1 - making it suitable as a CI/CD gate.

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
- **Change thresholds**: edit the `Thresholds` variable in `main.go`.
- **Add metrics**: extend the `Metrics` struct in `main.go` and the export in `r/export_metrics.R`.
