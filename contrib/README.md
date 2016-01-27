CONTRIBUTION
-----------------

# Running full dev environnement

```bash
# Running Registry, Clair and Hyperclair-DEV-BOX
docker-compose --x-networking up -d

# Enter the hyperclair dev box
docker exec -ti hyperclair bash

# Run Any command ex:
go run main.go help
# Or
go run main.go pull --config .hyperclair.yml registry:5000/jgsqware/ubuntu-git
```
