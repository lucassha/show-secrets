## show-secrets

Show all decoded secrets in a Kubernetes namespace.

Kubernetes secrets are base64 encoded, and as such, generally require some manipulation with `jq` in order to properly parse what the secret values are. This CLI allows for a quick check of the decoded values for each secret.

### Usage

`show-secrets` can be used as a standalone binary, or as a [kubectl plugin](addlink).

Each of these commands are equivalent:

```bash
show-secrets
kubectl show secrets
```

In order to see all flags: `show-secrets --help`

Generally speaking, `show-secrets` implements flags with the same pattern as `kubectl`. As of now, output, namespace, and kubeconfig flags are enabled.

```bash
# show secrets for devs namespace
kubectl show secrets -n devs

# output json instead of wide output
kubectl show secrets -o json

# use a different kubeconfig file
kubectl show secrets --kubeconfig ~/some/file/here
```

### Sample Workflow

Create some secrets to test against 

```bash
kubectl create secret generic secret3 --from-literal=name=sam --from-literal=age=30 --from-literal=town=columbia,sc
kubectl create secret generic secret2 --from-literal=name=caitlyn --from-literal=age=28
kubectl create secret generic secret1 --from-literal=name=shannon
```

Compare the two outputs:

`kubectl get secrets -n default`

```bash
NAME                  TYPE                                  DATA   AGE
default-token-6cr5c   kubernetes.io/service-account-token   3      8d
secret1               Opaque                                1      8d
secret2               Opaque                                2      8d
secret3               Opaque                                3      8d
```

`kubectl show secrets -n default`

```bash
NAME           KEY         VALUE
secret1        name        shannon
secret2        age         28
└──            name        caitlyn
secret3        age         30
├──            name        sam
└──            town        columbia,sc
```

## Installation

Currently only implemented for Mac. Creates two binaries in `/usr/local/bin/`.

`brew install lucassha/homebrew/show-secrets`

## TODO

- Finish the script to push a tarball up to S3 for homebrew downloading
- Add support for more than Mac
- Update code to include regex search for secrets (i.e. -- `kubectl show secrets AWS_*`)