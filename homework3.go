package main
import(
	"fmt"
	"errors"
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
	serv *http.Server
	ctx context.Context
}

func (a *App) Start() error {
	a.serv = &http.Server{Addr: "0.0.0.0:8080"}
	http.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(resp, "Hello, QCon!")
	})
	return a.serv.ListenAndServe()
}

func (a *App) serveStop() error {
	a.serv.Shutdown(a.ctx);
	return errors.New("server stop interrupt");
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
    group, errCtx := errgroup.WithContext(ctx)

	a := new(App)
	a.ctx = ctx

	group.Go(func() error {
		return a.Start()
	})
	group.Go(func() error {
		<-errCtx.Done()
		a.serveStop()
		cancel();
		return nil;
	})

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)

	group.Go(func() error {
		for {
			select {
				case <-errCtx.Done():
					return errCtx.Err()
				case <-ch:
					cancel()
			}
		}
	})

	if err := group.Wait(); err != nil {
		fmt.Println(err)
	}
}