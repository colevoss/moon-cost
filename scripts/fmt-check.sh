#! /usr/bin/env bash

lines=$(gofmt -l .)

if [ ! -z "$lines" ]; then
  echo "Invalid go formatting. Please fix"
  echo "$lines"
  exit 1
fi
