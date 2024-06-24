package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/mbatimel/Data_Showcase_HW/internal/cache"
)
var ErrChannelClosed = errors.New("channel is closed")
type Server interface {
	Run(ctx context.Context) error
	Close() error

}
type server struct {
	srv *http.Server
	cache cache.ICache
}

func (s *server) Run(ctx context.Context) error{
	ch:=make(chan error, 1)
	defer close(ch)
	go func(){
		ch <- s.srv.ListenAndServe()
	}()
	select  {
	case err, ok := <-ch:
		if !ok{
			return ErrChannelClosed
		}
		if err != nil{
			return fmt.Errorf("failed to listen and serve: %w", err)
		}
	case <-ctx.Done():
		if err:=ctx.Err();err!=nil{
			return fmt.Errorf("context faild: %w", err)
		}
			
	}
	return nil
}
func (s *server) Close() error{
	return nil
}
