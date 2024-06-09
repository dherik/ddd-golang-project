#!/bin/bash

OLLAMA_HOST=127.0.0.1:11435 OLLAMA_MODELS=/usr/share/ollama/.ollama/models HSA_OVERRIDE_GFX_VERSION=10.3.0 ollama serve &> /dev/null &
pid1=$!
echo "Initializing server for chat"

OLLAMA_HOST=127.0.0.1:11436 OLLAMA_MODELS=/usr/share/ollama/.ollama/models HSA_OVERRIDE_GFX_VERSION=10.3.0 ollama serve &> /dev/null &
pid2=$!
echo "Initializing server for autocomplete"

echo -n "Ollama servers started"

trap "kill -INT $pid1 $pid2" SIGINT
wait $pid1 $pid2
echo -n "Ollama servers stopped"