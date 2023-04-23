package main

import "gorm-demo/query"

func main() {
	// Basic DAO API
	user, err := query.User.Where(query.User.Name.Eq("modi")).First()

	// Dynamic SQL API
	users, err := query.User.FilterWithNameAndRole("modi", "admin")

	query.Use().Transaction()
}
