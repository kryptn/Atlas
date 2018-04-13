# Atlas

This microservice collects annotations on kubernetes services within the configured namespace

###Initial Goal

to give the client an endpoint that would return a list of graphql schemas/endpoints

the k8s service for each graphql microservice will need to be updated to add a custom annotation:

```yaml
kind: Service
apiVersion: v1
metadata:
  name: some-service
  annotations:
    atlas.domain.io/graph.client: "/graphql"
    atlas.domain.io/graph.admin: "/admin/graphql"
``` 

When all graphql services are annotated the client is now able to collect a list of each service and its graphql endpoint using `GET /full/atlas.domain.io/graph/client`. Examples below

host endpoint should be constructable using `'http://'+item.service_name+item.annotation.value`

## Special Setup?

I had to run this to create the serviceaccount, role, and rolebinding:

`kubectl create clusterrolebinding user-cluster-admin-binding --clusterrole=cluster-admin --user=user@domain.com`

## Details

Annotations in form `fq.dn/key.value.value: string_value` will be collected

valid keys:

    atlas.domain.io/key.value
    domain.io/key.value.is.continued
    
invalid keys:
    
    anything.without.a.slash
    name/only-a-key

This is the regex for the annotation key parser: `([a-zA-Z][-a-zA-Z0-9.]*)\/([-a-zA-Z0-9]+)\.([-a-zA-Z0-9.]*)`

Example services:
```yaml
kind: Service
apiVersion: v1
metadata:
  name: example-service-one
  annotations:
    atlas.domain.io/graph.client: "/graphql"
    atlas.domain.io/graph.admin: "/admin/graphql"
    atlas.ext.io/unrelated-key.unrelated-value: "nonsense"
---
kind: Service
apiVersion: v1
metadata:
  name: example-service-two
  annotations:
    atlas.domain.io/graph.client: "/graphql"
    atlas.domain.io/graph.admin: "/admin/graphql"
    atlas.domain.io/graph.external: "/api/v1/graph"
---
kind: Service
apiVersion: v1
metadata:
  name: example-service-three
  annotations:
    atlas.domain.io/graph.client: "/graphql"
    atlas.domain.io/graph.admin: "/admin/graphql"
    atlas.ext.io/graph.client: "/ext/graphql"
```

## Usage

#### Full usage with FQDN and key/value


```http request
GET /full/atlas.io/graph/client
```

```json
{
  "annotations": [
    {
      "key": "graph",
      "value": "client",
      "annotation": {
        "key": "atlas.domain.io/graph.client",
        "value": "/graphql"
      },
      "service_name": "atlas-service"
    },
    {
      "key": "graph",
      "value": "client",
      "annotation": {
        "key": "atlas.domain.io/graph.client",
        "value": "/graphql"
      },
      "service_name": "example-service-one"
    },
    {
      "key": "graph",
      "value": "client",
      "annotation": {
        "key": "atlas.domain.io/graph.client",
        "value": "/graphql"
      },
      "service_name": "example-service-three"
    },
    {
      "key": "graph",
      "value": "client",
      "annotation": {
        "key": "atlas.domain.io/graph.client",
        "value": "/graphql"
      },
      "service_name": "example-service-two"
    }
  ]
}
```

#### Get all using just FQDN and key

```http request
GET /full/atlas.domain.io/graph/
```

```json
{
  "annotations": [
    {
      "key": "graph",
      "value": "client",
      "annotation": {
        "key": "atlas.ext.io/graph.client",
        "value": "/ext/graphql"
      },
      "service_name": "example-service-three"
    }
  ]
}
```

#### Get only using the FQDN

```http request
GET /full/atlas.ext.io
```

```json

{
  "annotations": [
    {
      "key": "unrelated-key",
      "value": "unrelated-value",
      "annotation": {
        "key": "atlas.ext.io/unrelated-key.unrelated-value",
        "value": "nonsense"
      },
      "service_name": "example-service-one"
    },
    {
      "key": "graph",
      "value": "client",
      "annotation": {
        "key": "atlas.ext.io/graph.client",
        "value": "/ext/graphql"
      },
      "service_name": "example-service-three"
    }
  ]
}
```