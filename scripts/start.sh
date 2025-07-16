#!/bin/sh

# Start both servers in background
./http-server &
HTTP_PID=$!

./grpc-server &
GRPC_PID=$!

# Function to handle shutdown
shutdown() {
    echo "Shutting down servers..."
    kill $HTTP_PID $GRPC_PID
    wait $HTTP_PID $GRPC_PID
    echo "Servers stopped"
    exit 0
}

# Trap signals
trap shutdown SIGTERM SIGINT

# Wait for any process to exit
wait
