# E-Wallet Top Up System

E-Wallet Top Up System is a simple simulation of an e-wallet top up system, which tries to mimic the flow of popular e-wallets in Indonesia. This is a small project that utilizes Go (Golang) to build a REST API backend system.

## Table of Content
1. [Features](#features)
2. [Technologies](#technologies)
3. [Setup](#setup)
4. [Simulating Direct E-Wallet Top Up](#simulating-direct-e-wallet-top-up)
5. [Simulating Top Up from Bank](#simulating-top-up-from-bank)
6. [Unit Testing](#unit-testing)
7. [Troubleshooting](#troubleshooting)

    7.1 [CGO / gcc / C compiler error](#cgo--gcc--c-compiler-error)
    
    7.1.1 [Windows](#windows)

    7.1.2 [Linux](#linux)

## Features

- Direct top up simulation
- Top up from bank simulation
- Transaction history
- Rate limits

## Technologies

- Go + Fiber (web framework)
- GORM (ORM)
- SQLite (local, file-based DB)
- Swagger (API documentation)

## Setup 

1. Clone the repository into a directory

    HTTPS

        git clone https://github.com/karsteneugene/top-up-system.git

    SSH

        git clone git@github.com:karsteneugene/top-up-system.git

    or download the zip file and extract.

2. Open up the terminal inside the folder or cd into the folder

        cd top-up-system

3. Once inside, run the `main.go` file

        go run main.go

4. The API should now be running locally! To make things simpler, the Swagger API Documentation is implemented to easily test endpoints. Simply visit http://localhost:3000/swagger or http://127.0.0.1:3000/swagger.

    To close the server, just do Ctrl + C in the console

That's it for the setup! No database setup is needed as this project uses SQLite for its database that is already prepared for users, which is the `ewallet.db` file in the root folder.

## Simulating Direct E-Wallet Top Up

This simulates a direct top up from popular e-wallets in Indonesia, where a card or an account is already bound to the e-wallet.

1. Head to http://localhost:3000/swagger or http://127.0.0.1:3000/swagger.
2. Click the POST `/transactions/topup/direct/{id}` endpoint under transactions tag.
3. Click the **Try it out** button on the right.
4. Input the Wallet ID to top up and input an amount to add.

    Note: There are rate limits implemented, users should not be able to:
    
    - input anything under 1000 (Rp 1,000) and anything above 2000000 (Rp 2,000,000) (per transaction limit).
    - top up to a sum of more than 5000000 (Rp 5,000,000) in a day (daily limit - resets daily).
    - top up to a sum of more than 20000000 (Rp 20,000,000) in a month (monthly limit - resets monthly).

    Another note: If you are unsure what Wallet IDs are available, do check the GET `/wallets` endpoint under wallets tag.

5. Once done inputting, click the blue `Execute` button.
6. If it is successful, it should show the transaction being made.

    Check the GET `/transactions/wallet/{id}` endpoint under the transactions tag to see the transaction history of a wallet.

## Simulating Top Up from Bank

First, be sure to get the virtual account number of the wallet you want to top up by going to http://localhost:3000/swagger or http://127.0.0.1:3000/swagger. Then heading to the GET `/wallets/va/{id}` endpoint under the wallets tag.

There are 2 ways to simulate top up from bank. The easiest way is running the script `topup_bca.go` by doing:

    go run ./scripts/topup_bca.go

This will automatically input the recipient's bank code and account number. All you need to do is enter the virtual account number, the amount to top up, and the description (optional).

But if you want to input the bank code and account manually, while also looking at the JSON request and response body:

1. Head to http://localhost:3000/swagger or http://127.0.0.1:3000/swagger and go to the POST `/transactions/topup/bank/{va}` endpoint under the transactions tag.
2. Input the virtual account number and you can proceed to edit the JSON request body.

    For account number and bank codes, please refer to the `validate_bank.go` file in `utils/validate_bank.go`. This is to simulate bank and account validation.

3. Hit execute and done! You can check the GET `/transactions/wallet/{id}` endpoint under the transactions tag to see the transaction you just made.

## Unit Testing

A unit test was also implemented for this small project to test a few cases. However, this unit test is only done on the direct top up endpoint and not the bank one as they pretty much have the same logic aside from validating the banks and accounts, which are hardcoded anyways.

To run the unit test, simply type this in the console:

    go test ./tests/api/handlers

or

    go test ./tests/api/handlers -v 

for verbose mode.

This tests for:

- Successful direct top up (Transaction successful)
- Wallet not found (Invalid wallet ID)
- Less than minimum top up amount per transaction (Amount lower than Rp 1,000)
- Exceeds the maximum top up amount per transaction (Amount higher than Rp 2,000,000)
- Exceeds daily limit (Total transactions will exceed or are exceeding Rp 5,000,000)
- Exceeds monthly limit (Total transactions will exceed or are exceeding Rp 20,000,000)

## Troubleshooting

### CGO / gcc / C compiler error  

When trying to run the application and you get this error:

    failed to initialize database, got error Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub

Proceed to type this inside the console:

    go env -w CGO_ENABLED=1

If you try to run the application again and get this error:

    # runtime/cgo
    cgo: C compiler "gcc" not found: exec: "gcc": executable file not found in %PATH%

That means gcc is not installed on your machine. 

#### Windows

For my Windows machine, I installed [TDM-GCC](https://jmeubank.github.io/tdm-gcc/) as I feel like it's the simplest one. TDM-GCC should add itself to PATH when installing with defualt settings.

Now restart your console and you should be able to run the application now.


#### Linux

For my Ubuntu that was run in WSL, I ran this command:

    sudo apt-get install build-essential

If there was an error retrieving some files, run this command:

    sudo apt-get update

and then rerun the first command to install build-essential.

The application should be able to run now.