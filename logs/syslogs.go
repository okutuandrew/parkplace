package logs

import "log"
import "os"
import "io"

func SysLogs() {
	file, err := os.OpenFile("logs/syslogs", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Println("FILE MISSING:", err)
		return
	}

	multi := io.MultiWriter(os.Stdout, file)
	log.SetOutput(multi)

	log.Println("TEST")
}
