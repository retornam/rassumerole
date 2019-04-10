## Config

This assumes you have already setup the needed files  under `~/.aws` specifically
the `~/.aws/config` and the `~/.aws/credentials` files


##  Install


Clone this repo

```sh 
   cd  rassumerole
   go build 
   cp rassumerole /usr/local/bin/
```


## Usage

```sh
    eval $(rassumerole exampleprofile)
```

thats it.


