# evaluate.R
# Loads the trained workflow and computes evaluation metrics on the test set.

library(tidymodels)

wf    <- readRDS("r/model.rds")
split <- readRDS("r/data_split.rds")
test  <- split$test

preds <- augment(wf, new_data = test)

metrics <- metric_set(roc_auc, accuracy, f_meas)
results <- metrics(preds,
                   truth    = class,
                   estimate = .pred_class,
                   .pred_PS,
                   event_level = "first")

print(results)

saveRDS(results, "r/eval_results.rds")
cat("Evaluation complete. Results saved to r/eval_results.rds\n")
