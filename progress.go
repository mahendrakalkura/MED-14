package main

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"os"
)

func progress(settings *Settings) {
	fmt.Println("progress()")

	database := get_database(settings)

	var items = [][]string{}

	total, completed, pending, percentage := addresses_progress(database)
	items = append(items, []string{total, completed, pending, percentage})

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetAutoFormatHeaders(false)
	table.SetAutoWrapText(false)
	table.SetHeader([]string{"Total", "Completed", "Pending", "Percentage"})
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.AppendBulk(items)
	table.Render()
}
