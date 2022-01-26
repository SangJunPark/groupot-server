package cli

import (
	"flag"
	"fmt"
	"gotutorial/explorer"
	"gotutorial/rest"
	"os"
	"runtime"
)

func Start() {
	//go rest.Start(4000)
	//explorer.Start(5000)
	//fmt.Print("d")
	if len(os.Args) < 2 {
		fmt.Printf("Welcome to 이신영코인")
		fmt.Printf("Blah Blah : ")
		//defer를 이행하고 모든 함수를 제거
		runtime.Goexit()
		// 강제로 종료
		os.Exit(1)
	}

	// //rest command
	// restCmd := flag.NewFlagSet("rest", flag.ExitOnError)
	// portFlag := rest.Int("port", 4000, "port of server default 4000")

	// //js의 slice와 비슷함 [2:5]
	// fmt.Println(os.Args[2:])
	// switch os.Args[1] {
	// case "explorer":
	// 	break
	// case "rest":
	// 	restCmd.Parse(os.Args[2:])
	// 	break
	// default:
	// 	break
	// }
	// if restCmd.Parsed(){
	// 	rest.Start(*portFlag)
	// }else if

	port := flag.Int("port", 4000, "port")
	mode := flag.String("mode", "rest", "resfffadsasdt")
	flag.Parse()
	fmt.Println("ASdf")
	switch *mode {
	case "html":
		fmt.Println("ff")
		explorer.Start(*port)
		break
	case "rest":
		fmt.Println("dd")

		rest.Start(*port)
		break
	default:
		break
	}

	fmt.Println(*port, *mode)
}
