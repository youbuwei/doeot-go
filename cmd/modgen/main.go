package main

import (
	"context"
	"os"

	"github.com/youbuwei/doeot-go/internal/tools/modgen"
)

func main() {
	_ = modgen.NewCommand().Run(context.Background(), os.Args[1:])
}
