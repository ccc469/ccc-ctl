package util

import (
	"os"

	"github.com/olekukonko/tablewriter"
)

func Print(data [][]string, headers []string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)
	// table.SetFooter([]string{"", "", "Total", "$146.93"}) // Add Footer
	table.SetBorder(false) // Set Border to false

	table.SetHeaderColor(tablewriter.Colors{tablewriter.Bold, tablewriter.BgGreenColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.BgGreenColor})

	// table.SetColumnColor(tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiBlackColor},
	// 	tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiRedColor})

	// table.SetFooterColor(tablewriter.Colors{}, tablewriter.Colors{},
	// 	tablewriter.Colors{tablewriter.Bold},
	// 	tablewriter.Colors{tablewriter.FgHiRedColor})
	table.SetBorders(tablewriter.Border{Left: true, Top: true, Right: true, Bottom: true})
	table.AppendBulk(data)
	table.SetCenterSeparator("|")
	table.Render()
}
