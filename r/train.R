# train.R
# Trains a Random Forest classifier using tidymodels and saves the workflow.
# Dataset: modeldata::cells (binary classification - cell segmentation quality).

library(tidymodels)
library(modeldata)
library(ranger)

set.seed(42)

data("cells", package = "modeldata")

# Remove case column (train/test indicator from the original study)
cells <- cells |> select(-case)

split <- initial_split(cells, prop = 0.8, strata = class)
train <- training(split)
test  <- testing(split)

# Save split for evaluation step
saveRDS(list(train = train, test = test), "r/data_split.rds")

# Recipe
rec <- recipe(class ~ ., data = train) |>
  step_normalize(all_numeric_predictors())

# Model spec
spec <- rand_forest(trees = 500) |>
  set_engine("ranger") |>
  set_mode("classification")

# Workflow
wf <- workflow() |>
  add_recipe(rec) |>
  add_model(spec) |>
  fit(data = train)

saveRDS(wf, "r/model.rds")
cat("Training complete. Model saved to r/model.rds\n")
