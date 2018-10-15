package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	flag "github.com/spf13/pflag"
)

type selpg_args struct {
	start_page  int
	end_page    int
	in_filename string
	page_len    int  // default value, can be overriden by "-l number" on command line
	page_type   bool // flase for lines-delimited, true for form-feed-delimited, default is false
	print_dest  string
}

func main() {
	log := log.New(os.Stderr, "", 0) // error log
	var data string

	// initialize arguments
	var args selpg_args
	flag.IntVarP(&args.start_page, "startpage", "s", -1, "start page")
	flag.IntVarP(&args.end_page, "endpage", "e", -1, "end page")
	flag.IntVarP(&args.page_len, "pagelen", "l", 72, "page length")
	flag.BoolVarP(&args.page_type, "pagetype", "f", false, "page type")
	flag.StringVarP(&args.print_dest, "printdest", "d", "", "print destination")

	flag.Parse()

	// process arguments
	if args.start_page < 1 || args.end_page < 1 || args.start_page > args.end_page || args.page_len < 1 || (args.page_len != 72 && args.page_type) || (flag.NArg() > 1) {
		flag.Usage()
		os.Exit(1)
	}

	// read input
	var in string
	buf := make([]byte, 1024)
	if flag.NArg() != 0 {
		// read from file
		file, open_file_err := os.Open(flag.Args()[0])
		if open_file_err != nil {
			log.Println("Fail to open file.", open_file_err.Error())
		}

		_, read_file_err := file.Read(buf)
		for read_file_err != io.EOF {
			in += string(buf)
			_, read_file_err = file.Read(buf)
		}
	} else {
		// read form std
		reader := bufio.NewReader(os.Stdin)
		_, read_std_err := reader.Read(buf)
		for read_std_err != io.EOF {
			in += string(buf)
			_, read_std_err = reader.Read(buf)
		}
	}

	//process input
	if args.page_type {
		// -f
		pages := strings.SplitAfter(in, "\f")
		if args.end_page > len(pages) {
			log.Println("Page number error.")
			flag.Usage()
			os.Exit(2)
		}
		data = strings.Join(pages[args.start_page-1:args.end_page-1], "")
	} else {
		// -l
		lines := strings.SplitAfter(in, "\n")
		if args.end_page > (len(lines)/args.page_len + 1) {
			log.Println("Page number error.")
			flag.Usage()
			os.Exit(2)
		}
		if len(lines) < args.end_page*args.page_len {
			data = strings.Join(lines[(args.start_page-1)*args.page_len:len(lines)], "")
		} else {
			data = strings.Join(lines[(args.start_page-1)*args.page_len:args.end_page*args.page_len], "")
		}
	}

	writer := bufio.NewWriter(os.Stdout)
	if args.print_dest == "" {
		fmt.Printf("%s", data)
	} else {
		// lp
		cmd := exec.Command("lp", "-d"+args.print_dest)
		lpStdin, err := cmd.StdinPipe()
		if err != nil {
			log.Println("Fail to open lp stdin.", err.Error())
			os.Exit(3)
		}
		go func() {
			defer lpStdin.Close()
			io.WriteString(lpStdin, data)
		}()
		lpStdout, err := cmd.CombinedOutput()
		if err != nil {
			log.Println("Fail to open lp stdout.", err.Error())
			os.Exit(4)
		}

		_, err = writer.Write(lpStdout)
		if err != nil {
			log.Println("Fail to open stdout.", err.Error())
			os.Exit(5)
		}
	}
}
