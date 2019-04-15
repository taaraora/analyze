#!/bin/bash
set -ex

echo "Running linters"

make lint
