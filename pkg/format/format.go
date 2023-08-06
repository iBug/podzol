package format

import (
	"io"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"github.com/ustclug/podzol/pkg/docker"
)

func ShowContainer(w io.Writer, data docker.ContainerInfo) error {
	table := tablewriter.NewWriter(w)
	table.SetCenterSeparator(" ")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetTablePadding(" ")
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetNoWhiteSpace(true)

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
