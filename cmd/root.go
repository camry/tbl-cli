package cmd

import (
    "bytes"
    _ "embed"
    "fmt"
    "os"
    "strings"
    "text/template"

    "github.com/camry/g/v2/gerrors/gerror"
    "github.com/samber/lo"
    "github.com/spf13/cobra"
    "github.com/spf13/viper"

    "tbl-cli/internal/gdb"
    "tbl-cli/internal/model"
)

var (
    //go:embed .version
    Version string
    //go:embed template/schema.tpl
    TemplateTable string
)

func Execute() error {
    return rootCmd.Execute()
}

func init() {
    rootCmd.Flags().StringVarP(&server, "server", "s", "", "指定服务器。(格式: <user>:<password>@<host>:<port>)")
    rootCmd.Flags().StringVarP(&db, "db", "d", "", "指定数据库。")
    rootCmd.Flags().StringVarP(&icPath, "config", "c", "", "指定忽略配置文件路径。")
    rootCmd.Flags().StringVarP(&out, "out", "o", "", "指定输出目录。")
    rootCmd.Flags().StringVarP(&pkg, "pkg", "p", "pb", "指定包名。")
    rootCmd.Flags().StringVarP(&goPkg, "go-pkg", "g", "", "指定 Go 包名。")

    cobra.CheckErr(rootCmd.MarkFlagRequired("server"))
    cobra.CheckErr(rootCmd.MarkFlagRequired("db"))
    cobra.CheckErr(rootCmd.MarkFlagRequired("config"))
    cobra.CheckErr(rootCmd.MarkFlagRequired("out"))
}

var (
    server string
    db     string
    icPath string
    out    string
    pkg    string
    goPkg  string

    rootCmd = &cobra.Command{
        Use:     "tbl-cli",
        Short:   "打表代码生成工具。",
        Version: Version,
        Run: func(cmd *cobra.Command, args []string) {
            tblDB, err := gdb.NewTblDB(server, db)
            cobra.CheckErr(err)

            var ic *model.IgnoreConfig
            vi := viper.New()
            vi.SetConfigType("yaml")
            vi.SetConfigFile(icPath)
            err = vi.ReadInConfig()
            cobra.CheckErr(err)
            err = vi.Unmarshal(&ic)
            cobra.CheckErr(err)
            icMap := lo.SliceToMap(ic.Ignores, func(item *model.IgnoreTable) (string, *model.IgnoreTable) {
                return item.Table, item
            })
            icMap["base_json"] = &model.IgnoreTable{
                Table: "base_json",
            }
            icMap["base_json_proto"] = &model.IgnoreTable{
                Table: "base_json_proto",
            }

            pbPkg, ok := lo.Last(strings.Split(goPkg, "/"))
            if !ok {
                cobra.CheckErr(gerror.Newf("--go-pkg 参数错误（github.com/camry/pb/cbpb）。"))
            }

            templateSchema := &model.TemplateSchema{
                Package:   pkg,
                GoPackage: goPkg,
                PbPackage: pbPkg,
            }

            for _, table := range tblDB.GetTables() {
                if _, ok := icMap[table.TableName]; ok {
                    continue
                }
                templateTable := &model.TemplateTable{
                    TableName:           table.TableName,
                    TableNamePascalCase: lo.PascalCase(table.TableName),
                }
                templateSchema.Tables = append(templateSchema.Tables, templateTable)
            }

            tpl, err := template.New("TemplateSchema").Parse(TemplateTable)
            cobra.CheckErr(err)
            var wr bytes.Buffer
            err = tpl.Execute(&wr, templateSchema)
            cobra.CheckErr(err)
            err = os.WriteFile(fmt.Sprintf("%s/schema.go", out), wr.Bytes(), 0644)
            cobra.CheckErr(err)
        },
    }
)
