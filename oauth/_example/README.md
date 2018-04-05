go google oauth example
=======================

The goal of this example is to show a soup to nuts implementation of google oauth.

## Instructions:

1. You'll want to create a `config.yml` file in this folder, that will be where you put things like the google client / secret pair and other options.

An example `config.yml`:

```yaml
googleAuth:
  clientID: "< YOUR CLIENT_ID HERE >"
  clientSecret: "< YOUR CLIENT_SECRET HERE >"
  hostedDomain: "< YOUR HOSTED DOMAIN HERE (ex: blend.com) >"

web:
  port: 5000 #up to you, this is a pretty standard default.
  baseURL: "< YOUR BASE URL HERE >"
```

Caveats:
- You must configure an oauth credential set in the [google cloud console](https://console.cloud.google.com/apis/credentials). 
- You should add a `local.app.xyz` entry to that cloud configuration, and add a hosts entry map for that dns name to point to local host in `/etc/hosts`
    > 127.0.0.1 local.app.xyz
- You should then set your `web > baseURL` entry to `http://local.app.xyz`
    > it is not recommended to run your app on http outside your own machine. 

2. Start the app with:
```bash
> go run main.go
```

3. You should be able to go to `http://local.app.xyz` in a browser and have it talk to the (very basic) app.
