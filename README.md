
# Bitcoin data analysis tool

## Installation

`go build`

# Usage
1. download BTC blockchain data:

`./blockchain download <heights>`

`./blockchain download -r <start_height> <end_height>`

`./blockchain download -f <path_to_heights>`

2. cluster BTC address:

`./blockchain cluster -f <path_to_dataset> <address>`

3. visualize BTC transactions based on an entity's addresses:

`./blockchain viz -d <path_to_dataset>  -a <path_to_address>`
