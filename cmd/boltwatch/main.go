package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/user/boltwatch/internal/db"
	"github.com/user/boltwatch/internal/stats"
	"github.com/user/boltwatch/internal/watch"
)

func main() {
	interval := flag.Duration("interval", 2*time.Second, "poll interval (e.g. 1s, 500ms)")
	verbose := flag.Bool("verbose", false, "print raw stats to stderr")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: boltwatch [flags] <db-file>\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	dbPath := flag.Arg(0)
	inspector, err := db.Open(dbPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening db: %v\n", err)
		os.Exit(1)
	}
	defer inspector.Close()

	opts := watch.DefaultOptions().
		WithInterval(*interval).
		WithVerbose(*verbose)

	collector := stats.NewCollector(inspector)
	w := watch.NewWatcher(collector, opts.Interval)

	w.OnUpdate = func(buckets []stats.BucketStat) {
		fmt.Print("\033[H\033[2J") // clear terminal
		fmt.Println(stats.FormatTable(buckets))
		fmt.Println(stats.FormatSummary(buckets))
	}
	w.OnError = func(e error) {
		fmt.Fprintf(os.Stderr, "watch error: %v\n", e)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	fmt.Printf("Watching %s every %s — press Ctrl+C to quit\n", dbPath, opts.Interval)
	w.Start(ctx)
	fmt.Println("\nStopped.")
}
