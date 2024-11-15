package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/oanatmaria/ethblkcn-observer/parser"
	"github.com/oanatmaria/ethblkcn-observer/storage"
)

func TestHandleSubscribe(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockParser := parser.NewMockParser(ctrl)
	srv := NewHttpServer(":8080", mockParser)

	tests := []struct {
		name           string
		address        string
		subscribeResp  bool
		expectCall     bool
		expectedStatus int
	}{
		{"ValidAddressSubscribed", "0x1234567890abcdef1234567890abcdef12345678", true, true, http.StatusOK},
		{"ValidAddressAlreadySubscribed", "0x1234567890abcdef1234567890abcdef12345678", false, true, http.StatusBadRequest},
		{"InvalidAddress", "invalidAddress", false, false, http.StatusBadRequest},
		{"EmptyAddress", "", false, false, http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectCall {
				mockParser.EXPECT().Subscribe(tt.address).Return(tt.subscribeResp)
			}

			req := httptest.NewRequest("POST", "/subscribe?address="+tt.address, nil)
			w := httptest.NewRecorder()

			err := srv.(*HttpServer).handleSubscribe(w, req)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			resp := w.Result()
			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}
		})
	}
}

func TestHandleTransactions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockParser := parser.NewMockParser(ctrl)
	srv := NewHttpServer(":8080", mockParser)

	tests := []struct {
		name           string
		address        string
		mockResponse   []storage.Transaction
		expectCall     bool
		expectedStatus int
	}{
		{"ValidAddressWithTransactions", "0x1234567890abcdef1234567890abcdef12345678", []storage.Transaction{}, true, http.StatusOK},
		{"MissingAddress", "", nil, false, http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectCall {
				mockParser.EXPECT().GetTransactions(tt.address).Return(tt.mockResponse)
			}

			req := httptest.NewRequest("GET", "/transactions?address="+tt.address, nil)
			w := httptest.NewRecorder()

			err := srv.(*HttpServer).handleTransactions(w, req)
			if err != nil && tt.expectedStatus != http.StatusInternalServerError {
				t.Fatalf("Unexpected error: %v", err)
			}

			resp := w.Result()
			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}
		})
	}
}

func TestHandleCurrentBlock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockParser := parser.NewMockParser(ctrl)
	srv := NewHttpServer(":8080", mockParser)

	mockParser.EXPECT().GetCurrentBlock().Return(123456)

	req := httptest.NewRequest("GET", "/current_block", nil)
	w := httptest.NewRecorder()

	err := srv.(*HttpServer).handleCurrentBlock(w, req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

func TestStartServerAndShutdown(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockParser := parser.NewMockParser(ctrl)
	srv := NewHttpServer(":8080", mockParser)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		// Shutdown server after it starts
		time.Sleep(2 * time.Second)
		cancel()
	}()

	err := srv.Start(ctx)
	if err != nil && err != http.ErrServerClosed {
		t.Fatalf("Unexpected error: %v", err)
	}
}
