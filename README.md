# GoServer

Example curl to add a product:

curl --location --request POST 'http://localhost:8090/products' \
--header 'Content-Type: application/json' \
--data-raw '{
    "id": "1",
    "name": "Product 1",
    "price": 9.99
}'

data can be accessed in the browser under localhost:8090/products
