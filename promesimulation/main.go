package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"upperpeng.com/database"
)

type album struct {
	ID string `json:"id"`
	Title string `json:"title"`
	Artist string `json:"artist"`
	Price float64 `json:"price"`
}

func typeConvertToDatabse(in album) database.Album{
	var result database.Album
	result.Title = in.Title
	result.Artist = in.Artist
	result.Price = in.Price
	return result
}

var httpReqs  *prometheus.CounterVec
var all prometheus.Counter
var id prometheus.Counter
var del prometheus.Counter
var add prometheus.Counter
var ItemsNum int = 0
var ItemIds []int64 = make([]int64,10) //数据库里的id
func init() {
	httpReqs := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "upperpeng_first",
			Help: "How many HTTP requests processed, partitioned by status code and HTTP method.",
		},
		[]string{"name", "method"},
	)
	prometheus.MustRegister(httpReqs)
	//httpReqs.WithLabelValues("404", "POST").Add(42)
	all = httpReqs.WithLabelValues("all","GET")
	id = httpReqs.WithLabelValues("id","GET")
	del = httpReqs.WithLabelValues("del","GET")
	add = httpReqs.WithLabelValues("add","POST")
}
//curl http://localhost:8080/albums -X POST -H "Content-Type: application/json" -d '{"id":"100","title":"a song","artist":"elsesv","price":49.99}'
func main() {
	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.POST("/albums",postAlbums)
	router.GET("/albums/:id", getAlbumByID)
	router.GET("/albums/delete", delAlbumByID)

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	router.Run("0.0.0.0:8080")

}



func getAlbums(c *gin.Context) {
	albums,_:= database.GetALLAlbumsInfo()
	ItemsNum = len(albums)
	for _,al := range albums{
		ItemIds = append(ItemIds,al.ID)
	}
	c.IndentedJSON(http.StatusOK, albums)
	all.Inc()
}

func getAlbumByID(c *gin.Context) {

	Id := c.Param("id")
	intId,err := strconv.Atoi(Id)
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
	}else{
		album,_ := database.GetAlbumInfoByID(int64(intId))
		c.IndentedJSON(http.StatusOK,album)
	}
	id.Inc()
}

func delAlbumByID(c *gin.Context) {
	ItemIds = nil
	albums,_:= database.GetALLAlbumsInfo()
	ItemsNum = len(albums)
	for _,al := range albums{
		ItemIds = append(ItemIds,al.ID)
	}


	index := rand.Intn(ItemsNum)
	id := ItemIds[index]
	database.DelAlbumInfoByID(int64(id))
	c.IndentedJSON(http.StatusOK,gin.H{"message": fmt.Sprint(id)+"already delete"})

	del.Inc()
}


func postAlbums(c *gin.Context){
	var newAlbum album
	if err := c.BindJSON(&newAlbum);err!=nil{
		c.IndentedJSON(http.StatusOK,gin.H{"message": "json bind failure!"})
	}
	if _,err:=database.AddAlbum(typeConvertToDatabse(newAlbum));err!=nil{
		c.IndentedJSON(http.StatusExpectationFailed, gin.H{"message": "add failure!"})
	}else{
		c.IndentedJSON(http.StatusCreated,newAlbum)
	}

	ItemIds = nil
	albums,_:= database.GetALLAlbumsInfo()
	ItemsNum = len(albums)
	for _,al := range albums{
		ItemIds = append(ItemIds,al.ID)
	}

	add.Inc()
}
