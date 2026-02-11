package model

type TemplateSchema struct {
    Package   string
    GoPackage string
    PbPackage string
    Tables    []*TemplateTable
}

type TemplateTable struct {
    TableName           string
    TableNamePascalCase string
}
