package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/graphql-go/graphql"
)

type Tutorial struct {
	ID       int
	Title    string
	Author   Author
	Comments []Comment
}

type Author struct {
	Name      string
	Tutorials []int
}

type Comment struct {
	Body string
}

func populate() []Tutorial {
	author := Author{Name: "Yogie", Tutorials: []int{1}}
	tutorial := Tutorial{
		ID:     1,
		Title:  "Go GraphQL Tutorial",
		Author: author,
		Comments: []Comment{
			Comment{
				Body: "First Comment",
			},
		},
	}
	var tutorials []Tutorial
	tutorials = append(tutorials, tutorial)
	return tutorials
}

func main() {
	fmt.Println("GraphQL Tutorial")

	tutorials := populate()

	var commentType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Comment",
			Fields: graphql.Fields{
				"body": &graphql.Field{
					Type: graphql.String,
				},
			},
		},
	)

	var authorType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Author",
			Fields: graphql.Fields{
				"Name": &graphql.Field{
					Type: graphql.String,
				},
				"Tutorials": &graphql.Field{
					Type: graphql.NewList(graphql.Int),
				},
			},
		},
	)

	var tutorialType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Tutorial",
			Fields: graphql.Fields{
				"id": &graphql.Field{
					Type: graphql.Int,
				},
				"title": &graphql.Field{
					Type: graphql.String,
				},
				"author": &graphql.Field{
					Type: authorType,
				},
				"comments": &graphql.Field{
					Type: graphql.NewList(commentType),
				},
			},
		},
	)

	fields := graphql.Fields{
		"hello": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return "World", nil
			},
		},
		"tutorial": &graphql.Field{
			Type:        tutorialType,
			Description: "Get Tutorial By ID",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id, ok := p.Args["id"].(int)
				if ok {
					for _, tutorial := range tutorials {
						if int(tutorial.ID) == id {
							return tutorial, nil
						}
					}
				}

				return nil, nil
			},
		},
		"list": &graphql.Field{
			Type:        graphql.NewList(tutorialType),
			Description: "Get Full Tutorial List",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return tutorials, nil
			},
		},
	}

	// define the object config
	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	// defines a schema config
	schemaConfig := graphql.SchemaConfig{
		Query: graphql.NewObject(rootQuery),
	}
	// creates out schema
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("Failed to create new GraphQL Schema, err: %v", err)
	}

	// example query
	_ = `
	{
		tutorial(id:1) {
			id
			title
			author {
				Name
				Tutorials
			}
		}
	}
	`

	http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		var wrappedRequest map[string]string

		body, _ := ioutil.ReadAll(r.Body)
		err = json.Unmarshal(body, &wrappedRequest)

		resp := graphql.Do(graphql.Params{
			Schema:        schema,
			RequestString: wrappedRequest["query"],
		})

		b, _ := json.Marshal(resp)

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(200)
		w.Write(b)
	})

	http.ListenAndServe(":8080", nil)

}
