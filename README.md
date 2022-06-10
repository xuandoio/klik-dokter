## Installation instruction (required Docker)
```docker-compose up -d```
 
Default port is 8080, so ensure your 8080 port is available or change to your custom port on `docker-compose.yml` file

## Postman Document
`https://documenter.getpostman.com/view/16862423/Uz5MGaD3`

## Postman Environment

```
{
	"id": "039239ee-4131-4036-b28d-d3c9642cefb1",
	"name": "Klik Docter",
	"values": [
		{
			"key": "base_url",
			"value": "http://127.0.0.1:8080/",
			"type": "default",
			"enabled": true
		},
		{
			"key": "access_token",
			"value": "",
			"type": "default",
			"enabled": true
		},
		{
			"key": "product_id",
			"value": "",
			"type": "default",
			"enabled": true
		}
	],
	"_postman_variable_scope": "environment",
	"_postman_exported_at": "2022-06-10T09:22:29.441Z",
	"_postman_exported_using": "Postman/9.20.3"
}

```