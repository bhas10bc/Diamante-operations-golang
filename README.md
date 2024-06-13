Diamante Network Client Example

This repository provides an example of how to interact with the Diamante blockchain network using the Go programming language. The example includes various functions to perform operations such as creating accounts, making payments, managing data, setting account options, and creating assets.

Prerequisites
Before running this code, ensure you have the following installed:
Go (1.16+)
An IDE or text editor (such as Visual Studio Code)

Getting Started
Clone the Repository:
Install Dependencies:
Ensure you have the github.com/diamcircle/go package installed. 
 
Code Overview
The code provides several functions to interact with the Diamante network:
Payment:
Creates and submits a payment operation on the Diamante network.
ManageData:
Sets, updates, or deletes a data entry for an account.
SetOptions:
Sets various configuration options for an account on the Diamante network.
CreateAsset:
Creates and submits an asset on the Diamante network.
FundAndActivateAccount:
Activates an account by funding it using the Diamcircle Friendbot.
Running the Example
To run the provided example, use the main function. It demonstrates creating a new random keypair, funding and activating the account, making a payment, managing data, setting options, and creating an asset.
