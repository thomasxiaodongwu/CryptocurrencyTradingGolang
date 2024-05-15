package utils

import "fmt"

func Tqdm(cur_progress int, all int) {
	if all-cur_progress > 1 {
		var percent int = int(float32(cur_progress+1) / float32(all) * 100)
		var rate string
		for i := 0; i < percent; i += 1 {
			rate += "#"
		}
		fmt.Printf("\r[%-100s]%3d%%", rate, percent)
	} else if all-cur_progress == 1 {
		var rate string
		for i := 0; i < 100; i += 1 {
			rate += "#"
		}
		fmt.Printf("\r[%-100s]%3d%%\n", rate, 100)
	}
}
