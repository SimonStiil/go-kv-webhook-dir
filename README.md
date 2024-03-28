# Key Value Webhhook from Directory
Have you ever been in the situation that you need to share a secret og part of a secret between namespaces in kubernetes?
Do you already use External Secrets?
Then this may be the solution just for you.

## Insecure
This is a quite insecure solution as the whole idea is sharing "cluster public" information between namespaces like the CA Part of a secret or similar. If you need to add security consider adding an issue for what you need

## Configuration 
Base configuration is done with Environment to make it Kubernetes friendly
| Variable | Description | Default Value |
|---|---|---|
| PORT | Port number for the webhook | 8080 |
| KEY_DIRECTORY | directory on host to find keys (filenames) in | /keys |

All files in KEY_DIRECTORY will be made available through http requests as a json response 
## Example usage
When run bash in cluster:
```bash
curl http://secret-sharing.example:8080/hello
```
Returns reply:
```json
{"key":"hello","value":"world"}
```

Based on examples below
### Deployment configuration
```yaml
container: 
  - name: go-kv-webhook-dir
    image: ghcr.io/simonstiil/go-kv-webhook-dir:main
    ports:
      - containerPort: 8080
    volumeMounts:
      - name: hello
        mountPath: /keys/hello
        subPath: hello
        readOnly: true
    resources: # Image has a very small footprint ~0.00003% CPU and < 1.91Mb Ram based on deployment on RK1 ARM64 
      limits:
        cpu: 10m
        memory: 50Mi
      requests:
        cpu: 1m
        memory: 5Mi
    ...
volumes:
  - name: hello
    secret:
      secretName: hello-world
      items:
        - key: somekeyname
          path: hello
```

### External Secrets Configuration
If above has a service named `secret-sharing` and deployed in in namespace `example`

```yaml
apiVersion: external-secrets.io/v1beta1
kind: ClusterSecretStore
metadata:
  name: secret-sharing-example
spec:
  provider:
    webhook:
      url: "http://secret-sharing.example:8080/{{ .remoteRef.key }}" # file name to read from /keys/
      result:
        jsonPath: "$.value"
      headers:
        Content-Type: application/json
```

```yaml
apiVersion: external-secrets.io/v1beta1
kind: ClusterExternalSecret
metadata:
  name: secret-sharing-hello-world
spec:
  namespaceSelector:
    matchLabels:
      secret-sharing/hello: "true" # label to look for in namespaces to populate secret
  refreshTime: 1h
  externalSecretSpec:
    secretStoreRef:
      name: secret-sharing-example
      kind: ClusterSecretStore
    refreshInterval: 1h
    target:
      name: shared-hello-secret # Target secret name
      template:
        type: Opaque
        data:
          somekeyname: "{{ .hello | toString }}"
    data:
      - secretKey: hello
        remoteRef:
          key: hello # file name to read from /keys/
```

## Multi arch support
This is the first image i have made that supports both ARM64 add AMD64 using kaniko and manifest-tool