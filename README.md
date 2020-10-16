`avx` Aviatrix CLI
=

Disclaimer: This is a personal project and is in no way affiliated with Aviatrix Systems Inc.

`avx` is a CLI for interacting with the Aviatrix RPC API https://api.aviatrix.com/

This is done through the Aviatrix SDK that is maintained for use in the Aviatrix
terraform provider https://github.com/AviatrixSystems/terraform-provider-aviatrix

Build & Install
-

### MacOS
```shell script
$ git clone git@github.com:CyrusJavan/avx.git
$ cd avx
$ go build -o /usr/local/bin/avx
```

Required Configuration
-

Required environment variables:

- `AVIATRIX_CONTROLLER_IP`
- `AVIATRIX_USERNAME`
- `AVIATRIX_PASSWORD`

Usage
-

### `avx login`

`avx login` will attempt to login with the provided credentials. If
successful, the CID will be printed out.
```shell script
$ avx login
CID: "MMUyqYcNOjaWUWIFHmYA"
```

---

### `avx rpc <action>`

`avx rpc <action>` will attempt to login and send a POST request to
the API with the provided `action`. `avx` prints out debug information like the
controller IP, request body and response latency. `avx` then prints out the 
response body.
```shell script
$ avx rpc list_accounts
controller IP: 127.0.0.1
request body:
{
  "CID": "soPEtEopZlkC1Vwwzdl4",
  "action": "list_accounts"
}
latency: 153ms
response body:
{
  "return": true,
  "results": {
    "account_list": [
      {
        ...
```

---

### `avx rpc <action> <key>=<value> [<key>=<value>...]`

In this form `avx rpc` will send a POST request with the given action and any extra
params that were provided.
```shell script
$ avx rpc delete_account_profile account_name=john-gcloud
controller IP: 127.0.0.1
request body:
{
  "CID": "CgRVzRukvCtUGLwp80lw",
  "account_name": "tfa-byl0f",
  "action": "delete_account_profile"
}
latency: 113ms
response body:
{
  "return": true,
  "results": "Account deleted successfully."
} 
```

---
