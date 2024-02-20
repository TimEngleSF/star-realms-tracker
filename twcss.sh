#!/bin/bash

# Define the input and output paths
INPUT="./css/input.css"
OUTPUT="./css/style.css"

# Function to start Tailwind CSS in watch mode
watch() {
   ./tailwindcss -i $INPUT -o $OUTPUT --watch
}

# Function to build and minify CSS for production
build() {
  ./tailwindcss -i $INPUT -o $OUTPUT --minify
}

# Check the first argument passed to the script
case "$1" in
  watch)
    watch
    ;;
  build)
    build
    ;;
  *)
    echo "Usage: $0 {watch|build}"
    exit 1
    ;;
esac
