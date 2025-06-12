# go-blockchain: Submissão Acadêmica Distribuída com Blockchain
## Objetivo
Este projeto tem como objetivo demonstrar uma aplicação distribuída e paralela usando os conceitos de blockchain, adaptados ao contexto acadêmico (submissão de textos científicos). A ideia central é criar uma rede entre universidades/instituições para submissão segura e verificável de trabalhos acadêmicos (TCCs, artigos, relatórios, etc.), evitando plágio e garantindo autoria.

A aplicação foi construída como uma blockchain simples escrita em Go, com comunicação peer-to-peer, validação de blocos com PoW, e sincronização entre nós. Além disso, inclui uma interface web para que usuários possam submeter trabalhos diretamente de um navegador/intranet da universidade.

## Base Teórica: Artigo de Referência
O projeto inicialmente seria baseado no código do artigo:
```
An autonomous blockchain-based workflow execution broker for e-science
Alper Alimoğlu & Can Özturan, 2024
DOI: 10.1007/s10586-024-04534-z 
```
Link: https://doi.org/10.1007/s10586-024-04534-z

#### Sobre o artigo

Esse artigo propõe o uso de smart contracts em Ethereum para execução de workflows científicos com agendamento via SLURM.

Porém, a implementação original não funcionava por conta de dependências quebradas e complexidade de ambiente (Ethereum, IPFS, Slurm, DAG engine etc.).

Assim, decidimos construir nossa própria implementação baseada nos conceitos centrais do artigo e demais conceitos da área de blockchain, adaptando para um ambiente acadêmico:

- Submissão descentralizada e auditável

- Coordenação de ações via blockchain

- Validação de dados com hashes

- Propagação e consistência entre nós

## Conceitos Aplicados e Sobre nossa Aplicação
As submissões de trabalhos acadêmicos são registradas como blocos encadeados, garantindo integridade, autoria e auditabilidade.

Cada arquivo submetido tem seu hash SHA-256 calculado, funcionando como uma identificação única — o que impede submissões duplicadas e ajuda a combater o plágio.
A adição de novos blocos exige a realização de uma prova de trabalho (PoW), assegurando que apenas blocos válidos sejam inseridos na cadeia.

A comunicação entre os nós ocorre via conexões peer-to-peer (P2P), garantindo a propagação e sincronização dos dados.

Para facilitar o uso, a aplicação inclui uma interface web que permite o envio de arquivos diretamente do navegador, automatizando o cálculo do hash e a submissão para a blockchain.



## Funcionamento e Dados Usados
Cada trabalho submetido tem seu hash SHA-256 calculado no navegador. Esse hash é então enviado para o nó da instituição via requisição HTTP, que:

- Gera um novo bloco com a submissão

- Valida a unicidade do hash (evita plágios)

- Propaga o novo bloco aos peers

- Submissões repetidas (mesmo hash) são rejeitadas!

### Como rodar:

Compile o código com go build para gerar o executável.
``` bash
go build -o blockchain main.go
```

Inicie cada nó passando a porta como argumento, por exemplo:
``` bash
./blockchain 8000
```

Para iniciar um nó que já conhece outro par, passe o endereço do nó conhecido como segundo argumento, por exemplo:
``` bash
./blockchain 8001 localhost:8000
````
### Como funciona:
Exemplos de uso com curl:

Para enviar uma nova transação (JSON no corpo da requisição) a um nó na porta 8000:
``` bash
curl -X POST http://localhost:8000/add-data -d "Alice -> Bob: 10"
```

Para ver a cadiea de de blocos de um nó:
``` bash
curl http://localhost:8001/chain
```

## Autores
- Marlon Henrique Sanches

- Pedro Adorno Possebon

- Gabriel Lucchetta Garcia Sanchez

## Observações Finais
Este projeto é uma simplificação prática das ideias do artigo “eBlocBroker” para a disciplina de PDE.