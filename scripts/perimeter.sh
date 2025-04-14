#!/bin/bash
# Рух фігури квадратом із затримками

SERVER_URL="http://localhost:17000"

curl -s -X POST "$SERVER_URL" -d "reset"
curl -s -X POST "$SERVER_URL" -d "green"
curl -s -X POST "$SERVER_URL" -d "bgrect 0.4 0.4 0.6 0.6"
curl -s -X POST "$SERVER_URL" -d "figure 0.2 0.2"

move_and_update() {
    curl -s -X POST $SERVER_URL -d "move $1 $2"
    curl -s -X POST $SERVER_URL -d "update"
    sleep 0.5
}

curl -s -X POST $SERVER_URL -d "update"
sleep 1

move_and_update 0 0.6
move_and_update 0.6 0
move_and_update 0 -0.6
move_and_update -0.6 0
