import (
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// ...

	// Create a context to handle graceful shutdown.
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Wait for a shutdown signal or a key press.
	fmt.Println("Press Ctrl+C or send SIGINT/SIGTERM to stop the server.")
	select {
	case <-ctx.Done():
		fmt.Println("Shutting down the server gracefully...")
	case <-keyPress():
		fmt.Println("Key press detected, shutting down the server gracefully...")
		cancel() // Manually cancel the context to trigger shutdown.
	}

	// Wait for the server to shut down gracefully.
	select {
	case <-time.After(5 * time.Second):
		fmt.Println("Server did not shut down gracefully. Exiting...")
	case <-serverDone:
		fmt.Println("Server has been shut down.")
	}
}

// keyPress waits for a key press and returns a channel that signals when a key is pressed.
func keyPress() <-chan struct{} {
	done := make(chan struct{})
	go func() {
		var buf [1]byte
		os.Stdin.Read(buf[:])
		done <- struct{}{}
	}()
	return done
}
