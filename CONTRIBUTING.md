# Contributing

Thank you for contributing

## Development Flow

- Please submit an issue
  TODO: issue template
- Fork repo
- Checkout branch `username/#<issue-number>/<feature>`
  EX: poteto0/#1/first-commit
- Link issue into commit. `<your message> refs: #<issue-number>` to the commit message
- Write UT as much as possible
- PASS the UT (\*you need `-gcflags=all=-l`)
  `go test ./... -cover -gcflags=all=-l`
- PASS the linter
- Crate PR
