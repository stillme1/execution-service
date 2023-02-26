package run

import (
	"os/exec"
)

func Compile(inpath, outpath string) bool {
	return exec.Command("g++", inpath, "-o", outpath).Run() == nil
}
