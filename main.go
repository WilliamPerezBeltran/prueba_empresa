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
    
    /*
        Información de SSL y servidores
        https://api.ssllabs.com/api/v3/analyze?host=​<dominio>
        Ejemplo:
        https://api.ssllabs.com/api/v3/analyze?host=truora.com
    */

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

    /*
    el resultado de leer el json se encuentra en el
    struct Result que se pasas como result
    */

    // imprimo el resultado del fetch
    log.Printf("%+v",result)

    //declaro e inicializo mi objectServers vacio

    arrayServer := []ObjectServer{}
    // recorro la ip bajada del Endpoints

    //numero de ips que hay 

    nambersOfIps := len(result.Endpoints)
    fmt.Println("*******len(ip.IpAddress)*******")
    fmt.Println("*******len(ip.IpAddress)*******")
    fmt.Println("*******len(ip.IpAddress)*******")
    fmt.Println("*******len(ip.IpAddress)*******")
    fmt.Println("*******len(ip.IpAddress)*******")
    fmt.Println(nambersOfIps)
    fmt.Println("--------------")
    for _, ip := range result.Endpoints {
        getWhoIs_ip, err := whois.Whois(ip.IpAddress)
        var organizationName string
        var country string
        if err == nil {
            // imprimo las ips de los servidores
            log.Printf("%+v",result)

            arrayOfWhoIs := strings.Split(getWhoIs_ip, "\n")
            for index,data := range arrayOfWhoIs{
                if strings.Contains(data, "OrgName")  {
                        // data[16:] me muestra el nombre de la organización 
                        fmt.Println(index,"=>",data[16:])
                        organizationName = data[16:]

                }

                if strings.Contains(data, "Country")  {
                        // data[16:] me muestra el pais del servidor 
                        fmt.Println(index,"=>",data[16:])
                        country = data[16:]
                }
            }
            
            // declaro e inicializo mi objecto ObjectServer
            // con los valores adecuados 

            serverObject := ObjectServer{
                address:ip.IpAddress, 
                grade:ip.Grade,
                country:country,
                onwer:organizationName,
            }
            // agrego un nuevo objecto de ObjectServer al array 
            arrayServer = append(arrayServer,serverObject)

        }
        
    }
    fmt.Println("******** server server server server server ********")
    log.Printf("%+v",arrayServer)
    fmt.Println("---------------------------------------------")
    fmt.Println("---------------------------------------------")
    fmt.Println("---------------------------------------------")
    fmt.Println("---------------------------------------------")
    fmt.Println("---------------------------------------------")
    fmt.Println("---------------------------------------------")
    // webScrapingTitle(domainGet)
    // bodyOfRequest(domainGet)
    webScrapingLinks(domainGet)


}

func bodyOfRequest(domain string){
    // Make HTTP GET request
     url := "https://www."+domain
    response, err := http.Get(url)
    // response, err := http.Get("https://www.devdungeon.com/")
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
    fmt.Println("**********title***********")
    fmt.Println("**********title***********")
    fmt.Println("**********title***********")
    fmt.Println("**********title***********")
    fmt.Println("**********title***********")
    // Make HTTP GET request
    // response, err := http.Get("https://www.devdungeon.com/")
        url := "https://www."+domain
    response, err := http.Get(url)
    // response, err := http.Get("https://www.truora.com/")
    if err != nil {
        log.Fatal(err)
    }
    defer response.Body.Close()

    // Get the response body as a string
    dataInBytes, err := ioutil.ReadAll(response.Body)
    pageContent := string(dataInBytes)

    // Find a substr
    titleStartIndex := strings.Index(pageContent, "<title>")
    if titleStartIndex == -1 {
        fmt.Println("No title element found")
        // os.Exit(0)
    }
    // The start index of the title is the index of the first
    // character, the < symbol. We don't want to include
    // <title> as part of the final value, so let's offset
    // the index by the number of characers in <title>
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
     // Make HTTP request
     url := "https://www."+domain
    response, err := http.Get(url)
    // response, err := http.Get("https://www.devdungeon.com")
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





func main() {
    router := fasthttprouter.New()
    router.POST("/domain", ReceiveDomainName)

    log.Fatal(fasthttp.ListenAndServe(":8080", router.Handler))
}

