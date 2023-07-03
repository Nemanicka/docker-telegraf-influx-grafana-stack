#payload='{"flightNum":"TEST12", "airline":"TAM", "airport":"Confins", "status":"Confirmed", "expected":"20180713083500", "confirmed":"20180713083500"}'

#siege -c 1 -r 1 -g 'http://127.0.0.1/flights/' -p 'http://127.0.0.1/flights/' -d $payload -H 'Content-Type: application/json'
siege -c 20 -r 100 -d 0.5  'http://127.0.0.1/flights/all' -H 'Content-Type: application/json'
#siege -c 5 -r 10 -d 1 'http://127.0.0.1/flights POST {"flightNum":"TEST14", "airline":"TAM", "airport":"Confins", "status":"Confirmed", "expected":"20180713083500", "confirmed":"20180713083500"}'  -H 'Content-Type: application/json'


#siege -c 1 -r 1 -g 'http://127.0.0.1/'

#ab -n 1 -c 1 -H 'Content-Type: application/json' -m GET "http://127.0.0.1/flights" -d $payload "http://127.0.0.1/flights" 
