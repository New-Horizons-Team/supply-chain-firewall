package printer

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"strings"
)

type DefaultPrinter struct{}

// PrintDiff implements pkg/printer.Printer
func (DefaultPrinter) PrintDiff(ctx context.Context, str string) {
	_ = ctx

	buf := bytes.NewBufferString(str)
	scanner := bufio.NewScanner(buf)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "-") {
			fmt.Printf("\033[0;31m%s\033[0m\n", line)
		} else if strings.HasPrefix(line, "+") {
			fmt.Printf("\033[0;32m%s\033[0m\n", line)
		} else {
			fmt.Println(line)
		}
	}
}

// PrintCriticalReports implements pkg/printer.Printer
func (DefaultPrinter) PrintCriticalReports(ctx context.Context, str string) {
	_ = ctx
	fmt.Println(str)
}

// PrintWarningReports implements pkg/printer.Printer
func (DefaultPrinter) PrintWarningReports(ctx context.Context, str string) {
	_ = ctx
	fmt.Println(str)
}

// Print implements pkg/printer.Printer
func (DefaultPrinter) Print(ctx context.Context, str string) {
	_ = ctx
	fmt.Println(str)
}
