package main
import(
	"fmt"
	// "errors"
    // "database/sql"
	// "log"
	"context"
	"os"
	"os/signal"
	"syscall"
	"net/http"
	"golang.org/x/sync/errgroup"
    _ "github.com/go-sql-driver/mysql"
)

type App struct {
}

func serveApp() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(resp, "Hello, QCon!")
	})
	return http.ListenAndServe("0.0.0.0:8080", mux)
}

func serveStop() error {
	fmt.Println("server stop interrupt");
	return nil;
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
    group, errCtx := errgroup.WithContext(ctx)

	// done := make(chan error, 1)	
	group.Go(func() error {
		return serveApp()
	})
	group.Go(func() error {
		<-errCtx.Done()
		fmt.Println("Server Error3")
		serveStop()
		cancel();
		return nil;
	})

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)

	group.Go(func() error {
		fmt.Println("Server Error2")
		for {
			select {
				case <-errCtx.Done():
					return errCtx.Err()
				case <-ch:
					cancel()
			}
		}
	})

	if err := group.Wait(); err == nil {
		fmt.Println("Server Error")
	}
}