# graphql-tutorial
GraphQL tutorial from Scratch


## Curl Example
```
curl --location --request GET 'http://localhost:8080/graphql' \
--header 'Content-Type: application/json' \
--data-raw '{"query":"{\n\ttutorial(id:1) {\n        id\n        title\n        author {\n            Name\n            Tutorials\n            }\n\t}\n    hello\n    list {\n        id\n        title\n    }\n}\n    ","variables":{}}'
```

## Link Tutorials
- [https://www.youtube.com/watch?v=AlLBG6HrE7E](https://www.youtube.com/watch?v=AlLBG6HrE7E)
