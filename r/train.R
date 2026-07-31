# train.R
# Trains a Random Forest classifier using tidymodels and saves the workflow.

library(tidymodels)
library(ranger)

set.seed(42)

# Load data - replace with your own dataset
data <- iris |>
  mutate(target = ifelse(Species == "setosa", "yes", "no") |> factor())

split <- initial_split(data, prop = 0.8, strata = target)
train <- training(split)
test  <- testing(split)

# Save split for evaluation step
saveRDS(list(train = train, test = test), "r/data_split.rds")

# Recipe
rec <- recipe(target ~ Sepal.Length + Sepal.Width + Petal.Length + Petal.Width,
              data = train) |>
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
