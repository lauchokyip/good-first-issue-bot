package types

type IssueQueryWithNumber struct {
	IssueQuery
	Number int
}

type IssueQuery struct {
	Owner string
	Repo  string
}
