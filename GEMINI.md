# Go-Alchemy-Sdk

Golang sdk for alchemy, inspired by https://github.com/alchemyplatform/alchemy-sdk-js.

## Directory

- `/docs`: docs for user.
- `/internal`: external Connection Wrapper
- `/gas`: mainly API of go-alchemy-sdk.
- `/namespace`: API for various functions.
- `/types`: basic types for project.
- `/utils`: utils for project.
- `/ether`: ethereum api
- `/wallet`: wallet api
- `/constant`: constant definition
- `/playground`: playground. this is not for git log.
- `/.github`: github collections

## Rule

- Use TDD for development.
- This project internal use go-ethereum client.
  - If you can use go-ethereum, you should use go-ethereum.
- If you create public API, you should document.
