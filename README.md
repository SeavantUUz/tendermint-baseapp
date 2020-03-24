# BaseApp

A basic app from tendermint official tutorial. Extent it with auction.

## build, install

As same as official tutorial, follow the [tutorial](https://tutorials.cosmos.network/) step by step.

### auction

After committed a buy-name operation, using below commands to launch an auction:

```bash
./acli tx nameservice auction jack.id 10nametoken --from jack
```

### bid

After launched an auction, joining to bid by:

```bash
./acli tx nameservice bid jack.id 20nametoken --from alice
```

An auction will be automatically finished after last bid time plus 100 blocks.

### query

Query auctions by:

```bash
./acli query nameservice auctions
```

or a specify auction by:

```bash
./acli query nameservice auction jack.id
```

A rest server as well:

```bash
./acli rest-server
```

and access 

```bash
http://127.0.0.1:1317/nameservice/auctions
```

or

```bash
http://127.0.0.1:1317/nameservice/auction/jack.id
```