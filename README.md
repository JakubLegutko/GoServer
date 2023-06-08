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

Server accepts login to /login endpoint, after logging in the /restricted endpoint can be accessed
curl -X POST -d "username=jon&password=shhh!" http://localhost:8090/login
To access the /restricted section, you need to add the token returned from login endpoint to your request
curl -X GET -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiSm9uIFNub3ciLCJhZG1pbiI6dHJ1ZSwiZXhwIjoxNjg2NDg1ODQ3fQ.yty1LIr9LcHuECtG_A8FBMn3Zuw2qlkniO7qSxEXz-M" http://209.38.238.249:8090/restricted
