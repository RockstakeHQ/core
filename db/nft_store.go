package db

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"rockstake-core/types"
	"time"
)

type NFTStore interface {
	GenerateAndUploadNFT(ctx context.Context, betData types.NftNodeInfo) (string, error)
}

type PinataNFTStore struct {
	pinataKey    string
	pinataSecret string
}

func NewPinataNFTStore(pinataKey, pinataSecret string) *PinataNFTStore {
	return &PinataNFTStore{
		pinataKey:    pinataKey,
		pinataSecret: pinataSecret,
	}
}

func uploadToPinata(filename string, pinataAPIKey string, pinataAPISecret string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	fileContent, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return "", err
	}
	part.Write(fileContent)
	writer.Close()

	req, err := http.NewRequest("POST", "https://api.pinata.cloud/pinning/pinFileToIPFS", &requestBody)
	if err != nil {
		return "", err
	}
	req.Header.Set("pinata_api_key", pinataAPIKey)
	req.Header.Set("pinata_secret_api_key", pinataAPISecret)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to upload file to Pinata: %s", resp.Status)
	}

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	cid, ok := result["IpfsHash"].(string)
	if !ok {
		return "", fmt.Errorf("failed to get CID from Pinata response")
	}

	return cid, nil
}

func (s *PinataNFTStore) GenerateAndUploadNFT(ctx context.Context, betData types.NftNodeInfo) (string, error) {
	metadata := types.NFTMetadata{
		Description: "Bet Exchange",
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

	cid, err := uploadToPinata(tempFile, s.pinataKey, s.pinataSecret)
	if err != nil {
		return "", fmt.Errorf("error uploading to IPFS: %v", err)
	}

	return cid, nil
}
