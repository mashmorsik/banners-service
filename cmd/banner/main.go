package main

import (
	"github.com/mashmorsik/logger"
)

func main() {
	logger.BuildLogger(nil)

	//conf, err := config.LoadConfig()
	//if err != nil {
	//	logger.Errf("Error loading config: %v", err)
	//	return
	//}
	//
	//ctx, cancel := context.WithCancel(context.Background())
	//sigCh := make(chan os.Signal, 1)
	//signal.Notify(sigCh, syscall.SIGINT, syscall.SIGKILL)
	//
	//go func() {
	//	<-sigCh
	//	logger.Infof("context done")
	//	cancel()
	//}()

}
