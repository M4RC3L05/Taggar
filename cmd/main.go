package main

import (
	"context"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"os/signal"
	"syscall"

	"github.com/m4rc3l05/taggar/internal/cli/edit"
	"github.com/m4rc3l05/taggar/internal/cli/view"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:           "taggar",
	Short:         "View and edit audio tags",
	Args:          cobra.ExactArgs(1),
	SilenceErrors: false,
	Version:       "0.1.0",
	SilenceUsage:  true,
}

func init() {
	rootCmd.AddCommand(view.NewCommand(), edit.NewCommand())
}

func main() {
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGABRT)

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
