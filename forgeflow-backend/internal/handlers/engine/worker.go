package engine

import (
	"context"
	"log"
	"time"

	"github.com/rahuldev403/forgeflow/internal/database"
)

var JobQueue = make(chan string, 100)

func StartWorkerPool(numWorkers int) {
	for i := 1; i <= numWorkers; i++ {
		go worker(i)
	}
	log.Printf("Started %d background workers", numWorkers)
}

func worker(id int) {
	log.Printf("worker %d is ready and listening", id)

	for executionID := range JobQueue {
		log.Printf("[Worker %d] Picked up execution: %s", id, executionID)
		time.Sleep(2 * time.Second)
		completeExecution(executionID)
		log.Printf("[Worker %d] Finished execution: %s", id, executionID)
	}
}

func completeExecution(executionID string) {
	query := `UPDATE executions SET status = 'SUCCESS', completed_at = $1 WHERE id = $2`
	_, err := database.DB.Exec(context.Background(), query, time.Now(), executionID)
	if err != nil {
		log.Printf("Failed to update execution %s: %v", executionID, err)
	}
}
