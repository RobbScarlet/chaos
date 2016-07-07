package main

import (
	"chaos/api"
	"chaos/config"
	"github.com/emicklei/go-restful"
	"github.com/google/cayley/graph"
	"log"
	"net/http"
)

func main() {
	err := graph.InitQuadStore("bolt", config.Bolt_DB_Path, nil)
	if err != nil && err != graph.ErrDatabaseExists {
		panic(err)
	}

	container := restful.NewContainer()
	err = api.Register(container)
	if err != nil {
		panic(err)
	}

	//config := swagger.Config{
	//	WebServices:    container.RegisteredWebServices(), // you control what services are visible
	//	WebServicesUrl: "http://localhost:8080",
	//	ApiPath:        "/apidocs.json",
     //   //
	//	//// Optionally, specifiy where the UI is located
	//	//SwaggerPath:     "/apidocs/",
	//	//SwaggerFilePath: "/Users/emicklei/xProjects/swagger-ui/dist",
	//}
	//swagger.RegisterSwaggerService(config, container)

	log.Printf("start listening on localhost:8787")
	server := &http.Server{Addr: ":8787", Handler: container}
	log.Fatal(server.ListenAndServe())
}
