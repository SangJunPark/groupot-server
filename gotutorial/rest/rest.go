package rest

import (
	"encoding/json"
	"fmt"
	"gotutorial/blockchain"
	"gotutorial/p2p"
	"gotutorial/utils"
	"gotutorial/wallet"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var port string

type url string

func (u url) MarshalText() ([]byte, error) {
	url := fmt.Sprintf("http://localhost%s%s", port, u)
	return []byte(url), nil
}

type URLDescription struct {
	URL         url    `json:"url"`
	Method      string `json:"method"`
	Description string `json:"description"`
	Payload     string `json:"payload,omitempty"`
	ignore      string `json:"-"`
}

type addBlock struct {
	Message string
}

type addPeerPayload struct {
	Address string `json:"address"`
	Port    string `json:"port"`
}

type myWalletResponse struct {
	Address string `json:"address"`
}

type balanceResponse struct {
	Address string `json:"address"`
	Balance int    `json:"balance"`
}

type addTxPayload struct {
	To     string `json:"to"`
	Amount int    `json:"amount"`
}

type errorResponse struct {
	ErrorMessage string `json:"errorMessage"`
}

//ref Stringer interface it defines inside of GO fmt
func (u URLDescription) String() string {
	return "fuck you guys"
}

func documentation(rw http.ResponseWriter, r *http.Request) {
	data := []URLDescription{
		{
			URL:         url("/"),
			Method:      "GET",
			Description: "See documentation",
		},
		{
			URL:         url("/blocks"),
			Method:      "POST",
			Description: "",
			Payload:     "data:string",
		},
		{
			URL:         url("/status"),
			Method:      "GET",
			Description: "status",
			Payload:     "data:string",
		},
		{
			URL:         url("/blocks/{hash}"),
			Method:      "POST",
			Description: "",
			Payload:     "data:string",
		},
		{
			URL:         url("/balance/{address}"),
			Method:      "POST",
			Description: "",
			Payload:     "data:string",
		},
		{
			URL:         url("/ws"),
			Method:      "GET",
			Description: "upgrade to websocket",
		},
		{
			URL:         url("/peers"),
			Method:      "GET, POST",
			Description: "peers",
			Payload:     "data:string",
		},
	}

	//rw.Header().Add("Content-Type", "application/json")
	// b, err := json.Marshal(data)
	// utils.HandleErr(err)
	// fmt.Fprintf(rw, "%s", b)
	json.NewEncoder(rw).Encode(data)
}

func blocks(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		json.NewEncoder(rw).Encode(blockchain.Blocks(blockchain.Blockchain()))
	case "POST":
		var addBlock addBlock
		utils.HandleErr(json.NewDecoder(r.Body).Decode(&addBlock))
		fmt.Println(addBlock)
		block := blockchain.Blockchain().AddBlock()

		p2p.BroadcastNewBlock(block)

		//json.NewEncoder(rw).Encode(blockchain.GetBlockchain().AllBlocks())
		rw.WriteHeader(http.StatusCreated)
		http.Redirect(rw, r, "/", http.StatusMovedPermanently)
	}
}

func block(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Println(vars)

	hash := vars["hash"]
	block, err := blockchain.FindBlock(hash)
	encoder := json.NewEncoder(rw)

	if err == nil {
		encoder.Encode(block)
	} else {
		encoder.Encode(errorResponse{fmt.Sprint(err)})
	}
}

func status(rw http.ResponseWriter, r *http.Request) {
	json.NewEncoder(rw).Encode(blockchain.Status(blockchain.Blockchain()))
}

func balance(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]
	total := r.URL.Query().Get("total")
	switch total {
	case "true":
		amount := blockchain.BlanceByAddress(blockchain.Blockchain(), address)
		utils.HandleErr(json.NewEncoder(rw).Encode(balanceResponse{
			Address: address,
			Balance: amount,
		}))
	default:
		txOuts := blockchain.UTxOutsByAddress(blockchain.Blockchain(), address)
		utils.HandleErr(json.NewEncoder(rw).Encode(txOuts))
	}
}

func mempool(rw http.ResponseWriter, r *http.Request) {
	utils.HandleErr(json.NewEncoder(rw).Encode(blockchain.Mempool().Transactions()))
}

//datarace - 두개이상의 goroutine이 같은 데이터에 접근 했을때 발생
func transactions(rw http.ResponseWriter, r *http.Request) {
	var txPayload addTxPayload
	utils.HandleErr(json.NewDecoder(r.Body).Decode(&txPayload))
	tx, err := blockchain.Mempool().AddTx(txPayload.To, txPayload.Amount)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(errorResponse{err.Error()})
		return
	}
	p2p.BroadcastNewTransaction(tx)
	rw.WriteHeader(http.StatusCreated)
}

func peers(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		json.NewEncoder(rw).Encode(p2p.AllPeers(&p2p.Peers))
	case "POST":
		var payload addPeerPayload
		json.NewDecoder(r.Body).Decode(&payload)
		p2p.AddPeer(payload.Address, payload.Port, fmt.Sprint(port), true)
	}
}

//adapter pattern
func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(rw, r)
	})
}

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL)
		next.ServeHTTP(rw, r)
	})
}

func myWallet(rw http.ResponseWriter, r *http.Request) {
	address := wallet.Wallet().Address
	json.NewEncoder(rw).Encode(myWalletResponse{address})

	//it can be but looks dirty
	// json.NewEncoder(rw).Encode(struct {
	// 	Address string `json:"address"`
	// }{Address: address})
}

func Start(iPort int) {
	router := mux.NewRouter()
	port = fmt.Sprintf(":%d", iPort)
	router.Use(jsonContentTypeMiddleware)
	router.Use(loggerMiddleware)

	router.HandleFunc("/", documentation).Methods("GET")
	router.HandleFunc("/blocks", blocks).Methods("GET", "POST")
	router.HandleFunc("/status", status).Methods("GET")
	//hexadecimal find so useful
	router.HandleFunc("/blocks/{hash:[a-f0-9]+}", block).Methods("GET")
	router.HandleFunc("/balance/{address}", balance).Methods("GET")
	router.HandleFunc("/mempool", mempool).Methods("GET")
	router.HandleFunc("/wallet", myWallet).Methods("GET")
	router.HandleFunc("/transactions", transactions).Methods("POST")
	router.HandleFunc("/ws", p2p.Upgrade).Methods("GET")
	router.HandleFunc("/peers", peers).Methods("GET", "POST")

	fmt.Printf("%s port\n", port)
	log.Fatal(http.ListenAndServe(port, router))

}
