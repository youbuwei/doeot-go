package main

import (
	"context"
	"os"

	"github.com/youbuwei/doeot-go/internal/tools/bizgen"
)

func main() {
	_ = bizgen.NewCommand().Run(context.Background(), os.Args[1:])
}
