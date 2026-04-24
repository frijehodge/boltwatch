# boltwatch

A lightweight CLI for monitoring BoltDB files and visualizing bucket statistics in real time.

---

## Installation

```bash
go install github.com/yourusername/boltwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/boltwatch.git
cd boltwatch
go build -o boltwatch .
```

---

## Usage

Point `boltwatch` at any BoltDB file to start monitoring:

```bash
boltwatch --file ./mydata.db
```

### Options

| Flag | Default | Description |
|------|---------|-------------|
| `--file` | _(required)_ | Path to the BoltDB file |
| `--interval` | `2s` | Refresh interval for stats |
| `--bucket` | _(all)_ | Filter output to a specific bucket |

### Example Output

```
Watching: ./mydata.db  (refresh: 2s)
─────────────────────────────────────────
Bucket              Keys     Size
─────────────────────────────────────────
users               1,204    4.2 MB
sessions            389      1.1 MB
events              8,901    22.7 MB
─────────────────────────────────────────
Last updated: 2024-05-10 14:32:01
```

---

## Requirements

- Go 1.21+
- A valid [BoltDB](https://github.com/etcd-io/bbolt) (bbolt) database file

---

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

---

## License

[MIT](LICENSE)