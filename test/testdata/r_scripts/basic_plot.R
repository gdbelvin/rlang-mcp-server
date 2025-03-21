# Basic ggplot2 visualization
library(ggplot2)

# Create a simple scatter plot
ggplot(mtcars, aes(x = mpg, y = hp)) +
  geom_point() +
  theme_minimal() +
  labs(
    title = "MPG vs Horsepower",
    x = "Miles Per Gallon",
    y = "Horsepower",
    caption = "Source: mtcars dataset"
  )
