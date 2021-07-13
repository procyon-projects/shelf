package main

import "testing"

func TestPostgresSqlQueryBuilder_CreateQuery(t *testing.T) {
	queryBuilder := GetSqlQueryBuilder(Postgres)
	_, err := queryBuilder.Table("Users", "u").
		Select("firstName", "lastName").
		Limit(10).
		Offset(0).
		Where().
		Equals("firstName", "test", true).Or().
		Equals("lastName", "").
		OrderBy("firstName").Sort(ASC).
		OrderBy("lastName").Sort(DESC).
		CreateQuery()

	if err != nil {

	}

}
