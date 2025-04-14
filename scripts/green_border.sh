#!/bin/bash

SERVER_URL="http://localhost:17000"

curl -s -X POST "$SERVER_URL" -d "reset"
curl -s -X POST "$SERVER_URL" -d "white"
curl -s -X POST "$SERVER_URL" -d "bgrect 0.25 0.25 0.75 0.75"
curl -s -X POST "$SERVER_URL" -d "figure 0.5 0.5"
curl -s -X POST "$SERVER_URL" -d "green"
curl -s -X POST "$SERVER_URL" -d "figure 0.6 0.6"
curl -s -X POST "$SERVER_URL" -d "update"
