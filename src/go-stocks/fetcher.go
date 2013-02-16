package main
import "net/http"
import "io/ioutil"
import "fmt"
import "regexp"
import "strings"
import "strconv"

func fetchRatiosPage(stock *Stock){
    resp, err := http.Get("http://www.moneycontrol.com/financials/" + stock.URLName + "/ratios/" + stock.Name)
    
    if err != nil {
        panic("Could not fetch stock: " + stock.Name)
    }
    defer resp.Body.Close() 
    ratio_data, err := ioutil.ReadAll(resp.Body)
    stock.RatioPage = string(ratio_data)
}

func findRatio(stock *Stock) {
    lines := strings.Split (stock.RatioPage, "\n")
    
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
}

func main() {
    stock := &Stock {
              Name: "AL",
              URLName: "ashokleyland",
          }
    fetchRatiosPage(stock)
    findRatio(stock)
    fmt.Println(stock.CurrentRatio)
}
