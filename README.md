# Signbot
Signs messages with a provided ED25519 private key

## Routes
### /
Health check
Accepts no parameters
#### Returns
```
alive
```
### /sign
Signs a message with a provided private key in PEM format
#### Parameters
Accepts a JSON object with the following fields:

| Name    | Type   | Example       | Notes                                                                                                                       |
|---------|--------|---------------|-----------------------------------------------------------------------------------------------------------------------------|
| message | string | Hello, world! |                                                                                                                             |
| key     | string | MC4CAQA...    | Key must be provided in PEM format, beginning -----BEGIN PRIVATE KEY----- and ending -----END PRIVATE KEY----- are optional |

Example:
```json
{
	"message":  "Hello, world!",
	"key":  "MC4CAQAwBQYDK2VwBCIEIPWYB60kck3VNF+wDrdQSf60lwlOLC1OV3EHkllVnbzd"
}
```
#### Returns
Returns a JSON object with the following fields:

| Name      | Type   | Example      | Notes                      |
|-----------|--------|--------------|----------------------------|
| signature | string | TLqRIveio... | Signature in Base64 format |

Example:
```json
{
	"signature":  "TLqRIveiotDUPw3aYJSibvd3Np4Xo2djsJ+HgeOYP+Jo/O8uYKazPbyLF9WEfHnQhvsNjgQXOqTZB7Ut6NkcBw=="
}
```
### /verify
Verifies if a signed message is valid for a provided private key
#### Parameters
Accepts a JSON object with the following fields:

| Name      | Type   | Example       | Notes                                                                                                                       |
|-----------|--------|---------------|-----------------------------------------------------------------------------------------------------------------------------|
| message   | string | Hello, world! |                                                                                                                             |
| key       | string | MC4CAQA...    | Key must be provided in PEM format, beginning -----BEGIN PRIVATE KEY----- and ending -----END PRIVATE KEY----- are optional |
| signature | string | TLqRIveio...  | Signature in Base64 format                                                                                                  |

Example:
```json
{
	"message":  "Hello, world!",
	"key":  "MC4CAQAwBQYDK2VwBCIEIPWYB60kck3VNF+wDrdQSf60lwlOLC1OV3EHkllVnbzd",
	"signature":  "TLqRIveiotDUPw3aYJSibvd3Np4Xo2djsJ+HgeOYP+Jo/O8uYKazPbyLF9WEfHnQhvsNjgQXOqTZB7Ut6NkcBw=="
}
```
#### Returns
Returns a JSON object with the following fields:

| Name  | Type    | Example | Notes |
|-------|---------|---------|-------|
| valid | boolean | true    |       |

Example:
```json
{
	"valid":  true
}
```




