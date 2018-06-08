# graphb: A GraphQL query builder
__Focus on the query, not string manipulation__

## Motivation
As Go developers, when there is a GraphQL server, we often build our query string like this:
```go
// Define this in some file x.go
const queryTemplate = `
    "query": "
        query an_operation_name { 
            a_field_name (
                an_argument_name_if_any: \"%s\"
            ) {
                a_sub_field,
                another_sub_field
            }
        }"`
    
// DO this in some other file y.go
func buildQuery(value string) string {
    return fmt.Sprintf(queryTemplate, value)
}
```
This approach is verbose, inflexible and error prone.

It is verbose and inflexible because every time you want to query a different structure, you need to rewrite another template and another string interpolation function. Not to mention when you have argument types such lists or enums or when you want to use fragments, directives and other syntax of GraphQL.

It is error prone because there is no way to check the correctness of your syntax until the query string is send to the server.

It also wastes extra spaces and looks not beautiful.

Therefore, when it comes to GraphQL client in Go, a developer spends much of her time to fight the string manipulation war (that's easy to lose), rather than focusing on the query (aka the business logic).

This library solves the string building problem so that you focus on business logic.

## Example
All code are well documented. See [example](example) dir for more examples.

The lib provides 3 ways of constructing a query.
1. Method Chaining
2. Functional Options
3. Struct Literal

#### 1. Method Chaining
See [example/three_ways_to_construct_query_test.go#L10-L43](example/three_ways_to_construct_query_test.go#L10-L43)

#### 2. Functional Options
See [example/three_ways_to_construct_query_test.go#L45-L74](example/three_ways_to_construct_query_test.go#L45-L74)

#### 3. Struct Literal
See [example/three_ways_to_construct_query_test.go#L76-L100](example/three_ways_to_construct_query_test.go#L76-L100)

### Words from the author
__The library catches cycles.__ That is, if you have a `Field` whose sub Fields can reach the `Field` itself, the library reports an error.

I hesitate to make all fields private and only allow constructing a query through `NewQuery` and `MakeQuery`. I also don't know if `MakeQuery` is better than `NewQuery` or the contrary. __Please use it and give feedback__.

## Error Handling
All `graphb` errors are wrapped by [pkg/errors](https://github.com/pkg/errors).  
All error types are defined in [error.go](error.go)

## Test
`graphb` uses [testify/assert](https://github.com/stretchr/testify/#assert-package).
```bash
go test
```

## Todos
The library does not currently support:
1. Directive
2. Variables
3. Fragments

I do not know how useful would them be for a user of this library. Since the library builds the string for you, you sort of get the functionality of Variable and Fragment for free. You can just reuse a Field or the values of Fields and Arguments as normal Go code. Directive might be the most useful one for this library.
