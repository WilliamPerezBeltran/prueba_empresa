package main

import (
    "fmt"
    "log"
    "github.com/buaazp/fasthttprouter"
    "github.com/valyala/fasthttp"
    "encoding/json"
    "io/ioutil"
    "net/http"
    // "os"
    "github.com/likexian/whois-go"
    "strings"
    "io"
    "os"
    "github.com/PuerkitoBio/goquery"
    _ "github.com/lib/pq"
    "database/sql"
)

type ResultJson struct {
    Host string `json:"host"`
    Endpoints []MyEndpoint `json:"endpoints"`
}

type MyEndpoint struct {
    IpAddress string `json:"ipAddress"`
    Grade string `json:"grade"`
}

type ObjectServer struct{
    address string 
    grade string
    country string
    onwer string 
}


func ReceiveDomainName(ctx *fasthttp.RequestCtx){
    
    //extract data form the domainName key from POST
    domainGet := string(ctx.FormValue("domainName"))

    resp, err := http.Get("https://api.ssllabs.com/api/v3/analyze?host=​"+domainGet)
    // fmt.Println(resp)

    if err != nil {
        log.Fatal(err)

    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)

    if err != nil{
        log.Fatal(err)
    }

    var result ResultJson
    err = json.Unmarshal(body, &result)

    if err != nil {
        log.Fatal(err)
    }

    log.Printf("%+v",result)

    arrayServer := []ObjectServer{}

    nambersOfIps := len(result.Endpoints)
    fmt.Println(nambersOfIps)
    for _, ip := range result.Endpoints {
        getWhoIs_ip, err := whois.Whois(ip.IpAddress)
        var organizationName string
        var country string
        if err == nil {
            // print ips of servers
            log.Printf("%+v",result)

            arrayOfWhoIs := strings.Split(getWhoIs_ip, "\n")
            for index,data := range arrayOfWhoIs{
                if strings.Contains(data, "OrgName")  {
                        // get organizationName
                        fmt.Println(index,"=>",data[16:])
                        organizationName = data[16:]

                }

                if strings.Contains(data, "Country")  {
                        // get country
                        fmt.Println(index,"=>",data[16:])
                        country = data[16:]
                }
            }

            serverObject := ObjectServer{
                address:ip.IpAddress, 
                grade:ip.Grade,
                country:country,
                onwer:organizationName,
            }
            arrayServer = append(arrayServer,serverObject)
        }
        
    }
    fmt.Println("--server--")
    log.Printf("%+v",arrayServer)
    webScrapingTitle(domainGet)
    // bodyOfRequest(domainGet)
    webScrapingLinks(domainGet)

    fmt.Println("**************")
    fmt.Println("**************")
    fmt.Println("**************")
    fmt.Println("**************")
    fmt.Println("**************")
    fmt.Println("**************")
    // conectDb()
    conectDb1()
}

func bodyOfRequest(domain string){
    // Make HTTP GET request
     url := "https://www."+domain
    response, err := http.Get(url)
    if err != nil {
        log.Fatal(err)
    }
    defer response.Body.Close()
    // Copy data from the response to standard output
    n, err := io.Copy(os.Stdout, response.Body)
    if err != nil {
        log.Fatal(err)
    }
    log.Println("Number of bytes copied to STDOUT:", n)
}


func webScrapingTitle(domain string){
    url := "https://www."+domain
    response, err := http.Get(url)
    if err != nil {
        log.Fatal(err)
    }
    defer response.Body.Close()

    dataInBytes, err := ioutil.ReadAll(response.Body)
    pageContent := string(dataInBytes)

    // Find a substr
    titleStartIndex := strings.Index(pageContent, "<title>")
    if titleStartIndex == -1 {
        fmt.Println("No title element found")
        // os.Exit(0)
    }
    // ignore <title>
    titleStartIndex += 7

    // Find the index of the closing tag
    titleEndIndex := strings.Index(pageContent, "</title>")
    if titleEndIndex == -1 {
        fmt.Println("No closing tag for title found.")
        os.Exit(0)
    }

    // (Optional)
    // Copy the substring in to a separate variable so the
    // variables with the full document data can be garbage collected
    pageTitle := []byte(pageContent[titleStartIndex:titleEndIndex])

    // Print out the result
    fmt.Printf("Page title: %s\n", pageTitle)
}


func webScrapingLinks(domain string){
     url := "https://www."+domain
    response, err := http.Get(url)
    if err != nil {
        log.Fatal(err)
    }
    defer response.Body.Close()

    // Create a goquery document from the HTTP response
    document, err := goquery.NewDocumentFromReader(response.Body)
    if err != nil {
        log.Fatal("Error loading HTTP response body. ", err)
    }
    // Find all links and process them with the function
    // defined earlier
    document.Find("link").Each(processElement)
}

func processElement(index int, element *goquery.Selection) {
    // See if the href attribute exists on the element
    href, exists := element.Attr("href")
    if exists {
        // fmt.Println(href)
         if strings.Contains(href,".png"){
            fmt.Println(index,"=>",href)

        }
    }
}

func conectDb(){
    fmt.Println("----------")
    fmt.Println("----------")
    fmt.Println("----------")
    fmt.Println("----------")
    fmt.Println("----------")
    fmt.Println("----------")
    fmt.Println("----------")
    fmt.Println("----------")
    fmt.Println("----------")
    fmt.Println("----------")
    fmt.Println("----------")
    fmt.Println("----------")
    fmt.Println("----------")
    fmt.Println("----------")
    fmt.Println("----------")
    fmt.Println("----------")
    // db, err := sql.Open("postgres", "postgresql://abi@localhost:8081/userdatabase?sslmode=disable")
    // fmt.Println("----------")
    // if err != nil {
    // log.Fatal("error connecting to the database: ", err)
    // }else{
    //     fmt.Println("CONNECT TO THE DATA BASE SUCCESSED")
    // }
    // fmt.Println("------dd----")

    // defer db.Close()

    // stmt, err := db.Prepare("UPDATE tblusers SET name = $1 WHERE id=$2")
    // if err != nil {
    // log.Fatal(err)
    // }

    // defer stmt.Close()

    // res, err := stmt.Exec("marks", 1)
    // if err != nil {
    // log.Fatal(err)
    // }

    // affect, err := res.RowsAffected()
    // if err != nil {
    // log.Fatal(err)
    // }

    // fmt.Println(affect, "rows changed")



}
func conectDb1(){
    // Connect to the "defaultdb" database.
    db, err := sql.Open("postgres", "postgresql://root@localhost:26257/defaultdb?sslmode=disable")
    if err != nil {
        log.Fatal("error connecting to the database: ", err)
    }
    fmt.Println(db)

    if _, err := db.Exec(
        "INSERT INTO domain_list (domain) VALUES ('truora.com'),('facebook.com'),('stackoverflow.com'),('github.com')"); err != nil {
        log.Fatal(err)
    }
    // Print out the balances.
    // rows, err := db.Query("SELECT id, balance FROM accounts")
    rows, err := db.Query("SELECT * FROM domain_list")
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()
    fmt.Println("Initial balances:")
    for rows.Next() {
        var id int
        var domain string
        if err := rows.Scan(&id, &domain); err != nil {
            log.Fatal(err)
        }
        fmt.Printf("%d %d\n", id, domain)
    }

}

func main() {
    router := fasthttprouter.New()
    router.POST("/domain", ReceiveDomainName)
    log.Fatal(fasthttp.ListenAndServe(":8080", router.Handler))
}

