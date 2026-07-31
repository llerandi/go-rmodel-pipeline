# export_metrics.R
# Exports evaluation results to metrics.json for Go validation.

library(tidyr)
library(jsonlite)

results <- readRDS("r/eval_results.rds")

metrics <- results |>
  select(.metric, .estimate) |>
  pivot_wider(names_from = .metric, values_from = .estimate) |>
  rename(f_measure = f_meas)

write_json(metrics, "metrics.json", auto_unbox = TRUE, digits = 6)
cat("Metrics exported to metrics.json\n")
