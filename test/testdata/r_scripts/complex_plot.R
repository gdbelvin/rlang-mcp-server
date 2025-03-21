# Complex ggplot2 visualization
library(ggplot2)

# Create a more complex plot with facets, colors, and smoothing
ggplot(mtcars, aes(x = mpg, y = hp, color = factor(cyl))) +
  geom_point(size = 3, alpha = 0.7) +
  geom_smooth(method = "lm", se = TRUE, linetype = "dashed") +
  facet_wrap(~gear, nrow = 1, labeller = label_both) +
  scale_color_brewer(palette = "Set1", name = "Cylinders") +
  theme_bw() +
  theme(
    plot.title = element_text(size = 16, face = "bold", hjust = 0.5),
    plot.subtitle = element_text(size = 12, hjust = 0.5),
    axis.title = element_text(size = 12, face = "bold"),
    legend.position = "bottom",
    strip.background = element_rect(fill = "lightgray"),
    strip.text = element_text(face = "bold")
  ) +
  labs(
    title = "MPG vs Horsepower by Cylinder Count",
    subtitle = "Faceted by Gear Count",
    x = "Miles Per Gallon",
    y = "Horsepower",
    caption = "Source: mtcars dataset"
  )
