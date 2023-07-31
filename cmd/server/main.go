package main

import (
	"context"
	"fmt"

	cachepkg "github.com/Jamshid90/tcppow/internal/pkg/cache"
	configpkg "github.com/Jamshid90/tcppow/internal/pkg/config"
	serverpkg "github.com/Jamshid90/tcppow/internal/server"
)

func main() {
	// loading config from env
	config, err := configpkg.Load()
	if err != nil {
		fmt.Println("error load config:", err)
		return
	}
	// initialization cache
	cache, err := cachepkg.Init(config.Redis.Host, config.Redis.Port, config.Redis.Password, config.Redis.DB)
	if err != nil {
		fmt.Println("error init cache:", err)
		return
	}
	// initialization handler
	handler := serverpkg.NewHandler(cache, config.HashDuration, config.HashZerosCount)
	// initialization handler
	server := serverpkg.New(config.Server.Host, config.Server.Port, handler)
	// initialization context with cancel
	ctx, cancal := context.WithCancel(context.Background())
	defer cancal()

	// run server
	fmt.Println("listening", fmt.Sprintf("%s:%s", config.Server.Host, config.Server.Port))
	if err = server.Run(ctx); err != nil {
		fmt.Println("server error:", err)
	}
}
