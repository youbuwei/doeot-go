package main

import (
	"context"
	"os"

	"github.com/youbuwei/doeot-go/internal/tools/dev"
)

func main() {
	_ = dev.NewCommand().Run(context.Background(), os.Args[1:])
}
