package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/oanatmaria/ethblkcn-observer/parser"
)

type HttpServer struct {
	parser parser.Parser
	addr   string
	server *http.Server
}

func NewHttpServer(addr string, parser parser.Parser) Server {
	return &HttpServer{
		parser: parser,
		addr:   addr,
	}
}

func (s *HttpServer) Start(ctx context.Context) error {
	go s.startBlockProcessing(ctx)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /subscribe", s.wrapHandler(s.handleSubscribe))
	mux.HandleFunc("GET /transactions", s.wrapHandler(s.handleTransactions))
	mux.HandleFunc("GET /current_block", s.wrapHandler(s.handleCurrentBlock))

	s.server = &http.Server{
		Addr:    s.addr,
		Handler: mux,
	}

	go func() {
		<-ctx.Done()
		log.Println("Shutting down server...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := s.server.Shutdown(shutdownCtx); err != nil {
			log.Printf("Error shutting down server: %v", err)
		}
	}()

	log.Printf("Server is running at %s\n", s.addr)
	return s.server.ListenAndServe()
}

func (s *HttpServer) startBlockProcessing(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Stopping block processing...")
			return
		case <-ticker.C:
			log.Println("Processing blocks...")
			s.parser.ProcessNewBlocks(ctx)
			log.Println("Done processing latest blocks...")
		}
	}
}

func (s *HttpServer) wrapHandler(handler func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("Panic recovered: %v", rec)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()

		if err := handler(w, r); err != nil {
			log.Printf("Error: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (s *HttpServer) handleSubscribe(w http.ResponseWriter, r *http.Request) error {
	address := r.URL.Query().Get("address")
	if address == "" {
		http.Error(w, "Missing address parameter", http.StatusBadRequest)
		return nil
	}

	if !isValidEthAddress(address) {
		http.Error(w, "Invalid Ethereum address", http.StatusBadRequest)
		return nil
	}

	subscribed := s.parser.Subscribe(address)
	if subscribed {
		w.WriteHeader(http.StatusOK)
		if _, err := fmt.Fprintf(w, "Subscribed to address: %s\n", address); err != nil {
			log.Printf("Error writing response: %v", err)
		}
	} else {
		http.Error(w, fmt.Sprintf("Address already subscribed: %s", address), http.StatusBadRequest)
	}
	return nil
}

func (s *HttpServer) handleTransactions(w http.ResponseWriter, r *http.Request) error {
	address := r.URL.Query().Get("address")
	if address == "" {
		http.Error(w, "Missing address parameter", http.StatusBadRequest)
		return nil
	}

	transactions := s.parser.GetTransactions(address)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(transactions)
}

func (s *HttpServer) handleCurrentBlock(w http.ResponseWriter, r *http.Request) error {
	currentBlock := s.parser.GetCurrentBlock()
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(currentBlock)
}

func isValidEthAddress(address string) bool {
	regex := `^0x[0-9a-fA-F]{40}$`
	matched, _ := regexp.MatchString(regex, address)
	return matched
}
