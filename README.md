# Structured Logging Demo in Go

This project demonstrates **structured logging** in Go using the `slog` package, showcasing both **top-down** and **bottom-up** logging approaches. The demo includes a simple client-server application where logs are captured and stored in a `logs.txt` file. You can use `jq` to filter and analyze the logs.

## Prerequisites

- Go 1.21 or higher
- `jq` (for filtering JSON logs)
- Make (for running commands via the `Makefile`)

## How to Use

### 1. **Running the Demo**
Run the demo, which starts a simple client-server application that logs various events to `logs.txt`.

To run the demo, use the following command:

```sh
make demo
```

This will:
- Start the **client-server application**.
- Logs will be written to `logs.txt` as the server handles requests.

### 2. **Viewing Logs**
Once the demo is running, you can **filter and view logs** using `jq`. For example, you can extract logs related to a specific message:

```sh
cat logs.txt | jq

```

Or filter logs for a specific log level:

```sh
cat logs.txt | jq 'select(.level == "error")'
```

Or look at logs live:

```sh
tail -f logs.txt | jq
```

### 3. **Presenting the Demo**
If you're ready to present the structured logging concepts, use the `make present` command. This will open the **presentation** in your browser.

To start the presentation:

```sh
make present
```

This will start a local server to serve the presentation using Go Present.

