Golang Name Parser
===========

[![Build Status](https://travis-ci.org/blendlabs/go-name-parser.svg)](https://travis-ci.org/blendlabs/go-name-parser)

This is a simple library to parse names into their constituent parts.

It is (very largely) based off of [PHP-Name-Parser](https://github.com/joshfraser/PHP-Name-Parser).

##Example

```go
import "github.com/blendlabs/go-name-parser"
//...
name := names.Parse("Mr. Potato McTater, III")
fmt.Printf("%#v\n", name) 
/*
> name{Salutation:"Mr.", FirstName:"Potato", MiddleName:"", LastName:"McTater", Suffix:"III"}
*/
```