## show-secrets

TODO: ADD README :) 

### Development

Create some secrets to test against 

```bash
kubectl create secret generic secret3 --from-literal=name=sam --from-literal=age=30 --from-literal=town=columbia,sc
kubectl create secret generic secret2 --from-literal=name=caitlyn --from-literal=age=28
kubectl create secret generic secret1 --from-literal=name=shannon
```