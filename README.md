# Go Clean Architecture Example

## Run API

  make go/run/api

## Run tests

  make go/test

## API requests 

### Add book

```
curl -X "POST" "http://localhost:9000/v1/book" \
     -H 'Content-Type: application/json' \
     -H 'Accept: application/json' \
     -d $'{
  "title": "I Am Ozzy",
  "author": "Ozzy Osbourne",
  "pages": 294,
  "quantity":10
}'
```
### Search book

```
curl "http://localhost:9000/v1/book?title=ozzy" \
     -H 'Content-Type: application/json' \
     -H 'Accept: application/json'
```

### Show books

```
curl "http://localhost:9000/v1/book" \
     -H 'Content-Type: application/json' \
     -H 'Accept: application/json'
```

### Add user

```
curl -X "POST" "http://localhost:9000/v1/user" \
     -H 'Content-Type: application/json' \
     -H 'Accept: application/json' \
     -d $'{
  "email": "ozzy@metal.net",
  "first_name": "Ozzy",
  "last_name": "Osbourne",
  "password": "bateater666"
}'

```
### Search user

```
curl "http://localhost:9000/v1/user?name=ozzy" \
     -H 'Content-Type: application/json' \
     -H 'Accept: application/json'
```

### Show users

```
curl "http://localhost:9000/v1/user" \
     -H 'Content-Type: application/json' \
     -H 'Accept: application/json'
```


### Borrow a book

```
curl "http://localhost:9000/v1/loan/borrow/be8b1757-b043-4dbd-b873-63fa9ecd8bb1/282885d7-5d5e-4205-87eb-edc2b2ac5022" \
     -H 'Content-Type: application/json' \
     -H 'Accept: application/json'
```

### Return a book

```
curl "http://localhost:9000/v1/loan/return/be8b1757-b043-4dbd-b873-63fa9ecd8bb1" \
     -H 'Content-Type: application/json' \
     -H 'Accept: application/json'
```
