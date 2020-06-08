# go-coinspaid #

go-coinspaid is a Go client library for accessing the [Coinspaid API](https://docs.coinspaid.com/docs/).

## Installing

### *go get*

```golang
go get -u github.com/oakeshq/go-coinspaid
```

## Example

### Receive cryptocurrency
```golang
import (
  "fmt"

  "github.com/oakeshq/go-coinspaid"
)

func main() {
  client := coinspaid.NewClient("YOUR_API_KEY_HERE", "YOUR_API_SECRET_HERE")

  takeAddressInput := &coinspaid.TakeAddressInput{
    ForeignID: "user-id:2048",
    Currency:  "BTC",
  }

  address, err := client.TakeAddress(takeAddressInput)

  if err != nil {
    fmt.Printf("%s\n", err)
    return
  }

  fmt.Printf("Address: %s\n", address.Address)
}
```