package main

import (
	"context"
	"fmt"

	clientpkg "github.com/Jamshid90/tcppow/internal/client"
	configpkg "github.com/Jamshid90/tcppow/internal/pkg/config"
)

func main() {
	// loading config from env
	config, err := configpkg.Load()
	if err != nil {
		fmt.Println("error load config:", err)
		return
	}

	// initialization client
	client := clientpkg.New(config.Server.Host, config.Server.Port)
	// initialization context with cancel
	ctx, cancal := context.WithCancel(context.Background())
	defer cancal()

	// run client
	if err = client.Run(ctx, config.HashMaxIterations); err != nil {
		fmt.Println("server error:", err)
	}
}
