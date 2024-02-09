#! /bin/bash
echo "Running game config script..."

loc=$(pwd)

time=$(date)

mkdir -p "$loc/logs/$time"

sed -i "s/\"game_server\":.*/\"game_server\": \"${GAME_HOST}\",/g" config.json
sed -i "s/\"game_port\":.*/\"game_port\": \"${GAME_INNER}\",/g" config.json
sed -i "s/\"registry_server\":.*/\"registry_server\": \"${SERVER_HOST}\",/g" config.json
sed -i "s/\"registry_port\":.*/\"registry_port\": \"${SERVER_OUTER}\",/g" config.json
sed -i "s/\"game_mode\":.*/\"game_mode\": \"${GAME_MODE}\"/g" config.json

buildLog="$loc/logs/$time/build.log"
runLog="$loc/logs/$time/run.log"

echo "Building..."
go build >> "$buildLog"
echo "Done."

echo "Starting ${GAME_MODE} game..."
go run main.go >> "$runLog"
echo "Game stopped."