# Continue extension + Ollama

With [Continue](https://www.continue.dev/) extension you can have a similar experience as you have with Github Copilot, but running the models locally on your machine. For that we need to use two different models, one specialized for chat and another one for autocompletion. The following instructions will help you on that.

The script [start-ollama-continue.sh](start-ollama-continue.sh) is able to run two different [Ollama](https://ollama.com) servers. The first one is for the chat and the second one is for the auto-complete. The first one will run in the port 11435 and the second one will run in port 11436. The script will run the two servers and wait until the user close the script with `CTRL+C` command.

Before start the script, you need to load the models, e.g.:

```sh
ollama pull llama3          # size around 5.8GB (default for continue extension)
ollama pull starcoder:3b    # size around 2.6GB (default for continue extension)
ollama pull starcoder:1b    # size around 1.7GB
ollama pull phi3            # size around 4.8GB
```

You can choose the models that can fit in your GPU memory for a good experience, e.g.:
- `phi3` + `starcoder:3b`: both models can fully fit in a GPU with 8GB GPU memory
- `llama3` + `starcoder:3b`: one of the models will partially fit in GPU with 8GB GPU memory

This is what happens when a model partially fit:

```sh 
$ OLLAMA_HOST=127.0.0.1:11436 ollama ps
NAME            ID              SIZE    PROCESSOR       UNTIL               
starcoder2:3b   f67ae0f64584    2.6 GB  14%/86% CPU/GPU 29 minutes from now
```

Using CPU is not recommended because we will experiment slow answers from the models.

This is what happens when a model can fully fit:

```sh
$ OLLAMA_HOST=127.0.0.1:11435 ollama ps
NAME            ID              SIZE    PROCESSOR       UNTIL               
llama3:latest   365c0bd3c000    5.8 GB  100% GPU        28 minutes from now
```

## Continue extension

For Visual Studio Code, the `config.json file is used to configure the Continue extension. For continue chat, you need to add the following entries to be able to point the models to your ollama server running locally on your machine.

```json
    {
      "title": "Llama 3 (Local)",
      "provider": "ollama",
      "model": "llama3",
      "apiBase": "http://127.0.0.1:11435"
    },
    {
      "title": "Phi 3 (Local)",
      "provider": "ollama",
      "model": "phi3",
      "apiBase": "http://127.0.0.1:11435"
    }
```

The available models will show in the chat window and you can select the one you want to use.

For autocomplete, you need to set the following entry to be able to point the desired model to your ollama server running locally on your machine.

```json
  "tabAutocompleteModel": {
    "title": "Starcoder 3b",
    "provider": "ollama",
    "model": "starcoder2:3b",
    "apiBase": "http://127.0.0.1:11436"
  },
```