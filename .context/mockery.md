# Rules here specify how to use mockery

* Mockery v3 is used, see https://vektra.github.io/mockery/latest/ for more details
* It's configured via mockery config [.mockery.yaml](../.mockery.yaml)
* If you need to regenerate mocks, always run `go run github.com/vektra/mockery/v3` from project root without args

## Mocks structure

Standard mockery configuration for all interfaces is defined. It includes the following:
- Uses `testify` compatible mocks with typed expecter methods
- Output mocks in the same dir as interface
- All mocks are written to in a given package are written to `mocks_test.go`
- Rare exceptions are possible.