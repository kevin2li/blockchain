
# BTC data analysis tool

# Installation

`go build -o tool`

# Usage
1. download BTC tool data:

`./tool download <heights>`

`./tool download -r <start_height> <end_height>`

`./tool download -f <path_to_heights>`

2. cluster BTC address:

`./tool cluster -f <path_to_dataset> <address>`

3. visualize BTC transactions based on an entity's addresses:

`docker run --name kevin-neo4j -p7474:7474 -p7687:7687 --rm -v $HOME/neo4j/data:/data -v $HOME/neo4j/logs:/logs -v $HOME/neo4j/import:/var/lib/neo4j/import -v $HOME/neo4j/plugins:/plugins --env NEO4J_AUTH=neo4j/neo4j neo4j:latest`

`./tool viz -d <path_to_dataset>  -a <path_to_address>`
