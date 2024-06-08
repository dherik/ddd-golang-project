#!/bin/bash

# Instructions
# 
# Before start the script, you need to load the model before running this script, e.g.,
# $ ollama pull llama3          # size around 5.8GB (default for continue extension)
# $ ollama pull starcoder:3b    # size around 2.6GB (default for continue extension)
# $ ollama pull starcoder:1b    # size around 1.7GB
# $ ollama pull phi3            # size around 4.8GB
#
# You can choose the models that can fit in your GPU memory, e.g.,
# phi3 + starcoder:3b can fully fit in a GPU with 8GB GPU memory
# llama3 + starcoder:3b can partially fit in GPU with 8GB GPU memory, but not fully. This is what happens in this case:
# 
# $ OLLAMA_HOST=127.0.0.1:11436 ollama ps
# NAME            ID              SIZE    PROCESSOR       UNTIL               
# starcoder2:3b   f67ae0f64584    2.6 GB  14%/86% CPU/GPU 29 minutes from now
# 
# $ OLLAMA_HOST=127.0.0.1:11435 ollama ps
# NAME            ID              SIZE    PROCESSOR       UNTIL               
# llama3:latest   365c0bd3c000    5.8 GB  100% GPU        28 minutes from now

OLLAMA_HOST=127.0.0.1:11435 OLLAMA_MODELS=/usr/share/ollama/.ollama/models HSA_OVERRIDE_GFX_VERSION=10.3.0 ollama serve &> /dev/null &
pid1=$!
echo "Initializing server for chat"

OLLAMA_HOST=127.0.0.1:11436 OLLAMA_MODELS=/usr/share/ollama/.ollama/models HSA_OVERRIDE_GFX_VERSION=10.3.0 ollama serve &> /dev/null &
pid2=$!
echo "Initializing server for autocomplete"

echo -n "servers started"

trap "kill -INT $pid1 $pid2" SIGINT
wait $pid1 $pid2
echo -n "servers stopped"