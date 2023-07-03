# Setting Up

```bash
$ docker compose up
```

# Important

This setup uses host network, so make sure all necessary port are available.
These ports are:
80, 27017, 28017, 9200, 8086, 8125, 3000, 3030

# API

The project backend is mostly copied from [this repo](https://github.com/brunaobh/go-mongodb-rest-api-crud)
It emulates CRUD operation for flight management.
Other than that, all logs are written to elastic search.

The only two requests used for load testing are:
* GET /flights/all - get all flights in db
* POST /flights    - store a flight information
