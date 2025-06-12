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

#### Sobre o artigo e Correlação

Esse artigo propõe o uso de smart contracts em Ethereum para execução de workflows científicos com agendamento via SLURM.

Porém, a implementação original não funcionava por conta de dependências quebradas e complexidade de ambiente (Ethereum, IPFS, Slurm, DAG engine etc.).

Assim, decidimos construir nossa própria implementação baseada nos conceitos centrais do artigo e demais conceitos da área de blockchain:

- **Utilização de Blockchain para Garantia de Integridade e Autoria:**
Assim como o artigo propõe o uso de smart contracts em Ethereum para gerenciar workflows científicos e garantir sua execução confiável, o nosso sistema usa blockchain para registrar submissões acadêmicas, assegurando a integridade, autoria e auditabilidade dos trabalhos submetidos. A hash SHA-256 de cada arquivo é um mecanismo semelhante ao hash dos outputs de jobs no artigo, garantindo que cada submissão seja única e resistente a duplicações ou plágios. (Páginas 3, 7 e 9)

- **Controle Descentralizado e Segurança na Submissão:** O artigo destaca a coordenação automática e descentralizada do processamento de workflows via smart contracts, eliminando a dependência de uma entidade central. Da mesma forma, nosso sistema distribui o registro de submissões por meio de uma rede P2P, validando e propagando blocos com trabalhos acadêmicos, garantindo transparência e resistência a manipulações. (Páginas 9 e 7)
  
- **Prova de Trabalho (PoW) para Validação dos Blocos:**
No artigo, blocos só são adicionados após a realização de PoW, garantindo validação da autoria e resistência a ataques. Nossa implementação incorpora PoW ao validar submissões antes de adicioná-las ao blockchain, reforçando a segurança e a integridade do sistema de submissão acadêmica.

- **Propagação, Sincronização e Confiança entre Nós:**
A comunicação P2P do nosso sistema para sincronizar blocos reflete o mecanismo de propagação de blocos e validações do artigo, que usa a rede Ethereum e IPFS para compartilhar dados de workflows (Página 3). Porém, na nossa rede os nós por si só garantem que todas as universidades/instituições tenham uma visão consistente do histórico de submissões.

- **Resistência a Plágio e Submissões Duplication:**
A verificação do hash e a rejeição de submissões repetidas em seu sistema espelham o mecanismo do artigo de impedir duplicidade de outputs ou jobs, promovendo uma plataforma confiável e imutável para registros acadêmicos. 

- **Automação e Interface Web:**
O artigo propõe a automação do registro e validação usando smart contracts e interfaces de usuário (via IPFS), e o nosso projeto projeto também incorpora uma interface web que automatiza a submissão, alinhando-se à ideia de uma plataforma automatizada, segura e transparente para submissões acadêmicas como aborado no artigo. (Página 9)

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