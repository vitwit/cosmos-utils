# proposal-vote-script

This script will help you to vote for the proposal if it's not voted already.
You just have to clone the repo, configure the `config.toml` and run the the script.

Clone the repo
```sh
git clone https://github.com/vitwit/cosmos-utils.git

cd cosmos-utils/proposal-vote-script
```

#### Configure the config.toml

```sh
cp example.config.toml config.toml
```
and replace these values with your validator details.

- *lcd_endpoint*

     LCD endpoint of your validator, which will be used to get proposals info and votes of it.

- *deamon*

     Deamon name of the network (ex: regen, akash). This will be used to execute the vote tx command.

- *key_name*

     Key name of your validator. Name of the account from which you want to vote.

- *account_address*

     Account address of your validator. This address will be used to get your votes for particular proposal.

- *chain_id*

     Chain ID of your node.

- *fees*

     Fees to execute the tx. 

#### Run the script

-  Build and run the using binary

```sh
go build -o proposal-script && ./proposal-script
```
- or run
```sh
go run main.go
```