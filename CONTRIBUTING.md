# Contributing

Thank you for contributing

## Building Environment

If you want to try all uts in local, you need to create AlchemyAccount and `.env`

```bash
cp .example.env .env
```

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
