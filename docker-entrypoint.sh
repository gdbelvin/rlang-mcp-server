#!/bin/bash
set -e

# Check if an R Markdown file was provided
if [ "$1" = "render" ] && [ -n "$2" ]; then
  # Render the R Markdown file
  Rscript -e "rmarkdown::render('$2', output_format='html_document', output_dir='/rmd/output')"
  
  # Output the path to the rendered file
  echo "Rendered file: /rmd/output/$(basename ${2%.*}).html"
elif [ "$1" = "render_pdf" ] && [ -n "$2" ]; then
  # Render the R Markdown file to PDF
  Rscript -e "rmarkdown::render('$2', output_format='pdf_document', output_dir='/rmd/output')"
  
  # Output the path to the rendered file
  echo "Rendered file: /rmd/output/$(basename ${2%.*}).pdf"
elif [ "$1" = "render_word" ] && [ -n "$2" ]; then
  # Render the R Markdown file to Word
  Rscript -e "rmarkdown::render('$2', output_format='word_document', output_dir='/rmd/output')"
  
  # Output the path to the rendered file
  echo "Rendered file: /rmd/output/$(basename ${2%.*}).docx"
else
  # If no command is provided, start an R session
  exec R "$@"
fi
