package services

import (
	"context"
	"time"

	"github.com/padiazg/broken-hexagon/internal/config"
	"github.com/padiazg/broken-hexagon/pkg/logger"
)

// ProcessorConfig holds the processor configuration passed to the factory function
type ProcessorConfig struct {
	Config *config.Config
	Logger logger.Logger
}

// Processor coordinates the main business logic of the service.
// It implements the long-running service pattern with graceful shutdown.
type Processor struct {
	config *config.Config
	Logger logger.Logger
	// TODO: Add your dependencies here
	// Examples:
	//   - Database repositories
	//   - Message queue clients
	//   - External API clients
	//   - File system watchers
}

// NewProcessor creates a new Processor with its dependencies
func NewProcessor(config *ProcessorConfig) (*Processor, error) {
	// TODO: Initialize your dependencies here
	// Example:
	//   repo, err := database.NewRepository(config)
	//   if err != nil {
	//       return nil, fmt.Errorf("failed to initialize repository: %w", err)
	//   }

	p := &Processor{
		config: config.Config,
		Logger: config.Logger,
		// TODO: Set your dependencies here
	}

	return p, nil
}

// Start runs the processor and blocks until context is canceled.
// Implement your main service logic here.
func (p *Processor) Start(ctx context.Context) error {
	p.Logger.Info("Processor starting...")

	// TODO: Implement your service logic here
	//
	// Common patterns:
	//
	// 1. Message Queue Listener:
	//    return p.listenToQueue(ctx)
	//
	// 2. File System Watcher:
	//    return p.watchFiles(ctx)
	//
	// 3. Periodic Task Scheduler:
	//    return p.runScheduledTasks(ctx)
	//
	// 4. Event Stream Processor:
	//    return p.processEventStream(ctx)

	// Example: Simple ticker that runs every 10 seconds
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			p.Logger.Info("Processor shutting down...")
			return nil

		case <-ticker.C:
			if err := p.processTask(ctx); err != nil {
				p.Logger.Error("Task processing error: %v", err)
				// Decide whether to continue or return error
				// return fmt.Errorf("task failed: %w", err)
			}
		}
	}
}

// processTask is an example task that runs periodically
func (p *Processor) processTask(ctx context.Context) error {
	p.Logger.Info("Processing task...")

	// TODO: Implement your task logic here
	// Example:
	//   - Fetch data from source
	//   - Transform/process data
	//   - Store results
	//   - Update metrics

	return nil
}

// Close releases resources and performs cleanup
func (p *Processor) Close() {
	p.Logger.Info("Closing processor...")

	// TODO: Close your resources here
	// Examples:
	//   - Close database connections
	//   - Disconnect from message queues
	//   - Close file handles
	//   - Cancel background goroutines
}
