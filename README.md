# go-blockchain

O objetivo desse projeto é implementar uma blockchain simples em Go com comunicação peer-to-peer entre vários nós.
Cada nó pode receber blocos, validar hashes e sincronizar com outros nós da rede (existe POW e validação de hashes), permitindo um sistema distribuído básico.

# Como rodar:

Compile o código com go build para gerar o executável.

Inicie cada nó passando a porta como argumento, por exemplo:
``` bash
./blockchain 8000
```

Para iniciar um nó que já conhece outro par, passe o endereço do nó conhecido como segundo argumento, por exemplo:
``` bash
./blockchain 8001 localhost:8000
````
# Como funciona:
Exemplos de uso com curl:

Para enviar uma nova transação (JSON no corpo da requisição) a um nó na porta 8000:
``` bash
curl -X POST http://localhost:8000/add-data -d "Alice -> Bob: 10"
```

Para ver a cadiea de de blocos de um nó:
``` bash
curl http://localhost:8001/chain
```