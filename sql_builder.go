package main

import (
	"errors"
	"strconv"
	"strings"
)

type Sort int

const (
	ASC Sort = iota
	DESC
)

const (
	Postgres = "Postgres"
)

type Table struct {
	Name  string
	Alias string
}

type Query struct {
	Text string
}

type SqlMultiConditions interface {
	Or() SqlConditions
	And() SqlConditions
	OrderBy(column string) SqlSort
}

type SqlOrder interface {
	OrderBy(column string) SqlSort
	CreateQuery() (Query, error)
}

type SqlSort interface {
	Sort(sort Sort) SqlOrder
	CreateQuery() (Query, error)
}

type SqlConditions interface {
	Equals(column string, value string, ignoreCase ...bool) SqlMultiConditions
	Not(column string, value string, ignoreCase ...bool) SqlMultiConditions
	GreaterThan(column string, value string) SqlMultiConditions
	GreaterThanOrEqual(column string, value string) SqlMultiConditions
	LessThan(column string, value string) SqlMultiConditions
	LessThanOrEqual(column string, value string) SqlMultiConditions
	Between(column string, value1 string, value2 string) SqlMultiConditions
	IsNull(column string) SqlMultiConditions
	Null(column string) SqlMultiConditions
	IsNotNull(column string) SqlMultiConditions
	NotNull(column string) SqlMultiConditions
	In(column string, values ...string) SqlMultiConditions
	NotIn(column string, values ...string) SqlMultiConditions
	True(column string) SqlMultiConditions
	False(column string) SqlMultiConditions
	Like(column string, value string) SqlMultiConditions
	StartWith(column string, value string) SqlMultiConditions
	EndWith(column string, value string) SqlMultiConditions
	GroupConditions() SqlConditions
	OrderBy(column string) SqlSort
}

type SqlMultipleJoins interface {
	Join(table string, tableAlias ...string) SqlJoin
	OrderBy(column string) SqlSort
	Where() SqlConditions
	CreateQuery() (Query, error)
}

type SqlJoin interface {
	InnerJoin(otherTable string, tableKey string, otherTableKey string) SqlMultipleJoins
	LeftJoin(otherTable string, tableKey string, otherTableKey string) SqlMultipleJoins
	RightJoin(otherTable string, tableKey string, otherTableKey string) SqlMultipleJoins
	FullJoin(otherTable string, tableKey string, otherTableKey string) SqlMultipleJoins
}

type SqlSelect interface {
	Select(columns ...string) SqlSelect
	Limit(limit uint) SqlSelect
	Join(table string, tableAlias ...string) SqlJoin
	Offset(offset uint) SqlSelect
	OrderBy(column string) SqlSort
	Where() SqlConditions
	CreateQuery() (Query, error)
}

type SqlQueryBuilder interface {
	Table(name string, alias ...string) SqlSelect
}

func GetSqlQueryBuilder(database string) SqlQueryBuilder {

	if database == Postgres {
		return &postgresSqlQueryBuilder{}
	}

	return nil
}

type postgresSqlQueryBuilder struct {
	table         *Table
	selectColumns []string
	useLimit      bool
	useOffset     bool
	limit         uint
	offset        uint
	orderColumn   string
	orderSort     Sort
	orders        []string
}

func (builder *postgresSqlQueryBuilder) Table(name string, alias ...string) SqlSelect {
	aliasName := ""

	if len(alias) > 0 {
		aliasName = alias[0]
	}

	builder.table = &Table{
		Name:  name,
		Alias: aliasName,
	}
	return builder
}

func (builder *postgresSqlQueryBuilder) Select(columns ...string) SqlSelect {
	builder.selectColumns = columns
	return builder
}

func (builder *postgresSqlQueryBuilder) Limit(limit uint) SqlSelect {
	builder.useLimit = true
	builder.limit = limit
	return builder
}

func (builder *postgresSqlQueryBuilder) Offset(offset uint) SqlSelect {
	builder.useOffset = true
	builder.offset = offset
	return builder
}

func (builder *postgresSqlQueryBuilder) Sort(sort Sort) SqlOrder {
	builder.orderSort = sort
	return builder
}

func (builder *postgresSqlQueryBuilder) Where() SqlConditions {
	return builder
}

func (builder *postgresSqlQueryBuilder) Equals(column string, value string, ignoreCase ...bool) SqlMultiConditions {
	return builder
}

func (builder *postgresSqlQueryBuilder) Not(column string, value string, ignoreCase ...bool) SqlMultiConditions {
	return builder
}

func (builder *postgresSqlQueryBuilder) GreaterThan(column string, value string) SqlMultiConditions {
	return builder
}

func (builder *postgresSqlQueryBuilder) GreaterThanOrEqual(column string, value string) SqlMultiConditions {
	return builder
}

func (builder *postgresSqlQueryBuilder) LessThan(column string, value string) SqlMultiConditions {
	return builder
}

func (builder *postgresSqlQueryBuilder) LessThanOrEqual(column string, value string) SqlMultiConditions {
	return builder
}

func (builder *postgresSqlQueryBuilder) Between(column string, value1 string, value2 string) SqlMultiConditions {
	return builder
}

func (builder *postgresSqlQueryBuilder) IsNull(column string) SqlMultiConditions {
	return builder
}

func (builder *postgresSqlQueryBuilder) Null(column string) SqlMultiConditions {
	return builder
}

func (builder *postgresSqlQueryBuilder) IsNotNull(column string) SqlMultiConditions {
	return builder
}

func (builder *postgresSqlQueryBuilder) NotNull(column string) SqlMultiConditions {
	return builder
}

func (builder *postgresSqlQueryBuilder) In(column string, values ...string) SqlMultiConditions {
	return builder
}

func (builder *postgresSqlQueryBuilder) NotIn(column string, values ...string) SqlMultiConditions {
	return builder
}

func (builder *postgresSqlQueryBuilder) True(column string) SqlMultiConditions {
	return builder
}

func (builder *postgresSqlQueryBuilder) False(column string) SqlMultiConditions {
	return builder
}

func (builder *postgresSqlQueryBuilder) Like(column string, value string) SqlMultiConditions {
	return builder
}

func (builder *postgresSqlQueryBuilder) StartWith(column string, value string) SqlMultiConditions {
	return builder
}

func (builder *postgresSqlQueryBuilder) EndWith(column string, value string) SqlMultiConditions {
	return builder
}

func (builder *postgresSqlQueryBuilder) Or() SqlConditions {
	return builder
}

func (builder *postgresSqlQueryBuilder) And() SqlConditions {
	return builder
}

func (builder *postgresSqlQueryBuilder) GroupConditions() SqlConditions {
	return builder
}

func (builder *postgresSqlQueryBuilder) Join(table string, tableAlias ...string) SqlJoin {
	return builder
}

func (builder *postgresSqlQueryBuilder) InnerJoin(otherTable string, tableKey string, otherTableKey string) SqlMultipleJoins {
	return builder
}

func (builder *postgresSqlQueryBuilder) LeftJoin(otherTable string, tableKey string, otherTableKey string) SqlMultipleJoins {
	return builder
}

func (builder *postgresSqlQueryBuilder) RightJoin(otherTable string, tableKey string, otherTableKey string) SqlMultipleJoins {
	return builder
}

func (builder *postgresSqlQueryBuilder) FullJoin(otherTable string, tableKey string, otherTableKey string) SqlMultipleJoins {
	return builder
}

func (builder *postgresSqlQueryBuilder) OrderBy(column string) SqlSort {
	if builder.orderColumn != "" {
		order := builder.orderColumn

		if builder.orderSort == ASC {
			order = order + " ASC"
		} else {
			order = order + " DESC"
		}

		builder.orders = append(builder.orders, order)

		builder.orderColumn = ""
		builder.orderSort = ASC
	}

	builder.orderColumn = column
	return builder
}

func (builder *postgresSqlQueryBuilder) CreateQuery() (Query, error) {
	query := "SELECT"

	if len(builder.selectColumns) == 0 {
		query = query + " * "
	} else {
		query = query + " " + strings.Join(builder.selectColumns, ", ") + " "
		builder.selectColumns = []string{}
	}

	if builder.table == nil {
		return Query{}, errors.New("table name cannot be empty")
	}

	query = query + "FROM " + builder.table.Name

	if builder.table.Alias != "" {
		query = query + " AS " + builder.table.Alias + " "
	}

	if builder.orderColumn != "" {
		order := builder.orderColumn

		if builder.orderSort == ASC {
			order = order + " ASC"
		} else {
			order = order + " DESC"
		}

		builder.orders = append(builder.orders, order)

		builder.orderColumn = ""
		builder.orderSort = ASC
	}

	if len(builder.orders) > 0 {
		query = query + "ORDER BY "
		query = query + strings.Join(builder.orders, " ")
		builder.orders = []string{}
	}

	if builder.useLimit {
		query = query + " LIMIT " + strconv.Itoa(int(builder.limit))
	}

	if builder.useOffset {
		query = query + " OFFSET " + strconv.Itoa(int(builder.limit))
	}

	builder.useOffset = false
	builder.useLimit = false

	return Query{}, nil
}
