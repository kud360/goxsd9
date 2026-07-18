// Command specsync downloads and converts the pinned specification manifest.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/kud360/goxsd9/internal/specsync"
)

func main() {
	if err := run(context.Background(), os.Args[1:], os.Stderr); err != nil {
		if _, writeErr := fmt.Fprintf(os.Stderr, "specsync: %v\n", err); writeErr != nil {
			os.Exit(1)
		}
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string, stderr io.Writer) error {
	flags := flag.NewFlagSet("specsync", flag.ContinueOnError)
	flags.SetOutput(stderr)
	manifestPath := flags.String("manifest", "specs/manifest.json", "path to the pinned JSON manifest")
	cacheDir := flags.String("cache", "specs/cache", "directory for verified source XHTML")
	outputDir := flags.String("out", "specs/markdown", "directory for generated Markdown")
	if err := flags.Parse(args); err != nil {
		return fmt.Errorf("parse command arguments: %w", err)
	}

	manifest, err := readManifest(*manifestPath, func(path string) (io.ReadCloser, error) {
		return os.Open(path)
	})
	if err != nil {
		return err
	}

	options := specsync.Options{
		CacheDir:  *cacheDir,
		OutputDir: *outputDir,
		IDs:       flags.Args(),
		Client:    specsync.NewHTTPClient(),
	}
	if err := specsync.Sync(ctx, manifest, options); err != nil {
		return err
	}
	return nil
}

type manifestOpener func(string) (io.ReadCloser, error)

func readManifest(path string, open manifestOpener) (specsync.Manifest, error) {
	if open == nil {
		return specsync.Manifest{}, fmt.Errorf("open manifest %q: opener must not be nil", path)
	}
	manifestFile, err := open(path)
	if err != nil {
		return specsync.Manifest{}, fmt.Errorf("open manifest %q: %w", path, err)
	}
	if manifestFile == nil {
		return specsync.Manifest{}, fmt.Errorf("open manifest %q: opener returned a nil reader", path)
	}
	manifest, readErr := specsync.ReadManifest(manifestFile)
	closeErr := manifestFile.Close()
	if readErr == nil && closeErr == nil {
		return manifest, nil
	}
	var causes []error
	if readErr != nil {
		causes = append(causes, fmt.Errorf("read manifest %q: %w", path, readErr))
	}
	if closeErr != nil {
		causes = append(causes, fmt.Errorf("close manifest %q: %w", path, closeErr))
	}
	return specsync.Manifest{}, errors.Join(causes...)
}
