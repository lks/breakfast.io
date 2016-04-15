package main

import (
  "database/sql"
  "gopkg.in/gorp.v1"
  _ "github.com/lib/pq"
  "github.com/gin-gonic/gin"
  "time"
  "log"
)

func main() {
  router := gin.Default()
  router.POST("/members", CreateMember)
  router.Run(":8080")
}

func CreateMember(c *gin.Context) {
  var member Member
  c.Bind(&member)

  if member.Firstname != "" &&  member.Lastname != "" && member.Email != ""{
    content := newMember(member.Firstname, member.Lastname, member.Email)
    err := dbmap.Insert(&content)
    checkErr(err, "Insert failed")
    c.JSON(201, content)
  } else {
    c.JSON(422, gin.H{"error": "Fields are empty"})
  }
}

type Member struct {
  Id int64 `db: "member_id"`
  Firstname string `db:"firstname" json:"firstname"`
  Lastname  string `db:"lastname" json:"lastname"`
  Email  string `db:"email" json:"email"`
  Count int64 `db:"count" json:"count"`
  Created int64 `db:"created" json:"created"`
}

func newMember(Firstname string, Lastname string, Email string) Member {
  return Member{
    Firstname: Firstname,
    Lastname: Lastname,
    Email: Email,
    Count: 0,
    Created: time.Now().UnixNano(),
  }
}

func initDb() *gorp.DbMap {
  db, err := sql.Open("postgres", "user=j.calabrese dbname=breakfast-app sslmode=disable")
	checkErr(err, "sql.Open failed")

  // construct a gorp DbMap
  dbmap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}

  // add a table, setting the table name to 'posts' and
  // specifying that the Id property is an auto incrementing PK
  dbmap.AddTableWithName(Member{}, "members").SetKeys(true, "Id")

  // create the table. in a production system you'd generally
  // use a migration tool, or create the tables via scripts
  err = dbmap.CreateTablesIfNotExists()
  checkErr(err, "Create tables failed")

  return dbmap
}

func checkErr(err error, msg string) {
    if err != nil {
        log.Fatalln(msg, err)
    }
}

// initialize the DbMap
var dbmap = initDb()
