package main

import (
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v4/stdlib"

	"github.com/virsavik/alchemist-template/cmd/banner"
	"github.com/virsavik/alchemist-template/pkg/config"
	"github.com/virsavik/alchemist-template/pkg/system"
	"github.com/virsavik/alchemist-template/users"
)

type monolith struct {
	*system.System
	modules []system.Module
}

func main() {
	banner.Show()

	if err := run(); err != nil {
		fmt.Printf("alchemist-template exitted abnormally: %s\n", err.Error())
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.ReadConfigFromEnv()
	if err != nil {
		return err
	}

	s, err := system.New(cfg)
	if err != nil {
		return err
	}

	m := monolith{
		System: s,
		modules: []system.Module{
			users.Module{},
			// Add more module here
		},
	}

	if err = m.startupModules(); err != nil {
		return err
	}

	fmt.Println("started alchemist-template application")
	defer fmt.Println("stopped alchemist-template application")

	m.Waiter().Add(
		m.WaitForWeb,
	)

	//go func() {
	//	for {
	//		var mem runtime.MemStats
	//		runtime.ReadMemStats(&mem)
	//		m.logger.Msgf("Alloc = %v  TotalAlloc = %v  Sys = %v  NumGC = %v", mem.Alloc/1024, mem.TotalAlloc/1024, mem.Sys/1024, mem.NumGC)
	//		time.Sleep(10 * time.Second)
	//	}
	//}()

	return m.Waiter().Wait()
}

func (m *monolith) startupModules() error {
	for _, module := range m.modules {
		ctx := m.Waiter().Context()
		if err := module.Startup(ctx, m); err != nil {
			return err
		}
	}

	return nil
}
