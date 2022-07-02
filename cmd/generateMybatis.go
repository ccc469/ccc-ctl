/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	internal "github.com/ccc469/ccc-ctl/internal/mybatis"
	"github.com/spf13/cobra"
)

// generateMybatisCmd represents the generateMybatis command
var generateMybatisCmd = &cobra.Command{
	Use:   "generateMybatis",
	Short: "生成mybatis工程文件",
	Long:  `生成mybaits工程文件`,
	Run: func(cmd *cobra.Command, args []string) {
		internal.Run()
	},
}

func init() {
	rootCmd.AddCommand(generateMybatisCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generateMybatisCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generateMybatisCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	generateMybatisCmd.Flags().StringVarP(&internal.BatisType, "batis-type", "c", "tk", "tk, mybatis-plus")
	generateMybatisCmd.Flags().StringVarP(&internal.Table, "table", "t", "", "多个表名英文逗号隔开")
	generateMybatisCmd.Flags().BoolVarP(&internal.IsAllTables, "all-table", "a", false, "是否生成所有表，传值生成所有表")
	generateMybatisCmd.Flags().StringVarP(&internal.Database, "database", "d", "", "Database")
	generateMybatisCmd.Flags().StringVarP(&internal.Host, "host", "H", "127.0.0.1", "Mysql Host")
	generateMybatisCmd.Flags().IntVarP(&internal.Port, "port", "p", 3306, "Mysql Port")
	generateMybatisCmd.Flags().StringVarP(&internal.UserName, "username", "u", "root", "Mysql Username")
	generateMybatisCmd.Flags().StringVarP(&internal.Password, "password", "P", "123456", "Mysql Password")
	generateMybatisCmd.Flags().StringVar(&internal.ModelPackage, "model", "com.example.entity", "实体类包名")
	generateMybatisCmd.Flags().StringVar(&internal.MapperPackage, "mapper", "com.example.mapper", "Mapper接口包名")
	generateMybatisCmd.Flags().StringVar(&internal.XmlPackage, "xml", "com.example.xml", "Xml文件包名")

	generateMybatisCmd.MarkFlagRequired("database")
}
