package db

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"rockstake-core/types"
	"time"

	shell "github.com/ipfs/go-ipfs-api"
)

type NFTStore interface {
	GenerateAndUploadNFT(ctx context.Context, betData types.NftNodeInfo) (string, error)
}

type FilebaseStore struct {
	accessKey  string
	secretKey  string
	bucketName string
}

func NewFilebaseStore(accessKey, secretKey, bucketName string) *FilebaseStore {
	return &FilebaseStore{
		accessKey:  accessKey,
		secretKey:  secretKey,
		bucketName: bucketName,
	}
}

func uploadToIPFS(filename string) (string, error) {
	// Conectăm clientul IPFS la nodul local
	sh := shell.NewShell("localhost:5001")

	// Citim conținutul fișierului
	file, err := os.Open(filename)
	if err != nil {
		return "", fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	fileData, err := ioutil.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("error reading file: %v", err)
	}

	// Adăugăm fișierul la IPFS
	cid, err := sh.Add(bytes.NewReader(fileData))
	if err != nil {
		return "", fmt.Errorf("error uploading to IPFS: %v", err)
	}

	// Asigurăm propagarea globală a CID-ului
	err = sh.Pin(cid) // Pinăm CID-ul local pentru a ne asigura că rămâne disponibil
	if err != nil {
		return "", fmt.Errorf("error pinning CID: %v", err)
	}

	// Publicăm CID-ul în rețea (DHT - Distributed Hash Table)
	err = sh.Pin(cid)
	if err != nil {
		return "", fmt.Errorf("error providing CID to network: %v", err)
	}

	// Construim URL-ul
	ipfsURL := fmt.Sprintf("https://ipfs.io/ipfs/%s", cid)

	return ipfsURL, nil
}

func (s *FilebaseStore) GenerateAndUploadNFT(ctx context.Context, betData types.NftNodeInfo) (string, error) {
	metadata := types.NFTMetadata{
		Description: "RS-Exchange",
		Compiler:    "Rockstake",
		Attributes: []types.Attribute{
			{TraitType: "Fixture", Value: betData.FixtureID},
			{TraitType: "Market", Value: betData.Market},
			{TraitType: "Selection", Value: betData.Selection},
			{TraitType: "Type", Value: betData.Type},
			{TraitType: "Odd", Value: fmt.Sprintf("%.2f", betData.Odd)},
			{TraitType: "Stake", Value: fmt.Sprintf("%.2f", betData.Stake)},
			{TraitType: "Potential Win", Value: fmt.Sprintf("%.2f", betData.Odd*betData.Stake)},
			{TraitType: "Status", Value: "open"},
		},
	}

	jsonData, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error marshaling metadata: %v", err)
	}

	tempFile := fmt.Sprintf("bet_metadata_%d.json", time.Now().UnixNano())
	if err := ioutil.WriteFile(tempFile, jsonData, 0644); err != nil {
		return "", fmt.Errorf("error saving temp file: %v", err)
	}
	defer os.Remove(tempFile)

	cid, err := uploadToIPFS(tempFile)
	if err != nil {
		return "", fmt.Errorf("error uploading to IPFS: %v", err)
	}

	return cid, nil
}
