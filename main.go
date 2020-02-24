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
    "./db"
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

}



func main() {
    router := fasthttprouter.New()
    router.POST("/domain", ReceiveDomainName)

    log.Fatal(fasthttp.ListenAndServe(":8080", router.Handler))
}

