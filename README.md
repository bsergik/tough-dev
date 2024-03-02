# tough-dev

## SSO Service

- Using [keycloak](https://www.keycloak.org/)
- [Authentication and Authorization (OpenID Connect)of Go Rest Api’s using an open-source IAM called Keycloak](https://medium.com/@allusaiprudhvi999/authentication-and-authorization-in-golang-microservice-using-an-open-source-iam-called-keycloak-46f03a26248f)

```shell
docker run --name mykeycloak -p 8003:8080 -e KEYCLOAK_ADMIN=admin -e KEYCLOAK_ADMIN_PASSWORD=change_me quay.io/keycloak/keycloak start-dev

curl --location --request POST 'http://localhost:8003/auth/realms/demorealm/protocol/openid-connect/token' \
--header 'Content-Type: application/x-www-form-urlencoded' \
--data-urlencode 'client_id=myclientid' \
--data-urlencode 'grant_type=password' \
--data-urlencode 'username=admin' \
--data-urlencode 'password=admin' \
--data-urlencode 'scope=openid' \
--data-urlencode 'client_secret=Cd4QUcjUmbsc5liiKOa8cQPmQpWbwqmV'
```