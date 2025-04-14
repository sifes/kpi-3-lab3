#!/bin/bash

curl -X POST http://localhost:17000 -d "reset"
curl -X POST http://localhost:17000 -d "green"
curl -X POST http://localhost:17000 -d "figure 0.9 0.1"
curl -X POST http://localhost:17000 -d "update"

while true; do
    curl -X POST http://localhost:17000 -d "update"
    for ((i = 0; i < 16; i++)); do
        curl -X POST http://localhost:17000 -d "move -0.05 0.05"
        curl -X POST http://localhost:17000 -d "update"
    done
    for ((i = 0; i < 16; i++)); do
        curl -X POST http://localhost:17000 -d "move 0.05 -0.05"
        curl -X POST http://localhost:17000 -d "update"
    done
done