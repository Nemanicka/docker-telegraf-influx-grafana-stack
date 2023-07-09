# Stress Tests Result

Siege commmands for getting results below are:
```bash
$ siege -c <concurrency> -r 10 -d 0.1 -f urls.txt
```


| Concurrency | Transactions | Availability | Elapsed Time | Data Transferred | Response Time | Transaction Rate | Throughput | Concurrency | Successful Transactions | Failed Transactions | Longest Transaction | Shortest Transaction |
|-------------|--------------|--------------|--------------|-----------------|---------------|------------------|------------|-------------|------------------------|---------------------|---------------------|----------------------|
| 10          | 100          | 100.00       | 0.66         | 0.09            | 0.00          | 151.52           | 0.13       | 0.17        | 100                    | 0                   | 0.01                | 0.00                 |
| 50          | 500          | 100.00       | 0.74         | 2.59            | 0.00          | 675.68           | 3.50       | 1.92        | 500                    | 0                   | 0.03                | 0.00                 |
| 100         | 1000         | 100.00       | 0.89         | 12.45           | 0.01          | 1123.60          | 13.99      | 11.11       | 1000                   | 0                   | 0.14                | 0.00                 |
| 300         | 2652         | 88.40        | 6.80         | 665.95          | 0.67          | 390.00           | 97.93      | 261.03      | 2652                   | 348                 | 6.06                | 0.00                 |
| 400         | 660          | 39.19        | 8.64         | 946.45          | 4.53          | 76.39            | 109.54     | 345.70      | 660                    | 1024                | 8.48                | 0.00                 |

## Setting Up

```bash
$ docker compose up
```

## Important

This setup uses host network, so make sure all necessary port are available.
These ports are:
80, 27017, 28017, 9200, 8086, 8125, 3000, 3030

## API

The project backend is mostly copied from [this repo](https://github.com/brunaobh/go-mongodb-rest-api-crud)
It emulates CRUD operation for flight management.
Other than that, all logs are written to elastic search.

The only two requests used for load testing are:
* GET /flights/all - get all flights in db
* POST /flights    - store a flight information

### Load Testing

Load Testing was performed via command

```bash
$ ./siegeLoadTestGET.sh & ./siegeLoadTestPOST.sh
```

All the results can be found in /screenshots directory
