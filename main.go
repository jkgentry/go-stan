package main

import (
    "context"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "math/rand"
    "net/http"
    "os"

    "github.com/gin-gonic/gin"
    "github.com/shurcooL/githubql"
    "golang.org/x/oauth2"
)

type StanLeeRequest struct {
    Orgs []string `formData:"orgs" json:"orgs" binding:"required"`
}

type Characters struct {
    Characters []string
}

type RepoQuery struct {
    Search struct {
        Nodes []Nodes
    } `graphql:"search(first: 100, type: REPOSITORY, query: $query)"`
}

type Nodes struct {
    Repository struct {
        Name githubql.String
    } `graphql:"... on Repository"`
}

func main() {
    r := gin.Default()
    r.POST("/stanlee", func(c *gin.Context) {

        var request StanLeeRequest
        err := c.BindJSON(&request)
        if err != nil {
            log.Fatal(err)
        }
        log.Printf("%v", request.Orgs[0])
        repos := getRepos(request.Orgs[0])
        fmt.Println(repos.Search.Nodes)

        unparsedJson, err := ioutil.ReadFile("./marvel-characters.json")
        if err != nil {
            log.Fatal(err)
        }
        var characters Characters
        json.Unmarshal(unparsedJson, &characters)

        number := rand.Intn(len(characters.Characters))
        for contains(repos.Search.Nodes, characters.Characters[number]) {
            number = rand.Intn(len(characters.Characters))
        }
        c.JSON(http.StatusOK, gin.H{
            "Name": characters.Characters[number],
        })
    })

    r.Run("0.0.0.0:8090") // listen and serve on 0.0.0.0:8080
}

func getRepos(owner string) RepoQuery {
    src := oauth2.StaticTokenSource(
        &oauth2.Token{AccessToken: os.Getenv("OAUTH_TOKEN")},
    )
    httpClient := oauth2.NewClient(context.Background(), src)
    client := githubql.NewEnterpriseClient(os.Getenv("GITHUB_API_URL"), httpClient)
    query := fmt.Sprintf("org:%s", owner)

    var q RepoQuery
    variables := map[string]interface{}{
        "query": githubql.String(query),
    }

    err := client.Query(context.Background(), &q, variables)
    if err != nil {
        log.Fatal(err)
    }
    return q
}

func contains(arr []Nodes, str string) bool {
    for _, a := range arr {
        if string(a.Repository.Name) == str {
            return true
        }
    }
    return false
}
