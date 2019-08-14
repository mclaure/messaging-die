# messaging-die

GO project that implements repply correlation (RPC) pattern for systems integration

![RPC-repply-pattern](https://user-images.githubusercontent.com/24611413/62911259-bc42b600-bd51-11e9-8033-5020fdaff14e.jpg)

## Getting Started

Simple API that provides soccer statistics using a public database
- Source Data: https://www.kaggle.com/hugomathien/soccer
- Data description: http://www.football-data.co.uk/notes.txt

Repository: https://github.com/mclaure/die-mannschaft

![Sync-Communication](https://user-images.githubusercontent.com/24611413/62910994-b8626400-bd50-11e9-923b-ef0d5d8f3c1f.jpg)

### Installing

1) Download the code
2) run the following command:
- E: \Go-Projects\src\github.com\mclaure\messaging-die\rpc>go run client.go
- Starting Public Service

3) run the following command:
- E:\Go-Projects\src\github.com\mclaure\messaging-die\rpc>go run server.go
- 2019/08/12 21:34:03 [*] Consuming from queue (rpc-die-queue)
- 2019/08/12 21:34:04 [S] Waiting for RPC requests

### C4 Model

![C4 Model-01](https://user-images.githubusercontent.com/24611413/62994080-1b72fa00-be28-11e9-8dd0-db8c9b944eda.jpg)
  
![C4 Model-02](https://user-images.githubusercontent.com/24611413/62994093-2ded3380-be28-11e9-90c0-1ea9a4efb5ae.jpg)

![Async-Communication](https://user-images.githubusercontent.com/24611413/62910317-10e43200-bd4e-11e9-8e77-70e31d3794ae.jpg)

### Available APIs

1)  GET /countries
2)  GET /teams
3)  GET /leagues
