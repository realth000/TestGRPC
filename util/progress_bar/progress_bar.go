package progress_bar

import (
	"fmt"
)

var (
	posChar = "|/-\\"
	times   = 0
)

func UpdateProgress(taskName string, progress int) {
	fmt.Printf("[%s %3d%%] %s\r", string(getChunk(progress)), progress, taskName)
	if progress >= 100 {
		fmt.Printf("\n")
	}
}

func getChunk(progress int) []byte {
	ret := make([]byte, 10, 10)
	for s := 0; s < len(ret); s++ {
		ret[s] = ' '
	}

	if progress < 10 {
		return insertChunkCursor(ret, 1)
	}

	pos := 1
	for ; pos <= progress/10; pos++ {
		ret[pos-1] = '='
	}
	if progress < 100 {
		return insertChunkCursor(ret, pos)
	}
	return ret
}

func insertChunkCursor(chunk []byte, pos int) []byte {
	if times > 3 {
		times = 0
	}
	chunk[pos-1] = posChar[times]
	times++
	return chunk
}
