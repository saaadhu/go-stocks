package fetcher

import "net/http"
import "io/ioutil"
import "regexp"
import "strings"
import "strconv"

var RatioChannel = make(chan *Stock)

func RatioFetcher(nextChan chan *Stock) {
    for {
        stock := <- RatioChannel
        ratiosPage := fetchRatiosPage(stock)
        go findRatio(stock, ratiosPage, nextChan)
    }
}

func fetchRatiosPage(stock *Stock) string {
    resp, err := http.Get("http://www.moneycontrol.com/financials/" + stock.URLName + "/ratios/" + stock.Name)
    
    if err != nil {
        panic("Could not fetch stock: " + stock.Name)
    }
    defer resp.Body.Close() 
    ratio_data, err := ioutil.ReadAll(resp.Body)
    return string(ratio_data)
}

func findRatio(stock *Stock, ratiosPage string, nextChan chan *Stock) {
    lines := strings.Split (ratiosPage, "\n")
    
    current_ratio_line_index := -1

    for line_index, line := range lines{
        if strings.Contains(line, "Current Ratio") {
            current_ratio_line_index = line_index
            break
        }
    }
    
    if current_ratio_line_index == -1 {
        panic("Could not scrape Current Ratio")
    }
    
    current_ratio_line_index += 2
    re, _ := regexp.Compile(`(\d+\.\d+)`)
    
    for {    
        data_line := lines[current_ratio_line_index]
        match := re.FindStringSubmatch(data_line)
        
        if match == nil { break }
        data, _ := strconv.ParseFloat (match[0], 32)
        stock.CurrentRatio = append(stock.CurrentRatio, data)
        current_ratio_line_index++
    }
    
    nextChan <- stock
}
