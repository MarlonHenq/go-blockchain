package main

import (
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"
    "os"
    "strconv"
    "strings"
    "time"
)

type Block struct {
    Index      int
    Timestamp  string
    Data       string
    PrevHash   string
    Hash       string
    Nonce      int
    Difficulty int
}

var Blockchain []Block
var Peers []string
var Difficulty = 3

// pow
func calculateHash(block Block) string {
    record := strconv.Itoa(block.Index) + block.Timestamp + block.Data + block.PrevHash + strconv.Itoa(block.Nonce)
    h := sha256.Sum256([]byte(record))
    return hex.EncodeToString(h[:])
}

func mineBlock(block Block) Block {
    prefix := strings.Repeat("0", block.Difficulty)
    for {
        hash := calculateHash(block)
        if strings.HasPrefix(hash, prefix) {
            block.Hash = hash
            return block
        }
        block.Nonce++
    }
}

func createGenesisBlock() Block {
    genesis := Block{
        Index:      0,
        Timestamp:  "2026-06-07T00:00:00Z",
        Data:       "Genesis",
        PrevHash:   "",
        Nonce:      2272,
        Difficulty: Difficulty,
    }
    genesis.Hash = calculateHash(genesis)
    return genesis
}


func generateBlock(oldBlock Block, data string) Block {
    newBlock := Block{
        Index:      oldBlock.Index + 1,
        Timestamp:  time.Now().String(),
        Data:       data,
        PrevHash:   oldBlock.Hash,
        Difficulty: Difficulty,
    }
    return mineBlock(newBlock)
}

func isBlockValid(newBlock, oldBlock Block) bool {
    if oldBlock.Index+1 != newBlock.Index || oldBlock.Hash != newBlock.PrevHash || calculateHash(newBlock) != newBlock.Hash {
        return false
    }
    return true
}

// HTTP Handlers
func getChain(w http.ResponseWriter, _ *http.Request) {
    json.NewEncoder(w).Encode(Blockchain)
}

func addData(w http.ResponseWriter, r *http.Request) {
    body, _ := io.ReadAll(r.Body)
    data := string(body)
    newBlock := generateBlock(Blockchain[len(Blockchain)-1], data)
    if isBlockValid(newBlock, Blockchain[len(Blockchain)-1]) {
        Blockchain = append(Blockchain, newBlock)
        propagateBlock(newBlock)
    }
    json.NewEncoder(w).Encode(newBlock)
}

func receiveBlock(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Erro ao ler o corpo", http.StatusBadRequest)
		return
	}

	var block Block
	err = json.Unmarshal(body, &block)
	if err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	lastBlock := Blockchain[len(Blockchain)-1]

	if isBlockValid(block, lastBlock) {
		Blockchain = append(Blockchain, block)
		fmt.Println("Bloco recebido e adicionado à blockchain:")
		fmt.Printf("Index: %d | Data: %s | Hash: %s\n", block.Index, block.Data, block.Hash)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Bloco adicionado com sucesso"))
	} else {
		fmt.Println("Bloco recebido é inválido!")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bloco inválido"))
	}
}


func registerNode(w http.ResponseWriter, r *http.Request) {
    body, _ := io.ReadAll(r.Body)
    peer := string(body)

    for _, p := range Peers {
        if p == peer {
            json.NewEncoder(w).Encode(Peers) // já registrado, mas retorna peers
            return
        }
    }

    Peers = append(Peers, peer)
    fmt.Println("Novo peer registrado:", peer)

    // retorna os peers atuais para o novo nó
    json.NewEncoder(w).Encode(Peers)
}



func propagateBlock(block Block) {
    for _, peer := range Peers {
        jsonData, _ := json.Marshal(block)
        resp, err := http.Post("http://"+peer+"/receive-block", "application/json", strings.NewReader(string(jsonData)))
        if err != nil {
            fmt.Println("Erro ao propagar para", peer, ":", err)
        } else {
            fmt.Println("Propagado para", peer, "Status:", resp.Status)
        }
    }
}

func isValidChain(chain []Block) bool {
    for i := 1; i < len(chain); i++ {
        if !isBlockValid(chain[i], chain[i-1]) {
            return false
        }
    }
    return true
}

func syncWithPeers() {
    for _, peer := range Peers {
        fmt.Println("Sincronizando com peer:", peer)
        resp, err := http.Get("http://" + peer + "/chain")
        if err != nil {
            fmt.Println("Erro ao puxar blockchain do peer", peer, ":", err)
            continue
        }
        defer resp.Body.Close()

        var receivedChain []Block
        if err := json.NewDecoder(resp.Body).Decode(&receivedChain); err != nil {
            fmt.Println("Erro ao decodificar chain:", err)
            continue
        }

        if isValidChain(receivedChain) && len(receivedChain) > len(Blockchain) {
            Blockchain = receivedChain
            fmt.Println("Blockchain substituída pela versão do peer", peer)
        }
    }
}

func registerWithPeer(peer string) {
    self := "localhost:" + os.Args[1]
    resp, err := http.Post("http://"+peer+"/register-node", "text/plain", strings.NewReader(self))
    if err != nil {
        fmt.Println("Erro ao registrar com peer:", err)
        return
    }
    defer resp.Body.Close()

    var knownPeers []string
    err = json.NewDecoder(resp.Body).Decode(&knownPeers)
    if err != nil {
        fmt.Println("Erro ao decodificar peers:", err)
        return
    }

    for _, p := range knownPeers {
        if p != self && !contains(Peers, p) {
            Peers = append(Peers, p)

            // REGISTRA-SE DE VOLTA NESSE PEER
            _, err := http.Post("http://"+p+"/register-node", "text/plain", strings.NewReader(self))
            if err != nil {
                fmt.Println("Erro ao se registrar com peer conhecido:", p, err)
            }
        }
    }

    if !contains(Peers, peer) {
        Peers = append(Peers, peer)
    }

    fmt.Println("Peers conhecidos após registro:", Peers)
}


func contains(slice []string, item string) bool {
    for _, v := range slice {
        if v == item {
            return true
        }
    }
    return false
}



func main() {
    port := "8000"
    if len(os.Args) > 1 {
        port = os.Args[1]
    }

    Blockchain = append(Blockchain, createGenesisBlock())

    http.HandleFunc("/chain", getChain)
    http.HandleFunc("/add-data", addData)
    http.HandleFunc("/receive-block", receiveBlock)
    http.HandleFunc("/register-node", registerNode)

    if len(os.Args) > 2 {
        registerWithPeer(os.Args[2]) // se passado, tenta registrar com outro nó
    }

    fmt.Println("Listening on port", port)
    log.Fatal(http.ListenAndServe(":"+port, nil))
}

