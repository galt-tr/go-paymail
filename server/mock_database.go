package server

import (
	"fmt"
	"log"
	"strings"

	"github.com/bitcoinschema/go-bitcoin"
	"github.com/mrz1836/go-logger"
)

// paymailAddress is a mock table of paymail addresses for the example server
type paymailAddress struct {
	Alias       string `json:"alias"`        // Alias or handle of the paymail
	ID          uint64 `json:"id"`           // Unique identifier
	PrivateKey  string `json:"private_key"`  // PrivateKey hex encoded
	PubKey      string `json:"pubkey"`       // PublicKey hex encoded
	LastAddress string `json:"last_address"` // This is used as a temp address for now (should be via xPub)
}

// paymailAddressTable is the mocked data for the example server (table: paymail_address)
var paymailAddressTable []*paymailAddress

// Create the list of mock aliases to create on load
var mockAliases = []struct {
	alias string
	id    uint64
}{
	{"mrz", 1},
	{"satchmo", 2},
}

// init run on load
func init() {

	// Generate a paymail addresses
	for _, mock := range mockAliases {
		if err := generatePaymail(mock.alias, mock.id); err != nil {
			log.Fatalf("failed to create paymail address in mock database for alias: %s id: %d", mock.alias, mock.id)
		}
	}

	// Log the paymail addresses in database
	logger.Data(2, logger.DEBUG, fmt.Sprintf("found %d paymails in the mock database", len(mockAliases)))
}

// generatePaymail will make a new row in the mock database
// creates a private key and pubkey
func generatePaymail(alias string, id uint64) error {

	// Start a row
	row := &paymailAddress{ID: id, Alias: alias}

	var err error

	// Generate new private key
	if row.PrivateKey, err = bitcoin.CreatePrivateKeyString(); err != nil {
		return err
	}

	// Get address
	if row.LastAddress, err = bitcoin.AddressFromPrivateKey(row.PrivateKey); err != nil {
		return err
	}

	// Derive a pubkey from private key
	if row.PubKey, err = bitcoin.PubKeyFromPrivateKey(row.PrivateKey); err != nil {
		return err
	}

	// Load some mock paymail addresses
	paymailAddressTable = append(paymailAddressTable, row)

	return nil
}

// getPaymail will find a paymail address given an alias
func getPaymailByAlias(alias string) *paymailAddress {
	alias = strings.ToLower(alias)
	for i, row := range paymailAddressTable {
		if alias == row.Alias {
			return paymailAddressTable[i]
		}
	}
	return nil
}
