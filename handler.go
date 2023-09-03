package main

// func (app *Application) AddProduct(writer http.ResponseWriter, req *http.Request) {
// 	var user = entity.UserAPIs{}
// 	reqBody, err := ioutil.ReadAll(req.Body)
// 	if err != nil {
// 		log.Printf("Error while reading request %s \n", err)
// 		return
// 	}
// 	err = json.Unmarshal(reqBody, &user)
// 	if err != nil {
// 		log.Printf("Error while unmarshalling the request %s \n", err)
// 		return
// 	}
// 	productInfo := fmt.Sprintf("%v", user)
// 	app.user = &user
// 	// productId := InsertProductTable(db, user)
// 	// producerChan <- productId
// 	// <-respChan
// 	fmt.Println(app.user)
// 	writer.Write([]byte(productInfo))
// }
