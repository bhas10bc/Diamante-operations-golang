package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/diamcircle/go/clients/auroraclient"
	"github.com/diamcircle/go/keypair"
	"github.com/diamcircle/go/network"
	"github.com/diamcircle/go/txnbuild"
)

// Payment creates and submits a payment operation on the Diamante network
func Payment(destination, amount string) error {
	networkPassphrase := network.TestNetworkPassphrase
	client := auroraclient.DefaultTestNetClient

	// Parse the secret key to initialize sourceKeypair
	sourceKeypair := keypair.MustParseFull("SBLTFVCSTVKJEH2WXY5HFXHOPEF3EDBYUK7X6I3MFDXC3PJAJX2L3BPE") // Replace with actual admin secret key

	// Fetch account details for sourceKeypair's address to initialize sourceAccount
	sourceAccount, err := client.AccountDetail(auroraclient.AccountRequest{AccountID: sourceKeypair.Address()})
	if err != nil {
		return err
	}

	// Initialize payment as a txnbuild.Payment struct
	payment := txnbuild.Payment{
		Destination: destination,
		Amount:      amount,
		Asset:       txnbuild.NativeAsset{},
	}

	// Initialize txParams as a txnbuild.TransactionParams struct
	txParams := txnbuild.TransactionParams{
		SourceAccount:        &sourceAccount,
		IncrementSequenceNum: true,
		Operations:           []txnbuild.Operation{&payment},
		BaseFee:              txnbuild.MinBaseFee,
		Timebounds:           txnbuild.NewInfiniteTimeout(),
	}

	// Initialize tx as a txnbuild.Transaction struct
	tx, err := txnbuild.NewTransaction(txParams)
	if err != nil {
		return err
	}

	// Sign tx
	tx, err = tx.Sign(networkPassphrase, sourceKeypair)
	if err != nil {
		return err
	}

	// Submit tx to the Diamante network
	_, err = client.SubmitTransaction(tx)
	if err != nil {
		return err
	}

	return nil
}

// ManageData sets, updates, or deletes a data entry for an account
func ManageData(name string, value string) error {
	// Initialize network passphrase and client
	networkPassphrase := network.TestNetworkPassphrase
	client := auroraclient.DefaultTestNetClient

	// Parse secret key of source account to obtain sourceKeypair
	sourceKeypair := keypair.MustParseFull("SBLTFVCSTVKJEH2WXY5HFXHOPEF3EDBYUK7X6I3MFDXC3PJAJX2L3BPE") // Replace with actual admin secret key

	// Obtain sourceAccount from client
	sourceAccount, err := client.AccountDetail(auroraclient.AccountRequest{AccountID: sourceKeypair.Address()})
	if err != nil {
		return err
	}

	// Create manageData operation
	manageData := txnbuild.ManageData{
		Name:  name,
		Value: []byte(value),
	}

	// Set transaction parameters
	txParams := txnbuild.TransactionParams{
		SourceAccount:        &sourceAccount,
		IncrementSequenceNum: true,
		Operations:           []txnbuild.Operation{&manageData},
		BaseFee:              txnbuild.MinBaseFee,
		Timebounds:           txnbuild.NewTimeout(300),
	}

	// Create transaction
	tx, err := txnbuild.NewTransaction(txParams)
	if err != nil {
		return err
	}

	// Sign transaction
	tx, err = tx.Sign(networkPassphrase, sourceKeypair)
	if err != nil {
		return err
	}

	// Submit transaction to Diamante network
	_, err = client.SubmitTransaction(tx)
	if err != nil {
		return err
	}

	return nil
}

// SetOptions sets various configuration options for an account on the Stellar network
func SetOptions(opts txnbuild.SetOptions) error {
	// Initialize network passphrase and client
	networkPassphrase := network.TestNetworkPassphrase
	client := auroraclient.DefaultTestNetClient

	// Initialize source keypair from secret seed
	sourceKeypair := keypair.MustParseFull("SBLTFVCSTVKJEH2WXY5HFXHOPEF3EDBYUK7X6I3MFDXC3PJAJX2L3BPE") // Replace with actual admin secret key

	// Get account object for source keypair's address
	sourceAccount, err := client.AccountDetail(auroraclient.AccountRequest{AccountID: sourceKeypair.Address()})
	if err != nil {
		return err
	}

	// Create set options object with configuration options set
	setOptions := txnbuild.SetOptions{
		InflationDestination: opts.InflationDestination,
		Signer:               opts.Signer,
		HomeDomain:           opts.HomeDomain,
		MasterWeight:         opts.MasterWeight,
		LowThreshold:         opts.LowThreshold,
		HighThreshold:        opts.HighThreshold,
		SetFlags:             opts.SetFlags,
		ClearFlags:           opts.ClearFlags,
	}

	// Create transaction parameters object
	txParams := txnbuild.TransactionParams{
		SourceAccount:        &sourceAccount,
		IncrementSequenceNum: true,
		Operations:           []txnbuild.Operation{&setOptions},
		BaseFee:              txnbuild.MinBaseFee,
		Timebounds:           txnbuild.NewTimeout(300),
	}

	// Create new transaction object from transaction parameters
	tx, err := txnbuild.NewTransaction(txParams)
	if err != nil {
		return err
	}

	// Sign transaction with source keypair's secret key and network passphrase
	tx, err = tx.Sign(networkPassphrase, sourceKeypair)
	if err != nil {
		return err
	}

	// Submit tx to the Diamante network
	_, err = client.SubmitTransaction(tx)
	if err != nil {
		return err
	}

	return nil
}

// CreateAsset creates and submits an asset on the Diamante network
func CreateAsset(distributorPublicKey string, distributorPrivateKey string, tokenSupply int, tokenName string) (string, string, error) {
	client := auroraclient.DefaultTestNetClient

	// Create Diamante keypairs for issuance and distributor
	issuanceKP, err := keypair.Random()
	if err != nil {
		return "", "", fmt.Errorf("error generating issuance keypair: %v", err)
	}

	distributorKP := keypair.MustParseFull(distributorPrivateKey)
	distAcctReq := auroraclient.AccountRequest{AccountID: distributorKP.Address()}
	distributorAccount, err := client.AccountDetail(distAcctReq)
	if err != nil {
		return "", "", fmt.Errorf("error in distributor account: %v", err)
	}

	// Activation transaction for issuance account
	activationTx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &distributorAccount,
			IncrementSequenceNum: true,
			BaseFee:              txnbuild.MinBaseFee,
			Timebounds:           txnbuild.NewInfiniteTimeout(),
			Operations: []txnbuild.Operation{
				&txnbuild.CreateAccount{
					Destination: issuanceKP.Address(),
					Amount:      "4",
				},
			},
		},
	)
	if err != nil {
		return "", "", fmt.Errorf("error creating activation transaction: %v", err)
	}

	// Sign and submit the activation transaction
	activationTx, err = activationTx.Sign(network.TestNetworkPassphrase, distributorKP)
	if err != nil {
		return "", "", fmt.Errorf("error signing activation transaction: %v", err)
	}

	_, err = client.SubmitTransaction(activationTx)
	if err != nil {
		return "", "", fmt.Errorf("error submitting activation transaction: %v", err)
	}

	// Change trust operation for the distributor
	issuerAcctReq := auroraclient.AccountRequest{AccountID: issuanceKP.Address()}
	issuerAccount, err := client.AccountDetail(issuerAcctReq)
	if err != nil {
		return "", "", fmt.Errorf("error getting issuer account details: %v", err)
	}

	tokens, _ := txnbuild.CreditAsset{Code: tokenName, Issuer: issuanceKP.Address()}.ToChangeTrustAsset()

	changeTrustTx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &distributorAccount,
			IncrementSequenceNum: true,
			BaseFee:              txnbuild.MinBaseFee,
			Timebounds:           txnbuild.NewInfiniteTimeout(),
			Operations: []txnbuild.Operation{
				&txnbuild.ChangeTrust{
					Line: tokens,
				},
			},
		},
	)
	if err != nil {
		return "", "", fmt.Errorf("error creating change trust transaction: %v", err)
	}

	// Sign and submit the change trust transaction
	changeTrustTx, err = changeTrustTx.Sign(network.TestNetworkPassphrase, distributorKP)
	if err != nil {
		return "", "", fmt.Errorf("error signing change trust transaction: %v", err)
	}

	_, err = client.SubmitTransaction(changeTrustTx)
	if err != nil {
		return "", "", fmt.Errorf("error submitting change trust transaction: %v", err)
	}

	// Transfer tokens from issuer to distributor
	_tokens := txnbuild.CreditAsset{Code: tokenName, Issuer: issuanceKP.Address()}
	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &issuerAccount,
			IncrementSequenceNum: true,
			BaseFee:              txnbuild.MinBaseFee,
			Timebounds:           txnbuild.NewInfiniteTimeout(),
			Operations: []txnbuild.Operation{
				&txnbuild.Payment{
					Destination: distributorPublicKey,
					Asset:       _tokens,
					Amount:      fmt.Sprintf("%d", tokenSupply),
				},
			},
		},
	)
	if err != nil {
		return "", "", fmt.Errorf("error creating transaction: %v", err)
	}

	// Sign and submit the transaction
	tx, err = tx.Sign(network.TestNetworkPassphrase, issuanceKP)
	if err != nil {
		return "", "", fmt.Errorf("error signing transaction: %v", err)
	}

	resp1, err := client.SubmitTransaction(tx)
	if err != nil {
		return "", "", fmt.Errorf("error submitting transaction: %v", err)
	}

	log.Printf("Token transferred successfully. Hash: %s\n", resp1.Hash)

	// Lock the issuance account
	lockIssuanceTx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &issuerAccount,
			IncrementSequenceNum: true,
			BaseFee:              txnbuild.MinBaseFee,
			Timebounds:           txnbuild.NewInfiniteTimeout(),
			Operations: []txnbuild.Operation{
				&txnbuild.SetOptions{
					MasterWeight: txnbuild.NewThreshold(0),
				},
			},
		},
	)
	if err != nil {
		return "", "", fmt.Errorf("error creating lock issuance transaction: %v", err)
	}

	// Sign and submit the lock issuance transaction
	lockIssuanceTx, err = lockIssuanceTx.Sign(network.TestNetworkPassphrase, issuanceKP)
	if err != nil {
		return "", "", fmt.Errorf("error signing lock issuance transaction: %v", err)
	}

	resp2, err := client.SubmitTransaction(lockIssuanceTx)
	if err != nil {
		return "", "", fmt.Errorf("error submitting lock issuance transaction: %v", err)
	}

	fmt.Println("Issuance account locked successfully.", resp2.Hash)

	return issuanceKP.Address(), resp1.Hash, nil
}

// FundAndActivateAccount activates the account by funding it using the Diamcircle Friendbot
func FundAndActivateAccount(address string) error {
	// Use Diamcircle Friendbot to fund and activate the account
	friendbotURL := "https://friendbot.diamcircle.io/?addr=" + address
	resp, err := http.Get(friendbotURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fund and activate account, status code: %d", resp.StatusCode)
	}

	return nil
}
func main() {
	// Create a new random keypair for the destination account
	destinationKP, err := keypair.Random()
	if err != nil {
		log.Fatal("Error generating random keypair:", err)
	}
	log.Printf("Generated new random keypair for destination account: Address - %s, Seed - %s\n", destinationKP.Address(), destinationKP.Seed())

	// Use the Diamcircle Friendbot to fund and activate the account
	err = FundAndActivateAccount(destinationKP.Address())
	if err != nil {
		log.Fatal("Error funding and activating account:", err)
	}
	log.Printf("Account %s successfully funded and activated using Friendbot\n", destinationKP.Address())

	// Example usage of other functions
	err = Payment(destinationKP.Address(), "10") // Example payment of 10 DIAM to the newly created account
	if err != nil {
		log.Fatal("Error making payment:", err)
	}
	log.Println("Payment successful")

	err = ManageData("example_data", "example_value") // Example of managing data on the account
	if err != nil {
		log.Fatal("Error managing data:", err)
	}
	log.Println("Data managed successfully")

	// Example usage of setting options
	options := txnbuild.SetOptions{
		Signer: &txnbuild.Signer{
			Address: "GAZT2S7GV6Z7QIS3HT32BRN42DFVSSULPA3SQWNE4NY2MAOKVSC3K7N5",
			Weight:  1,
		},
	}
	err = SetOptions(options)
	if err != nil {
		log.Fatal("Error setting options:", err)
	}
	log.Println("Options set successfully")

	// Example usage of creating an asset
	assetCode := "TEST"
	issuerSecret := "SBQWCMO7YBNGL73HHW4L3HUBGUKHFU3TXZAAQB5LIFU7MDG6I7BRBHCH" // Replace with actual issuer secret key
	tokenSupply := 10000
	sourceKeypair := keypair.MustParseFull("SBLTFVCSTVKJEH2WXY5HFXHOPEF3EDBYUK7X6I3MFDXC3PJAJX2L3BPE") // Replace with actual admin secret key

	addr, hash, _ := CreateAsset(sourceKeypair.Address(), issuerSecret, tokenSupply, assetCode)
	if err != nil {
		log.Fatal("Error creating asset:", err)
	}
	log.Println("Asset created successfully", addr, hash)
}
