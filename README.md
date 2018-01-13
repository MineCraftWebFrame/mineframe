# mineframe
A Web-based admin GUI for minecadt
## dev instructions
### get the code
* go get github.com/MineCraftWebFrame/mineframe
* Launch the revel server
  * cd $GOPATH/src/github.com/MineCraftWebFrame/mineframe/
  * revel run
  * This launches the server on port 9000
### Building React
* cd $GOPATH/src/github.com/MineCraftWebFrame/mineframe/react-project
  * npm run build

### Deving with React
* cd $GOPATH/src/github.com/MineCraftWebFrame/mineframe/react-project
* npm start
  * This launches the react server on 8080
  * The react server will proxy requests to Revel so all dev can be done via http://localhost:8080/
## build & commit instructions:
* before committing, copy the react built .js into the revel web root:
* copy the contents of /react-project/build/ into /public/
* commit
