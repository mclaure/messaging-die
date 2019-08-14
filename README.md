## messaging-die

GO project that implements repply correlation (RPC) pattern for systems integration

![RPC-repply-pattern](https://user-images.githubusercontent.com/24611413/62911259-bc42b600-bd51-11e9-8033-5020fdaff14e.jpg)

## Getting Started
### The Die Mannschaft Problem
```
Die Mannschaft is Germany's football association that over the past ten years has been including technology into the match field, keeping track of player's statistics, performance, predictions and overall metadata for their internal tournaments.

For a long time they have been keeping all of this data locally, without sharing it with any other organization. In December 2018, the local government has created new laws forcing Die Mannschaft to share all their data with other government institutions. But there's the risk of their infrastructure not being able to support the amount of external requests, they have estimated that at least 5000 organizations will be fetching their data continuously.
```
What approach would you suggest them to use for efficiently sharing their data with other organizations?

## Die MannSchaft : Sync API solution
* Simple API that provides soccer statistics using a public database
** [Source Data] (https://www.kaggle.com/hugomathien/soccer)
** [Data description] (http://www.football-data.co.uk/notes.txt)
** [Repository] (https://github.com/mclaure/die-mannschaft)

![Sync-Communication](https://user-images.githubusercontent.com/24611413/62910994-b8626400-bd50-11e9-923b-ef0d5d8f3c1f.jpg)

## Installing

* Download the code
* run the following command:
- E: \Go-Projects\src\github.com\mclaure\messaging-die\rpc>go run client.go
- Starting Public Service

* run the following command:
- E:\Go-Projects\src\github.com\mclaure\messaging-die\rpc>go run server.go
- 2019/08/12 21:34:03 [*] Consuming from queue (rpc-die-queue)
- 2019/08/12 21:34:04 [S] Waiting for RPC requests

### C4 Model

![C4 Model-01](https://user-images.githubusercontent.com/24611413/62994080-1b72fa00-be28-11e9-8dd0-db8c9b944eda.jpg)
  
![C4 Model-02](https://user-images.githubusercontent.com/24611413/62994093-2ded3380-be28-11e9-90c0-1ea9a4efb5ae.jpg)

![Async-Communication](https://user-images.githubusercontent.com/24611413/62910317-10e43200-bd4e-11e9-8e77-70e31d3794ae.jpg)

## Available APIs

*  GET /countries
*  GET /teams
*  GET /leagues

## Authors

* **Marcelo Claure** - *Initial work*
* **Boris Dominguez** - *Initial work*
* **Claudio Melendres** - *Initial work*
