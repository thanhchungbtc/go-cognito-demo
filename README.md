# Example of using AWS Cognito with custom authentication

## Prerequisite

- Create user pool in AWS Cognito
  - Do not check generate secret key
  - Create an app client
  - Enable verification method of your choice. We'll be using email in this example

## Features

### Login

User can authenticate with username and password using `USER_PASSWORD_AUTH` flow, or by refresh token using `REFRESH_TOKEN_AUTH`.

Example response

```json
  "access_token": "...",
  "expires_in": 3600,
  "id_token": "...",
  "refresh_token": "...",
  "token_type": "Bearer"
}
```

### Register & verification

Vertification code will be sent as of the method you choose in the first step

### Refresh token

To access protected resources, send the request along with `Authorization` header.
You'll get 401 when the token is expired, in that case just reauthenticate with the `refresh_token`

```js
axios.get("/api/me", {
  headers: { Authorization: `Bearer ${obj.id_token}` },
});
```

To automate this process, we'll create an interceptor.

```js
const createInterceptor = () => {
  const interceptor = axios.interceptors.response.use(
    (res) => res,
    (err) => {
      // not an unauthorized request, continue with error
      if (err.response.status !== 401) {
        return Promise.reject(err);
      }

      // get refresh token
      const data = window.localStorage.getItem("auth");
      const obj = JSON.parse(data);
      const { refresh_token } = obj;
      if (!data) {
        return Promise.reject(err);
      }
      // eject interceptor to prevent loop in case refresh token itself returns 401
      axios.interceptors.response.eject(interceptor);
      return axios
        .post("/api/login", {
          refresh_token,
        })
        .then(({ data }) => {
          // update the authenticated token
          window.localStorage.setItem("auth", JSON.stringify(data));
        })
        .catch((err) => {
          // remove the invalid authorization info
          window.localStorage.removeItem("auth");
          return Promise.reject(err);
        })
        .finally(createInterceptor);
    }
  );
};
```

### Parse the jwt using public key

When creating an user pool, the public keys will be available at `https://cognito-idp.{region}.amazonaws.com/{user_pool_id}/.well-known/jwks.json`

Find the public key with `kid` in the JWT's header. Leverage the existing jwt library, the parsing is simple as

```go
jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
  kid := token.Header["kid"].(string)
  jsonWebKey := a.keys.LookupKeyID(kid)[0]

  publicKey := &rsa.PublicKey{}

  if err := jsonWebKey.Raw(publicKey); err != nil {
    return nil, err
  }
  return publicKey, nil
})
```
