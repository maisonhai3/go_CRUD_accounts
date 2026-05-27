package main

import (
"net/http"
"time"
"github.com/gin-gonic/gin"
"github.com/google/uuid"
"fmt"
)


type account struct {
    Id        string    `json:"id"`
    Name      string    `json:"name"`
    Currency  string    `json:"currency"`
    Balance   int       `json:"balance"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    DeletedAt time.Time `json:"deleted_at"`
}

var accounts = []account{
	{Id: uuid.NewString(), Name: "mai", Currency: "USD", Balance: 1000, CreatedAt: time.Now()},
}

func getAllAccounts(c *gin.Context){
	c.IndentedJSON(http.StatusOK, accounts)
	ctx := c.Request.Context()
	fmt.Println(c)
	fmt.Println(ctx)
}

func getAccount(c *gin.Context){
	id := c.Param("id")
	for _, a := range accounts{
		if a.Id == id {
			c.JSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"error":"account not found"})
}

func createAccount(c *gin.Context){
	var account account

	if err := c.BindJSON(&account); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accounts = append(accounts, account)
	c.IndentedJSON(http.StatusCreated, account)
}

func main(){
	router := gin.Default()
	router.GET("/accounts", getAllAccounts)
	router.GET("/accounts/:id", getAccount)
	router.POST("/accounts", createAccount)
	router.Run("localhost:8080")
}