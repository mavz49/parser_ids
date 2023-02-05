package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/mavz49/parser_ids"
)

var dirPath string
var serverHost string
var serverHostMetadata string
var Tocken string
var BulkCount int
var ParserJson string

var SaveCount int

type SearchMetadata struct {
	Id   string   `json:"id"`
	Data []string `json:"data"`
}

type ResponseMetadata struct {
	Id   string            `json:"id"`
	Data map[string]string `json:"data"`
}

const (
	defaultName = "world"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
	name = flag.String("name", defaultName, "Name to greet")
)

func main() {
	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}

	dirPath = viper.GetString("dirPath")
	serverHost = viper.GetString("serverHost")
	serverHostMetadata = viper.GetString("serverHostMetadata")
	Tocken = viper.GetString("tocken")
	BulkCount = viper.GetInt("bulkCount")
	ParserJson = viper.GetString("parserJson")

	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayPing(ctx, &pb.PingRequest{Name: *name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())
	r, err = c.SayPingAgain(ctx, &pb.PingRequest{Name: *name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())

	/*
	   	Ids, err := GetIdsJsonFile(ParserJson)
	   	if err != nil {
	   		fmt.Println(err)
	   	}

	   	//TdsTest := []string{"99993"}
	   	for _, id := range Ids {

	   		searchMetadata := SearchMetadata{
	   			Id: id,
	   			Data: []string{
	   				"cataltname",
	   				"brand",
	   				"packid",
	   				"cost",
	   				"profit",
	   				"packname",
	   				"packarticle",
	   				"indicate",
	   				"piecesinpack",
	   				"sizetype",
	   				"sales_notes",
	   				"alias",
	   				"cat_name",
	   				"costdelivery",
	   				"sid",
	   				"name",
	   				"color",
	   				"collection",
	   			},
	   		}

	   		_, err := getOneMetadata(searchMetadata)
	   		if err != nil {
	   			log.Println("err ", err)
	   		}
	   		SaveCount = SaveCount + 1
	   		var c int = SaveCount % 1000
	   		if c == 0 {
	   			log.Println("получено ID: ", SaveCount)
	   		}

	   		//log.Println("res ", res)

	   /*
	*/
}

func getOneMetadata(sm SearchMetadata) (*ResponseMetadata, error) {
	messagesJson := sm
	client := &http.Client{}

	bytesRepresentation, err := json.Marshal(messagesJson)
	if err != nil {
		return nil, err
	}

	bodyJson := bytes.NewReader(bytesRepresentation)
	//log.Println("json для отправки ", string(bytesRepresentation))
	//log.Println("адресс ", serverHostMetadata+"/api/metadata/get")
	req, err := http.NewRequest("POST", serverHostMetadata+"/api/metadata/get", bodyJson)
	req.Header.Add("User-Agent", `Mozilla/5.0 Gecko/20100101 Firefox/39.0`)
	req.Header.Add("Authorization", "Bearer "+Tocken)

	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	//log.Println(resp)
	var resultResponse ResponseMetadata
	json.NewDecoder(resp.Body).Decode(&resultResponse)

	return &resultResponse, nil
}

func GetIdsJsonFile(jsonFile string) (map[string]string, error) {
	var data []string
	file, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(file, &data)
	if err != nil {
		return nil, err
	}
	//log.Println("StructureJsonFile объект: ", data)
	ids := map[string]string{}
	for _, v := range data {
		ids[v] = v
	}
	return ids, nil
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
