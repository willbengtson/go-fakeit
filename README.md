# go-fakeit 
A simple redirector for red team simulation to redirect from fake site to real site with auto Let's Encrypt cert generation

## Usage

```
Usage of ./go-fakeit:
  -fake-domain string
    	Comma delimited list of domains to host
  -local
    	Local server for testing
  -redirect string
    	Domain to redirect to
```

## "Production"

`./gofake-it -fake-domain badsite.com -redirect realsite.com`

If you want to suppport `www` and the `apex` then,

`./gofake-it -fake-domain badsite.com,www.badsite.com -redirect realsite.com`

## Development - local 

Build to the binary

1. go get
2. go build

Run the server locally

`./go-fakeit -fake-domain badsite.com -redirect realsite.com -local`