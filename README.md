# MineFrame
A Web-based admin GUI for minecraft
## dev instructions
### Setup your Minecraft Spigot server
* download spigot.jar from https://getbukkit.org
* make a user for the minecraft server. We are using "miner"
** adduser miner
*** cd /home/miner/
** make a directory for the server to run from
*** mkdir server
*** cd server
** move the spigot jar file into the server dir
** Make sure you can run the minecraft server
*** java -jar spigot-1.12.2.jar
** You may have to accept the eula by running nano eula.txt and changing eula=false to eula=true
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
  * This launches the react server on 3000
  * The react server will proxy requests to Revel so all dev can be done via http://localhost:3000/
## build & commit instructions:
* before committing, copy the react built .js into the revel web root:
* copy the contents of /react-project/build/ into /public/
* commit
