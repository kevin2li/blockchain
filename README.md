
# Bitcoin data analysis tool

# Usage
1. download BTC blockchain data:

`tool download <heights>`

`tool download -r <start_height> <end_height>`

`tool download -f <path_to_heights>`

2. cluster BTC address:

`tool cluster -f <path_to_dataset> <address>`

3. visualize BTC transactions based on an entity's addresses:

`tool viz -d <path_to_dataset>  -a <path_to_address>`
