package main

import "fetcher"
import "fmt"

func createStocks() (stocks []*fetcher.Stock) {
    stocks = append(stocks, &fetcher.Stock {
              Name: "AL",
              URLName: "ashokleyland",
          })
    return
}

var outputChannel = make(chan *fetcher.Stock)

func main() {
    stocks := createStocks()
    stockCount := len(stocks)
    
    // Setup pipeline
   go fetcher.RatioFetcher (fetcher.ResultsChannel)
   go fetcher.ResultsFetcher (outputChannel)
    
    // Feed the pipe
    for _,stock := range stocks {
        fetcher.RatioChannel <- stock  
    }
    
    analyzedStocks := 0
    for {
        stock := <- outputChannel 
        fmt.Println(stock.Name, stock.URLName, stock.CurrentRatio, stock.EPS)
        if analyzedStocks++; analyzedStocks == stockCount {
            break
        }
    }
}
