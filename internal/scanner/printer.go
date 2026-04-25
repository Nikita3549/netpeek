package scanner

import (
	"fmt"
	"os"
	"strings"
)

type Printer struct {
	total int
	done  int
}

func NewPrinter(total uint16) *Printer {
	return &Printer{
		total: int(total),
	}
}

func (p *Printer) OutputInitialText(host string, ports string) {
	fmt.Fprintf(os.Stdout, "\033[90mtarget\033[0m \033[34m%s\033[0m \033[90mports\033[0m \033[34m%s\033[0m\n\n", host, ports)
}

func (p *Printer) OutputPort(port uint16, open bool) {
	p.done++

	pct := p.done * 100 / p.total
	filled := pct * 40 / 100
	p.renderOutput(pct, filled)

	if !open {
		return
	}

	fmt.Fprintf(os.Stderr, "\r\033[K")
	fmt.Fprintf(os.Stdout, "\033[32mOPEN\033[0m  %d/tcp   \033[36m%s\033[0m\n", port, "")
}

func (p *Printer) OutputFinalStats(opened, closed uint16, elapsed float64) {
	fmt.Fprintf(os.Stderr, "\r\033[K")

	fmt.Fprintf(os.Stdout, "\n")
	fmt.Fprintf(os.Stdout, "\033[90m%-6s  %-8s  %-6s\033[0m\n",
		"open", "closed", "elapsed")
	fmt.Fprintf(os.Stdout, "\033[32m%-6d\033[0m  \033[31m%-8d\033[0m  \033[34m%.2fs\033[0m\n",
		opened, closed, elapsed)
}

func (p *Printer) renderOutput(pct, filled int) {
	bar := strings.Repeat("█", filled) + strings.Repeat("░", 40-filled)
	fmt.Fprintf(os.Stderr, "\r%3d%%  %s  (%d/%d)", pct, bar, p.done, p.total)
}
