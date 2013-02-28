package fetcher

import "net/http"
import "io/ioutil"
import "regexp"
import "strings"
import "strconv"

var ResultsChannel = make(chan *Stock)

func ResultsFetcher(nextChan chan *Stock) {
    for {
        stock := <- ResultsChannel
        resultsPage := fetchResultsPage(stock)
        go findResults(stock, resultsPage, nextChan)
    }
}

func fetchResultsPage(stock *Stock) string {
    resp, err := http.Get("http://www.moneycontrol.com/financials/" + stock.URLName + "/results/yearly/" + stock.Name)
    
    if err != nil {
        panic("Could not fetch stock: " + stock.Name)
    }
    defer resp.Body.Close() 
    results_data, err := ioutil.ReadAll(resp.Body)
    
    if err != nil {
        panic("Could not read contents for stock: " + stock.Name)
    }
    return string(results_data)
}

func findResults(stock *Stock, resultsPage string, nextChan chan *Stock) {
    lines := strings.Split (resultsPage, "\n")
    
    eps_line_index := -1

    for line_index, line := range lines{
        if strings.Contains(line, "Earnings Per Share") {
            eps_line_index = line_index
            break
        }
    }
    
    if eps_line_index == -1 {
        panic("Could not scrape Results")
    }
    
    eps_line_index += 2
    re, _ := regexp.Compile(`(\d+\.\d+)`)
    
    for {    
        data_line := lines[eps_line_index]
        match := re.FindStringSubmatch(data_line)
        
        if match == nil { break }
        data, _ := strconv.ParseFloat (match[0], 32)
        stock.EPS = append(stock.EPS, data)
        eps_line_index++
    }
    
    nextChan <- stock
}
