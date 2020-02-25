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
    // "io"
    "os"
    "github.com/PuerkitoBio/goquery"
    _ "github.com/lib/pq"
    "database/sql"
    "time"
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

    counter_call := 1

    for counter_call <= 10 {
        time.Sleep(1000 * time.Millisecond)
        
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

        arrayServer := []ObjectServer{}

        if len(result.Endpoints) == 0{
            continue
        }

        nambersOfIps = len(result.Endpoints)
        fmt.Println(nambersOfIps)
        for _index, ip := range result.Endpoints {
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
                insertElementsToTableServer(serverObject.address,serverObject.grade,serverObject.country,serverObject.onwer,_index)
            }
            
        }

        counter_call = counter_call + 1
    }

    fmt.Println(nambersOfIps)

    check_if_change(nambersOfIps)


    // fmt.Println("--server--")
    // log.Printf("%+v",arrayServer)
    webScrapingTitle(domainGet)
    webScrapingLinks(domainGet)
    // conectDb1()
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


func conectDb1(){
    db, err := sql.Open("postgres", "postgresql://root@localhost:26257/defaultdb?sslmode=disable")
    if err != nil {
        log.Fatal("error connecting to the database: ", err)
    }

    fmt.Println(db)

    if _, err := db.Exec(
        "INSERT INTO domain_list (domain) VALUES ('truora.com'),('facebook.com'),('stackoverflow.com'),('github.com')"); err != nil {
        log.Fatal(err)
    }

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
        fmt.Printf("%d %s\n", id, domain)
    }
}

func insertElementsToTableServer(address string, sslGrade string, country string, ownwer string, call int){

    if sslGrade == ""{
        sslGrade = "--"
    }

    db, err := sql.Open("postgres", "postgresql://root@localhost:26257/defaultdb?sslmode=disable")
    if err != nil {
        log.Fatal("error connecting to the database: ", err)
    }

    if _, err := db.Exec("INSERT INTO server (address,ssl_grade,country,owner,call) VALUES($1, $2, $3, $4, $5)",address,sslGrade,country,ownwer,call)
    err != nil {
        log.Fatal(err)
    }
}

func check_if_change(ips int){

    db, err := sql.Open("postgres", "postgresql://root@localhost:26257/defaultdb?sslmode=disable")
    if err != nil {
        log.Fatal("error connecting to the database: ", err)
    }

    // if _, err := db.Exec("SELECT * FROM server LIMIT $1",ips)
    // err != nil {
    //     log.Fatal(err)
    // }

    rows, err := db.Query("SELECT * FROM server LIMIT $1",ips)

    if err != nil {
        log.Fatal(err)
    }

    defer rows.Close()

    fmt.Println("Initial balances:")
    fmt.Println("Initial balances:")
    fmt.Println("Initial balances:")
    fmt.Println("Initial balances:")
    fmt.Println("Initial balances:")
    fmt.Println("Initial balances:")
    fmt.Println("Initial balances:")
    
    for rows.Next() {
        var id, call int
        var address,ssl_grade,country,owner string
        if err := rows.Scan(&id, &address, &ssl_grade, &country, &owner, &call); err != nil {
            log.Fatal(err)
        }
        fmt.Printf("%d %s %s %s %s %d\n", id, address,ssl_grade,country,owner,call)
    }


    // SELECT (id,address,ssl_grade,country,owner,call) FROM server where call = 0;
    // SELECT * FROM server LIMIT 1;

    // db, err := sql.Open("postgres", "postgresql://root@localhost:26257/defaultdb?sslmode=disable")
    // if err != nil {
    //     log.Fatal("error connecting to the database: ", err)
    // }

    // if _, err := db.Exec("SELECT (address,ssl_grade,country,owner) FROM server where call = 0",address,sslGrade,country,ownwer,call)
    // err != nil {
    //     log.Fatal(err)
    // }

    // for i := 0; i<= (len(ips)-1);i++{

    // }

    // if _, err := db.Exec("INSERT INTO server (address,ssl_grade,country,owner,call) VALUES($1, $2, $3, $4, $5)",address,sslGrade,country,ownwer,call)
    


}


var nambersOfIps int 

func main() {
    router := fasthttprouter.New()
    router.POST("/domain", ReceiveDomainName)
    log.Fatal(fasthttp.ListenAndServe(":8080", router.Handler))
}

