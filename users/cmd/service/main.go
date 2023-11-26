package main

import (
	"fmt"
	"os"

	"github.com/virsavik/alchemist-template/pkg/config"
	"github.com/virsavik/alchemist-template/pkg/system"
	"github.com/virsavik/alchemist-template/users"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("users exitted abnormally: %s\n", err)
		os.Exit(1)
	}
}

func run() (err error) {
	cfg, err := config.ReadConfigFromEnv()
	if err != nil {
		return err
	}

	s, err := system.New(cfg)
	if err != nil {
		return err
	}

	// call the module composition root
	if err = users.Root(s.Waiter().Context(), s); err != nil {
		return err
	}

	fmt.Println("started users service")
	defer fmt.Println("stopped users service")

	s.Waiter().Add(
		s.WaitForWeb,
		//s.WaitForRPC,
		//s.WaitForStream,
	)

	//go func() {
	//	for {
	//		var mem runtime.MemStats
	//		runtime.ReadMemStats(&mem)
	//		s.Logger().Infof("Alloc = %v, TotalAlloc = %v, Sys = %v, NumGC = %v", mem.Alloc/1024, mem.TotalAlloc/1024, mem.Sys/1024, mem.NumGC)
	//		time.Sleep(10 * time.Second)
	//	}
	//}()

	return s.Waiter().Wait()
}
