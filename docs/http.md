# HTTP API Specification

### HTTP Status Code

| HTTP Status Code | Description  |
|------------------|--------------|
| `200`            | OK           |
| `400`            | Bad Request  |
| `401`            | Unauthorized |
| `403`            | Forbidden    |
| `404`            | Not Found    |
| `500`            | Server Error |


### Response JSON Structure

| Field Name   | Type   | Requirement | Description       |
|--------------|--------|-------------|-------------------|
| `code`       | int    | Required    | Error Code        |
| `message`    | string | Required    |                   |
| `data`       | object | Optional    |                   |
| `errors`     | object | Optional    | Validation Errors |
| `request_id` | string | Required    |                   |


### Error Code

| Error Code | HTTP Status Code | Description           |
|------------|------------------|-----------------------|
| `102000`   | `200`            | Success               |
| `102001`   | `200`            | Failure               |
| `104000`   | `400`            | Bad Request           |
| `104010`   | `401`            | Unauthorized          |
| `104011`   | `401`            | Access Token Expired  |
| `104012`   | `401`            | Refresh Token Expired |
| `104030`   | `403`            | Forbidden             |
| `104040`   | `404`            | Not Found             |
| `105000`   | `500`            | Server Error          |


### Examples

Example 1:

```json
{
    "code": 102000,
    "message": "The message of success.",
    "request_id": "37654b5b-e0c4-431b-9d85-ef61e5bf3a08"
}
```

Example 2:

```json
{
    "code": 102000,
    "message": "Success.",
    "data": {
        "username": "Jeff"
    },
    "request_id": "660e2755-b7ae-4bab-9ee8-f1efba096b60"
}
```

Example 3:

```json
{
    "code": 102001,
    "message": "The message of failure.",
    "request_id": "d4aee94e-32a9-4252-997f-8f8304210eb9"
}
```

Example 4:

```json
{
    "code": 104000,
    "message": "The verification failed, please check errors for details.",
    "errors": {
        "account": [
            "The account field is required"
        ]
    },
    "request_id": "3d0ff64a-c99b-4c2c-bd4b-5206dec1935f"
}
```
