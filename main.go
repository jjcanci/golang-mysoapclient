package main

import (
	"bytes"
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

// Envelope Generado con https://www.onlinetool.io/xmltogo/
type Envelope struct {
	XMLName xml.Name `xml:"Envelope"`
	Text    string   `xml:",chardata"`
	Soap    string   `xml:"soap,attr"`
	Xsi     string   `xml:"xsi,attr"`
	Xsd     string   `xml:"xsd,attr"`
	Body    struct {
		Text                string `xml:",chardata"`
		GetLocationResponse struct {
			Text              string `xml:",chardata"`
			Xmlns             string `xml:"xmlns,attr"`
			GetLocationResult string `xml:"GetLocationResult"`
		} `xml:"GetLocationResponse"`
	} `xml:"Body"`
}

// GeoIP estructura interna del response
type GeoIP struct {
	XMLName xml.Name `xml:"GeoIP"`
	Text    string   `xml:",chardata"`
	Country string   `xml:"Country"`
	State   string   `xml:"State"`
}

type GetCountryISO2ByNameResponse struct {
	XMLName xml.Name `xml:"Envelope"`
	Text    string   `xml:",chardata"`
	Soap    string   `xml:"soap,attr"`
	Xsi     string   `xml:"xsi,attr"`
	Xsd     string   `xml:"xsd,attr"`
	Body    struct {
		Text                         string `xml:",chardata"`
		GetCountryISO2ByNameResponse struct {
			Text                       string `xml:",chardata"`
			Xmlns                      string `xml:"xmlns,attr"`
			GetCountryISO2ByNameResult string `xml:"GetCountryISO2ByNameResult"`
		} `xml:"GetCountryISO2ByNameResponse"`
	} `xml:"Body"`
}

func main() {
	getLocation()
	log.Println("////")
	getCountryISO2ByName("Argentina")
}

func getLocation() {
	log.Println("============= Inicio GetLocation ================")

	// wsdl service url
	url := fmt.Sprintf("%s%s",
		"http://wsgeoip.lavasoft.com",
		"/ipservice.asmx",
	)
	log.Println("URL:", url)

	// payload
	payload := []byte(strings.TrimSpace(`
		<Envelope xmlns="http://www.w3.org/2003/05/soap-envelope">
    		<Body>
        		<GetLocation xmlns="http://lavasoft.com/"/>
    		</Body>
		</Envelope>		
		`,
	))

	httpMethod := "POST"

	// soap action
	//soapAction := "urn:GetLocation"

	log.Println("-> Preparing the request")

	// prepare the request
	req, err := http.NewRequest(httpMethod, url, bytes.NewReader(payload))
	if err != nil {
		log.Fatal("Error on creating request object. ", err.Error())
		return
	}

	// set the content type header, as well as the oter required headers
	req.Header.Set("Content-type", "text/xml")

	// req.Header.Set("SOAPAction", soapAction)
	/*req.Header.Set("Authorization", fmt.Sprintf(
		"Basic %s",
		base64.StdEncoding.EncodeToString([]byte(
			username+":"+password,
		)),
	))*/

	// prepare the client request
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
			Dial: (&net.Dialer{
				Timeout: 5 * time.Second,
			}).Dial,
			TLSHandshakeTimeout: 5 * time.Second,
		},
		Timeout: time.Second * 10,
	}

	log.Println("-> Dispatching the request")

	// dispatch the request
	res, err := client.Do(req)
	if err != nil {
		log.Fatal("Error on dispatching request. ", err.Error())
		return
	}
	defer res.Body.Close()

	log.Println("-> Retrieving and parsing the response")

	log.Println(res.Status)
	log.Println("-> Everything is good, printing users data")

	// read and parse the response body
	result := new(Envelope)
	err = xml.NewDecoder(res.Body).Decode(result)
	if err != nil {
		log.Fatal("Error on unmarshaling xml. ", err.Error())
		return
	}
	log.Println("Country:", result.Body.GetLocationResponse.GetLocationResult)

	geoIP := new(GeoIP)
	err = xml.Unmarshal([]byte(result.Body.GetLocationResponse.GetLocationResult), &geoIP)
	if err != nil {
		log.Fatal("Error on unmarshaling xml. ", err.Error())
		return
	}
	log.Println("geo.Country:", geoIP.Country)

	log.Println("============= Fin GetLocation ================")
}

func getCountryISO2ByName(pais string) {
	log.Println("============= Inicio getCountryISO2ByName ================")

	// wsdl service url
	url := fmt.Sprintf("%s%s",
		"http://wsgeoip.lavasoft.com",
		"/ipservice.asmx",
	)
	log.Println("URL:", url)

	type QueryData struct {
		Pais string
	}

	// Template payload
	const payloadTemplate = `
		<Envelope xmlns="http://www.w3.org/2003/05/soap-envelope">
			<Body>
				<GetCountryISO2ByName xmlns="http://lavasoft.com/">
					<countryName>{{.Pais}}</countryName>
				</GetCountryISO2ByName>
			</Body>
		</Envelope>
		`

	// Create a new template and parse
	t := template.Must(template.New("payloadTemplate").Parse(payloadTemplate))

	// Completar con querydata
	querydata := QueryData{Pais: "Peru"}
	var doc bytes.Buffer
	err := t.Execute(&doc, querydata)
	if err != nil {
		log.Panic(err)
	}

	// payload
	payload := []byte(strings.TrimSpace(doc.String()))

	httpMethod := "POST"

	log.Println("-> Preparing the request")

	// prepare the request
	req, err := http.NewRequest(httpMethod, url, bytes.NewReader(payload))
	if err != nil {
		log.Fatal("Error on creating request object. ", err.Error())
		return
	}

	// set the content type header, as well as the oter required headers
	req.Header.Set("Content-type", "text/xml")

	// req.Header.Set("SOAPAction", soapAction)
	/*req.Header.Set("Authorization", fmt.Sprintf(
		"Basic %s",
		base64.StdEncoding.EncodeToString([]byte(
			username+":"+password,
		)),
	))*/

	// prepare the client request
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
			Dial: (&net.Dialer{
				Timeout: 5 * time.Second,
			}).Dial,
			TLSHandshakeTimeout: 5 * time.Second,
		},
		Timeout: time.Second * 10,
	}

	log.Println("-> Dispatching the request")

	// dispatch the request
	res, err := client.Do(req)
	if err != nil {
		log.Fatal("Error on dispatching request. ", err.Error())
		return
	}
	defer res.Body.Close()

	log.Println("-> Retrieving and parsing the response")

	log.Println(res.Status)
	log.Println("-> Everything is good, printing users data")

	// read and parse the response body
	result := new(GetCountryISO2ByNameResponse)
	err = xml.NewDecoder(res.Body).Decode(result)
	if err != nil {
		log.Fatal("Error on unmarshaling xml. ", err.Error())
		return
	}
	log.Println("Country:", result.Body.GetCountryISO2ByNameResponse.GetCountryISO2ByNameResult)

	geoIP := new(GeoIP)
	err = xml.Unmarshal([]byte(result.Body.GetCountryISO2ByNameResponse.GetCountryISO2ByNameResult), &geoIP)
	if err != nil {
		log.Fatal("Error on unmarshaling xml. ", err.Error())
		return
	}
	log.Println("geo.Country:", geoIP.Country)

	log.Println("============= Fin getCountryISO2ByName ================")
}
