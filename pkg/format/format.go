package format

import (
	"io"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"github.com/ustclug/podzol/pkg/docker"
)

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
			c.ID,
			strconv.Itoa(int(c.Port)),
			c.Deadline.String(),
		})
	}
	table.Render()
	return nil
}
