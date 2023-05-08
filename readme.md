# Anvil Estimate

Estimate Transaction Result of Token Change

## Usage:

### Deployment

First, run `build.sh`, then `docker build`.  
You need to set environment variables before run the docker image.  
UPSTREAM_RPC 

### Start Simulation

```http request
POST /api/v1/simulation/run HTTP/1.1
Content-Type: application/json

{
   "from": "from_address",
   "to": "contract_address",
   "amount": "empty string, or your value, ",
   "data": "",
   "token_change": [
       "0x0000000000000000000000000000000000000000",
       "0x0d500B1d8E8eF31E21C99d1Db9A6444d3ADf1270",
   ]
}
```

After start the simulation, you will get a uuid for the request, you can use this uuid to get the result.

### Get Simulation Result

```http request
GET /api/v1/simulation/result/{uuid} HTTP/1.1
Content-Type: application/json
```