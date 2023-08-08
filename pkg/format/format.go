package format

import (
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"github.com/ustclug/podzol/pkg/docker"
	"github.com/ustclug/podzol/pkg/utils"
)

func makeTable(w io.Writer) *tablewriter.Table {
	table := tablewriter.NewWriter(w)
	table.SetCenterSeparator("  ")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetTablePadding("  ")
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetNoWhiteSpace(true)
	return table
}

func ShowContainer(w io.Writer, data docker.ContainerInfo) error {
	table := makeTable(w)
	table.AppendBulk([][]string{
		{"Name:", data.Name},
		{"ID:", data.ID},
		{"Port:", strconv.Itoa(int(data.Port))},
		{"Timeout:", data.Deadline.String()},
	})
	table.Render()
	return nil
}

func ListContainers(w io.Writer, data []docker.ContainerInfo) error {
	table := makeTable(w)
	table.SetHeader([]string{"Name", "ID", "Port", "Deadline"})
	for _, c := range data {
		table.Append([]string{
			c.Name,
			c.ID[:12],
			strconv.Itoa(int(c.Port)),
			c.Deadline.String(),
		})
	}
	table.Render()
	return nil
}

var ErrNotWrapped = errors.New("error not wrapped")

func ListContainerActionErrors(w io.Writer, err error) error {
	es := utils.UnwrapErrors(err)
	if es == nil {
		return ErrNotWrapped
	}
	if len(es) > 0 {
		fmt.Fprintf(w, "Errors:\n")
		for _, e := range es {
			fmt.Fprintf(w, "  %v\n", e)
		}
	}
	return nil
}
