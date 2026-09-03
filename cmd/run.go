/*
Copyright © 2026
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/padiazg/broken-hexagon/internal/core/services"
	"github.com/padiazg/broken-hexagon/pkg/logger"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Start the broken-hexagon service",
	Long: `Start the broken-hexagon service and run until stopped.

The service will run continuously processing events, messages, or performing
scheduled tasks. It listens for SIGINT (Ctrl+C) and SIGTERM signals
and performs a graceful shutdown with a configurable timeout.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := GetConfig()

		// Initialize logger from config
		log := logger.New(&logger.Config{
			Level:  cfg.LogLevel,
			Format: cfg.LogFormat,
		})

		log.Info("Starting broken-hexagon service...")

		// Create your service/processor
		processor, err := services.NewProcessor(&services.ProcessorConfig{
			Config: cfg,
			Logger: log,
		})
		if err != nil {
			return fmt.Errorf("failed to create processor: %w", err)
		}
		defer processor.Close()

		// Configure context with cancellation for graceful shutdown
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// Channel to capture OS signals
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

		// Channel for processor errors
		errChan := make(chan error, 1)

		// Start processor in goroutine
		go func() {
			log.Info("Processor started")
			if err := processor.Start(ctx); err != nil {
				errChan <- err
			}
		}()

		// Wait for signal or error
		select {
		case sig := <-sigChan:
			log.Info("Received signal %v, initiating graceful shutdown...", sig)
			cancel()

			// Give time for processor to close gracefully
			shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer shutdownCancel()

			done := make(chan struct{})
			go func() {
				processor.Close()
				close(done)
			}()

			select {
			case <-done:
				log.Info("Service stopped gracefully")
			case <-shutdownCtx.Done():
				log.Warn("Shutdown timeout, forcing exit")
			}

			return nil

		case err := <-errChan:
			cancel()
			return fmt.Errorf("processor error: %w", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
