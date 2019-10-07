package main

import (
	"bufio"
	"fmt"
	"github.com/spf13/pflag"
	"os"
)

type selpg_args struct {
	start_page  int
	end_page    int
	in_filename string
	dest        string
	page_len    int
	page_type   int
}

var progname string

const INT_MAX = int(^uint(0) >> 1)

func (arg *selpg_args)process_args() {

	progname = os.Args[0]

	pflag.IntVarP(&arg.start_page, "start_page", "s", 0, "start page")
	pflag.IntVarP(&arg.end_page, "end_page", "e", 0, "end page")
	pflag.IntVarP(&arg.page_type, "page_type", "f", 0, "Page type")
	pflag.IntVarP(&arg.page_len, "page_len", "l", 36, "page len")
	pflag.StringVarP(&arg.dest, "dest", "d", "", "dest")

	pflag.Parse()
	file_arr := pflag.Args()
	arg.in_filename =""

	count := 0
	for{
		if arg.in_filename =="" && file_arr[count] != "<"{
			arg.in_filename = file_arr[count]
		}
		if file_arr[count] == ">"|| file_arr[count] == "|"{
			count++
			arg.dest = file_arr[count]
		}
		count++
		if count >= len(file_arr){
			break
		}
	}

	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "%s: please input enough arguments\n", progname)
		usage()
		os.Exit(1)
	}
	if arg.start_page < 1 || arg.start_page > INT_MAX {
		fmt.Fprintf(os.Stderr, "%s: please input positive integer for start_page\n", progname)
		usage()
		os.Exit(2)
	}
	if arg.end_page < 1 || arg.end_page > (INT_MAX-1) || arg.end_page < arg.start_page {
		fmt.Fprintf(os.Stderr, "%s: please input positive integer for end_page or start_page should not be greater than end_page\n", progname)
		usage()
		os.Exit(3)
	}
	if arg.page_len < 1 || arg.page_len > (INT_MAX-1) {
		fmt.Fprintf(os.Stderr, "%s: please input positive integer for page_len\n", progname)
		pflag.Usage()
		os.Exit(4)
	}
}

func (arg selpg_args)process_input() {

	var read *bufio.Reader
	var write *bufio.Writer

	if arg.in_filename == "" {
		read = bufio.NewReader(os.Stdin)
	} else {
		fin, err := os.Open(arg.in_filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: could not open input file %s\n", progname, arg.in_filename)
			os.Exit(5)
		}
		read = bufio.NewReader(fin)
		defer fin.Close()
	}

	if arg.dest == "" {
		write = bufio.NewWriter(os.Stdout)
	} else {
		fin1, err := os.Open(arg.dest)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: could not open output file %s\n", progname, arg.dest)
			os.Exit(6)
		}
		write = bufio.NewWriter(fin1)
		defer fin1.Close()
	}

	line_number, page_number, pLen := 1, 1, arg.page_len
	judge_Flag := '\n'
	if arg.page_type==1 {
		judge_Flag = '\f'
		pLen = 1
	}

	for {
		line,err:= read.ReadString(byte(judge_Flag));
		if err != nil && len(line) == 0 {
			break
		}
		if line_number > pLen {
			page_number++
			line_number = 1
		}
		if page_number >= arg.start_page && page_number <= arg.end_page {
			write.Write([]byte(line))
		}
		line_number++
	}

	if page_number < arg.end_page {
		fmt.Fprintf(os.Stderr,
			"\n%s: put exist error\n", progname)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "\nUSAGE: %s -sstart_page -eend_page [ -f | -llines_per_page ] [ -ddest ] [ in_filename ]\n", progname)
}

func main() {
	arg := selpg_args{}
	arg.process_args()
	arg.process_input()
}
