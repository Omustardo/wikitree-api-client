# WikiTree Proto and Clients

This project is a proof of concept for using Protobuf to simplify using [WikiTree's HTTP API](https://github.com/wikitree/wikitree-api) for non-Javascript code.

The WikiTree API uses GET and POST requests with JSON blobs. To use these in a type-safe manner (in other words, using structs or classes rather than dealing with JSON) would require writing a lot of boilerplate code to handle parsing JSON. By having the JSON blobs represented as messages in a [Protocol Buffer](https://developers.google.com/protocol-buffers), it is theoretically easy to generate boilerplate code for most well supported languages using the protobuf compiler (`protoc`).

A bit more context on why protocol buffers are nice, and the general thoughts behind creating this project: [WikiTree G2G Thread](https://web.archive.org/web/20220204004230/https://www.wikitree.com/g2g/1363173/is-there-proto-definition-for-the-api-would-worth-having-one)

## General Project Structure

`wikitree.proto` defines request/response messages and a service. Using the protobuf compiler (see the `Golang Proto Compilation` section below), the messages are turned into Golang structs while the service is turned into an interface.  The generated code is submitted under `go/proto/`

`go/wikiclient/client.go` implements the service interface. It supports converting request messages into URL parameters, executing HTTP GET requests for the URL+parameters, and converting the JSON responses into structs (the ones that were generated from Proto messages).

`go/examples/example.go` creates a new client and uses it to do a query.

## Proof of Concept Findings

The general finding is that it is possible but not trivial to use protos for accessing WikiTree's API. If there were a strong demand for non-javascript interaction with the API, I think it is a good way to go. The better alternative would be for WikiTree to adopt a gRPC / Proto API itself in some future release (WikiTree API v2?). This would still provide support a purely HTTP JSON API, though with all `snake_case` fields. It would make usage from other clients much easier.

### Naming

The main issue is inconsistent and incompatible naming of parameters in WikiTree API requests and in the JSON responses. Proto to JSON conversion (and vice-versa) requires that fields be named in `snake_case`. WikiTree uses a mix of `snake_case`, `lowerCamelCase`, `UpperCamelCase`, `lowerCamelCaseWithACRONYMInCaps`, and its own `Custom_Style`.

It would be possible to work around this with an extra layer of code that maps field names to the necessary format. This could be implemented as a Map between field names, like `{"resolveRedirect": "resolve_redirect", "Privacy_IsPrivate": "privacy_is_private"}`. To use this mapping in code, it would need to be used in both the request generation and to modify the JSON returned by the WikiTree API. While creating this mapping is a one-off task, using it significantly increases the amount of work required to write clients in each individual language.

### Auth

Another issue is that the WikiTree API requires authentication using a username and password for any restricted/private data. While not too difficult to implement, this is another thing that would need to be implemented in the client code for each individual language.

### Types

In the proof of concept, I only used string fields in the proto. This is obviously not ideal, as using bools/ints/floats/lists/etc would make usage much easier. Using the correct types might be complicated. For example, WikiTree inconsistently uses "true", "false", 0, 1 in boolean fields. Lists (like in getProfile's "fields" parameter) are expected to be comma separated values. Figuring out the right way to represent these in a proto would require additional thought.

## Usage

`go run go/examples/example.go`

Output:

```
Starting WikiTree client
== Raw JSON Response
{"page_name":"Kennedy-21529","profile":{"Id":30362890,"Name":"Kennedy-21529","FirstName":"Cunningham","MiddleName":"","LastNameAtBirth":"Kennedy","LastNameCurrent":"Kennedy","Nicknames":"","LastNameOther":"","RealName":"Cunningham","Prefix":"","Suffix":"","BirthLocation":"Shankill Parish, County Armagh, Ireland","DeathLocation":"Mexico","Gender":"Male","BirthDate":"1826-02-23","DeathDate":"1893-07-24","BirthDateDecade":"1820s","DeathDateDecade":"1890s","Photo":null,"IsLiving":0,"HasChildren":1,"NoChildren":0,"Privacy":60,"IsPerson":1,"Touched":"20220202003422","ShortName":"Cunningham Kennedy","BirthNamePrivate":"Cunningham Kennedy","LongNamePrivate":"Cunningham Kennedy","LongName":"Cunningham Kennedy","BirthName":"Cunningham Kennedy","Manager":30140149,"DataStatus":{"FirstName":"","MiddleName":"","LastNameAtBirth":"","LastNameCurrent":"","RealName":"","ColloquialName":"","Nicknames":"","Gender":"","BirthDate":"","DeathDate":"","BirthLocation":"","DeathLocation":"","Father":"","Mother":"","Spouse":"","Prefix":"","Suffix":"","LastNameOther":""},"Privacy_IsPrivate":false,"Privacy_IsPublic":false,"Privacy_IsOpen":true,"Privacy_IsAtLeastPublic":true,"Privacy_IsSemiPrivate":false,"Privacy_IsSemiPrivateBio":false,"Father":30363211,"Mother":30363195,"Parents":{"30363211":{"Id":30363211,"Name":"Kennedy-21535","FirstName":"George","MiddleName":"Edward","LastNameAtBirth":"Kennedy","LastNameCurrent":"Kennedy","Nicknames":"","LastNameOther":"","RealName":"George","Prefix":"","Suffix":"","BirthLocation":"Scotland","DeathLocation":"Scotland","Gender":"Male","BirthDate":"1791-00-00","DeathDate":"1840-00-00","BirthDateDecade":"1790s","DeathDateDecade":"1840s","Photo":null,"IsLiving":0,"HasChildren":1,"NoChildren":0,"Privacy":60,"IsPerson":1,"Touched":"20210507085134","ShortName":"George Kennedy","BirthNamePrivate":"George Kennedy","LongNamePrivate":"George E. Kennedy","LongName":"George Edward Kennedy","BirthName":"George Edward Kennedy","Manager":30140149,"DataStatus":{"FirstName":"","MiddleName":"","LastNameAtBirth":"","LastNameCurrent":"","RealName":"","ColloquialName":"","Nicknames":"","Gender":"","BirthDate":"guess","DeathDate":"guess","BirthLocation":"","DeathLocation":"","Father":"","Mother":"","Spouse":"","Prefix":"","Suffix":"","LastNameOther":""},"Privacy_IsPrivate":false,"Privacy_IsPublic":false,"Privacy_IsOpen":true,"Privacy_IsAtLeastPublic":true,"Privacy_IsSemiPrivate":false,"Privacy_IsSemiPrivateBio":false,"Father":0,"Mother":0},"30363195":{"Id":30363195,"Name":"Cunningham-14086","FirstName":"Jane","MiddleName":"","LastNameAtBirth":"Cunningham","LastNameCurrent":"Cunningham","Nicknames":"","LastNameOther":"","RealName":"Jane","Prefix":"","Suffix":"","BirthLocation":"Shankill Parish, County Armagh, Ireland","DeathLocation":"Murphysboro, Jackson County, Illinois, USA","Gender":"Female","BirthDate":"1792-06-05","DeathDate":"1875-06-23","BirthDateDecade":"1790s","DeathDateDecade":"1870s","Photo":null,"IsLiving":0,"HasChildren":1,"NoChildren":0,"Privacy":60,"IsPerson":1,"Touched":"20210507085134","ShortName":"Jane Cunningham","BirthNamePrivate":"Jane Cunningham","LongNamePrivate":"Jane Cunningham","LongName":"Jane Cunningham","BirthName":"Jane Cunningham","Manager":30140149,"DataStatus":{"FirstName":"","MiddleName":"","LastNameAtBirth":"","LastNameCurrent":"","RealName":"","ColloquialName":"","Nicknames":"","Gender":"","BirthDate":"","DeathDate":"","BirthLocation":"","DeathLocation":"","Father":"","Mother":"","Spouse":"","Prefix":"","Suffix":"","LastNameOther":""},"Privacy_IsPrivate":false,"Privacy_IsPublic":false,"Privacy_IsOpen":true,"Privacy_IsAtLeastPublic":true,"Privacy_IsSemiPrivate":false,"Privacy_IsSemiPrivateBio":false,"Father":30394824,"Mother":30394804}}},"status":0}
== Proto Response:
page_name:"Kennedy-21529" profile:{}
```

Note that the JSON to Proto conversion only works for fields which use `snake_case`, which is unfortunately very few fields.

## Golang Proto Compilation

Starting with a `.proto` file, the Proto Compiler (`protoc`) can turn it into your language of choice. This section is purely informational, as the generated code is checked-in to this repository. Look for `.pb.go` files.

This process requires installing `protoc` as well as a language-specific plugins. 

1. Download `protoc`: https://developers.google.com/protocol-buffers/docs/downloads
1. Install the plugin for Golang (or your language of choice) by starting the relevant tutorial: https://developers.google.com/protocol-buffers/docs/tutorials. For Golang, this involves running `go install google.golang.org/protobuf/cmd/protoc-gen-go@latest`
1. Run the compiler (`protoc`). The following command works to dump generated files into `go/proto/`: `protoc --go_out=paths=source_relative:./go/proto -I. wikitree.proto`

This works to generate code for the messages (aka structs), but in order to generate the service and client interfaces an additional plugin must be installed.

1. Run the commands to install gRPC. For Go, this follows https://www.grpc.io/docs/languages/go/quickstart/
   1. `go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1`
1. Use the `go-grpc_out` parameter when running `protoc`:
   1. For Go, `protoc --go-grpc_out=./go/proto -I. wikitree.proto`

These two invocations of `protoc` can be combined into one: `protoc --go_out=paths=source_relative:./go/proto --go-grpc_out=paths=source_relative:./go/proto -I. wikitree.proto`
