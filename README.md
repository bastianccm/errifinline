# errifinline

Go linter to check if single error assignments are inline.

[Google](https://google.github.io/styleguide/go/guide#concision) | [Uber](https://github.com/uber-go/guide/blob/master/style.md#errors)

## Install

```sh
go install github.com/bastianccm/errifinline/cmd/errifinline@latest
```

## Usage

```sh
errifinline ./...
```

## Example

```go
err := something()
if err != nil {  // this should inline _, err := something()...
    ...
}
```
