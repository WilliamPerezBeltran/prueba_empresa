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


type CheckRowStruct struct{
    check bool
    check_id int
    check_call int
}


func ReceiveDomainName(ctx *fasthttp.RequestCtx){
    
    //extract data form the domainName key from POST
    domainGet := string(ctx.FormValue("domainName"))
    InsertDomains(domainGet)

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

    number_ips, id_selected, call_selected, check_selected:= check_if_change(nambersOfIps,0,0,false)

    fmt.Println(id_selected)
    fmt.Println(call_selected)
    get_title := string(webScrapingTitle(domainGet))
    webScrapingLinks(domainGet)

    if check_selected == false {
        fmt.Println(" ENTRO A check_selected FALSE OSEA getServerNoChange")
        getServerNoChange(ctx,check_selected,get_title,imgUrl,number_ips)

    }else{
        fmt.Println("********++ ")
        // getServerChange(number_ips,id_selected)

    }
}

type Answer struct{
    Servers []Server
    Server_changed bool `json:"server_changed"`
    Ssl_grade string `json:"ssl_grade"`
    Previous_ssl_grade string `json:"previous_ssl_grade"`
    Logo string `json:"logo"`
    Title string `json:"title"`
    Is_down bool `json:"is_down"`
}

type Server struct{
    Address string `json:"Address"`
    Ssl_grade string `json:"ssl_grade"`
    Country string `json:"country"`
    Owner string `json:"Owner"`
}

func getServerNoChange(ctx *fasthttp.RequestCtx, check_selected bool, get_title string, imgUrl string, number_ips int){
    db, err := sql.Open("postgres", "postgresql://root@localhost:26257/defaultdb?sslmode=disable")
    if err != nil {
        log.Fatal("error connecting to the database: ", err)
    }

    rows, err := db.Query("SELECT address, ssl_grade, country, owner FROM server LIMIT $1",number_ips)
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()
    
    fmt.Println("*****getServerNoChangegetServerNoChange***++ ")

    answer := Answer{}

    counter := 0

    for rows.Next() {
        var  address, ssl_grade, country, owner string
        if err := rows.Scan(&address, &ssl_grade, &country, &owner); err != nil {
            log.Fatal(err)
        }
        item_server := Server{Address:address, Ssl_grade:ssl_grade, Country:country ,Owner:owner}
        answer.Servers = append(answer.Servers,item_server)
        if counter == 0{
            answer.Ssl_grade = ssl_grade

        }else if counter == 1{
            answer.Previous_ssl_grade = ssl_grade

        }
        counter ++
    }

    answer.Server_changed = check_selected
    answer.Logo = imgUrl
    answer.Title = get_title
    answer.Is_down = false
    deleteAllDataFromServer()

    enc := json.NewEncoder(ctx)
    enc.SetIndent("", "  ")
    enc.Encode(&answer)
}

func webScrapingTitle(domain string)[]byte{
    var pageTitle []byte
    url := "https://www."+domain
    response, err := http.Get(url)
    if err != nil {
        // Entra aca cuando el certificado de la pagina no es permitido 
        var str = []string{"No Data: Certificate not permitted"}
        var x = []byte{}
        for i:=0; i<len(str); i++{
            b := []byte(str[i])
            for j:=0; j<len(b); j++{
                x = append(x,b[j])
            }
        }
        return x
        // log.Fatal(err)
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
    pageTitle = []byte(pageContent[titleStartIndex:titleEndIndex])

    // Print out the result
    fmt.Printf("Page title: %s\n", pageTitle)

    return pageTitle
}

func webScrapingLinks(domain string){
     url := "https://www."+domain
    response, err := http.Get(url)
    if err != nil {
        imgUrl = "No Data: Certificate not permitted"
        // log.Fatal(err)
    }else{
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
    imgUrl = href
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

func check_if_change(ips_send int, id_send int, call_send int,check_send bool )(int ,int ,int ,bool){
    db, err := sql.Open("postgres", "postgresql://root@localhost:26257/defaultdb?sslmode=disable")
    if err != nil {
        log.Fatal("error connecting to the database: ", err)
    }

    rows, err := db.Query("SELECT * FROM server LIMIT $1",ips_send)

    if err != nil {
        log.Fatal(err)
    }

    defer rows.Close()

    fmt.Println("Initial balances:")
    
    var id, call int
    var address,ssl_grade,country,owner string
    check:=false
    for rows.Next() {

        if err := rows.Scan(&id, &address, &ssl_grade, &country, &owner, &call); err != nil {
            log.Fatal(err)
        }
        fmt.Printf("%d %s %s %s %s %d\n", id, address,ssl_grade,country,owner,call)

        // get_id, get_address ,get_ssl_grade ,get_country ,get_owner, get_call, get_check := checkRow(id, address ,ssl_grade ,country ,owner, call, check)
        get_id, _ ,_ ,_ ,_, get_call, get_check := checkRow(id, address ,ssl_grade ,country ,owner, call, check_send)

        if get_check == true {
            return ips_send, get_id, get_call,get_check

        }
    }
    return ips_send, id, call, check
}

func checkRow(id int, address string, ssl_grade string, country string, owner string, call int, check bool)(int, string, string, string, string, int, bool ){
    db, err := sql.Open("postgres", "postgresql://root@localhost:26257/defaultdb?sslmode=disable")
    if err != nil {
        log.Fatal("error connecting to the database: ", err)
    }

    rows, err := db.Query("SELECT * FROM server where id != $1 AND call = $2",id,call)

    if err != nil {
        log.Fatal(err)
    }

    defer rows.Close()

    origin_address := address
    origin_ssl_grade := ssl_grade
    origin_country := country
    origin_owner := owner


    for rows.Next() {
        var id, call int
        var address,ssl_grade,country,owner string

        if err := rows.Scan(&id, &address, &ssl_grade, &country, &owner, &call); err != nil {
            log.Fatal(err)
        }
            
        if origin_address != address || origin_ssl_grade != ssl_grade || origin_owner != owner || origin_country != country {
            return id, address,ssl_grade ,country ,owner, call, true
            break
        }

        
    }
return 0,"-","-","-","-",0,false
}

type Items struct{
    Items []Item
}

type Item struct{
    Name string `json:"domain"`
}

func ListAllDomains(ctx *fasthttp.RequestCtx){
    db, err := sql.Open("postgres", "postgresql://root@localhost:26257/defaultdb?sslmode=disable")
    if err != nil {
        log.Fatal("error connecting to the database: ", err)
    }

    rows, err := db.Query("SELECT DISTINCT domain FROM domain_list;")

    if err != nil {
        log.Fatal(err)
    }

    defer rows.Close()

    items := Items{}
    
    for rows.Next() {
        var domain string
        if err := rows.Scan(&domain); err != nil {
            log.Fatal(err)
        }

        item := Item{Name: domain}

        items.Items = append(items.Items,item)
    }
    json.NewEncoder(ctx).Encode(&items)
}

func InsertDomains(domain string){
    db, err := sql.Open("postgres", "postgresql://root@localhost:26257/defaultdb?sslmode=disable")
    if err != nil {
        log.Fatal("error connecting to the database: ", err)
    }

    if _, err := db.Exec("INSERT INTO domain_list (domain) VALUES ($1)",domain)
    err != nil {
        log.Fatal(err)
    }
}

func deleteAllDataFromServer(){
    db, err := sql.Open("postgres", "postgresql://root@localhost:26257/defaultdb?sslmode=disable")
    if err != nil {
        log.Fatal("error connecting to the database: ", err)
    }

    if _, err := db.Exec("DELETE FROM server")
    err != nil {
        log.Fatal(err)
    }
}


var (
 nambersOfIps int   
 imgUrl string   
) 
// var nambersOfIps int 

func main() {
    router := fasthttprouter.New()
    router.POST("/domain", ReceiveDomainName)
    router.GET("/listAllDomains", ListAllDomains)
    log.Fatal(fasthttp.ListenAndServe(":8080", router.Handler))
}

