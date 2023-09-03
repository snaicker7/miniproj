package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"miniproj/entity"
	"miniproj/imageprocess"
	"miniproj/repository"
	"miniproj/repository/dbrepo"
	"net/http"
	"os"
	"strconv"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/segmentio/kafka-go"
)

const (
	host               = "localhost"
	port               = 5432
	user               = "postgres"
	password           = "admin"
	dbname             = "shoppingdb"
	migrationUrl       = "file://db/migration"
	compressfolderName = "compress"
	downloadfolderName = "download"
)

var (
	ProducerChan chan int
	ConsumerChan chan int
	RespChan     chan entity.Product
)

type ImageInterface interface {
	ImageProcessing(product *entity.Product) ([]string, error)
	DownloadImageFile(URL, fileName string) (string, error)
	ImageCompression(fileName string) (string, error)
}

type Application struct {
	DSN                string
	MigrationURL       string
	DB                 repository.DatabaseRepoInterface
	route              *mux.Router
	img                imageprocess.ImageInterface
	compressfolderName string
}

type KafkaConfiguration struct {
	BrokerAdd []string
	Topic     string
	Group     string
}

func main() {
	app := Application{}
	kafkaconfig := KafkaConfiguration{
		BrokerAdd: []string{"localhost:9092"},
		Topic:     "quickstart-events",
		Group:     "my-group",
	}
	app.DSN = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	//set up db connection
	conn, err := app.ConnectToDB()
	if err != nil {
		return
	}
	defer conn.Close()
	ctx, cancel := context.WithCancel(context.Background())
	app.DB = &dbrepo.PostgresDatabaseRepo{DB: conn}
	// go app.kafkaConf.Produce(ctx, *app.chans.ProducerChan)
	// go app.kafkaConf.Consume(ctx, *app.chans.ConsumerChan)
	// go app.getProductURL(ctx)

	go Produce(ctx, kafkaconfig, ProducerChan)
	go Consume(ctx, kafkaconfig, ConsumerChan)
	go app.getProductURL(ctx)
	// err = app.DB.CreateUserDetails()
	// if err != nil {
	// 	log.Println("Error while creating users details")
	// }

	// Create the folder
	cfolderPath := "./" + compressfolderName
	err = createCompressFolder(cfolderPath)
	if err != nil {
		log.Fatal(err)
	}
	app.compressfolderName = compressfolderName

	mux := app.RegisterService()

	//Start a server
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
	defer cancel()
}
func openDB(dsn string) (*sql.DB, error) {

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		errStmt := fmt.Sprintf("Error %s when opening DB\n", err)
		log.Print(errStmt)
		return db, err
	}
	err = db.Ping()
	if err != nil {
		errStmt := fmt.Sprintf("Error %s when pinging DB\n", err)
		log.Print(errStmt)
		return db, err
	}
	return db, nil
}

func (app *Application) ConnectToDB() (*sql.DB, error) {
	conn, err := openDB(app.DSN)
	if err != nil {
		return nil, err
	}
	log.Print("Connected to Postgres")
	return conn, nil
}
func (app *Application) RegisterService() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/addProduct", app.AddProduct).Methods("POST")
	return router
}

func (app *Application) AddProduct(writer http.ResponseWriter, req *http.Request) {
	var user = entity.UserAPIs{}
	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Printf("Error while reading request %s \n", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(reqBody, &user)
	if err != nil {
		log.Printf("Error while unmarshalling the request %s \n", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	// productInfo := fmt.Sprintf("%v", user)
	app.img = &imageprocess.ProductImages{FolderName: app.compressfolderName}
	product, err := app.processData(user)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	resp, err := json.Marshal(product)
	if err != nil {
		log.Printf("Error while unmarshalling the request %s \n", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	// producerChan <- productId
	// <-respChan
	fmt.Println(user)
	writer.WriteHeader(http.StatusOK)
	writer.Write(resp)
}

func (app *Application) processData(user entity.UserAPIs) (entity.Product, error) {
	productId, err := app.DB.InsertProductTable(&user)
	if err != nil {
		return entity.Product{}, err
	}
	ProducerChan <- productId
	product, ok := <-RespChan
	if !ok {
		log.Printf("Error while getting product details from response channel", err)
		return entity.Product{}, errors.New("Error while getting product details from response channel")
	}

	product.CompressedProductImages, err = app.img.ImageProcessing(&product)

	err = app.DB.UpdateProductTable(&product)
	if err != nil {
		return entity.Product{}, err
	}
	productDetail, err := app.DB.FetchProductDetails(product.ProductID)
	if err != nil {
		return entity.Product{}, err
	}
	return *productDetail, nil
}

func createCompressFolder(folderPath string) error {
	_, err := os.Stat(folderPath)
	if os.IsNotExist(err) {
		// Folder does not exist, so create it
		err := os.Mkdir(folderPath, os.ModePerm)
		if err != nil {
			fmt.Printf("Failed to create folder: %v\n", err)
			return err
		}
		log.Printf("Folder '%s' created successfully\n", folderPath)
	} else if err != nil {
		log.Printf("Error checking folder existence: %v\n", err)
		return err
	} else {
		log.Printf("Folder '%s' already exists\n", folderPath)
	}
	return nil
}
func (app *Application) getProductURL(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case prodId := <-ConsumerChan:
			{
				prod, err := app.DB.FetchProductImgUrl(prodId)
				if err != nil {
					log.Printf("error while fetch product image URL")
				}
				RespChan <- *prod
			}
		}
	}
}

func Consume(ctx context.Context, kafkconfig KafkaConfiguration, consumerChan chan int) {
	log.Println("Connecting to Kafka Consumer")
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: kafkconfig.BrokerAdd,
		Topic:   kafkconfig.Topic,
		GroupID: kafkconfig.Group,
	})
	go func() {
		for {
			msg, err := r.ReadMessage(ctx)
			if err != nil {
				log.Panic("could not read message " + err.Error())
			}
			log.Println("received: ", string(msg.Value))
			productId, err := strconv.Atoi(string(msg.Value))
			if err != nil {
				log.Println("Error", err)
			}
			consumerChan <- productId
		}
	}()
}

func Produce(ctx context.Context, kafkconfig KafkaConfiguration, productIdChan chan int) {
	log.Println("Connecting to Kafka Producer")
	i := 0
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: kafkconfig.BrokerAdd,
		Topic:   kafkconfig.Topic,
	})

	go func() {
		for {
			select {
			case chanVal := <-productIdChan:
				err := w.WriteMessages(ctx, kafka.Message{
					Key:   []byte(strconv.Itoa(i)),
					Value: []byte(strconv.Itoa(chanVal)),
				})
				if err != nil {
					log.Panic("could not write message " + err.Error())
				}
				log.Println("writes:", i)
				i++
			}
		}
	}()
}
