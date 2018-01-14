import React, { Component } from 'react';
import logo from './block.png';
import './App.css';
import axios from 'axios';
import ServerConfigEditor from './ServerConfigEditor';

class App extends Component {
  constructor(props) {
    super(props);
    this.state = {
      ServerStatus: "Checking...",
      btnStartDisabled: true,
      btnStopDisabled: true,
      btnRestartDisabled: true
    };
  }

  minecraftServerAction(action){
    console.log('minecraftServerAction! '+action);

    switch(action){
      default:
        break;
      case "start":
        this.startService();
        break;
      case "restart":
        this.setState({ServerStatus: "Running"});
        this.setButtonStates("Running");
      break;
      case "stop":
        this.setState({ServerStatus: "Stopped"});
        this.setButtonStates("Stopped");
        break;
    }
  }

  startService(){
  var me = this;
  axios.post('/MfApi/ServiceStart') //, {Post:Name}
  .then(function (response) {
    //console.log(response);
    
    me.setState({ServerStatus: response.data.ServerStatus});
    me.setButtonStates(response.data.ServerStatus);
  })
  .catch(function (error) {
    console.log(error);
    me.setState({ServerStatus: "Error!"});
  });
}

  getServerStatus(){
    var me = this;
    axios.post('/MfApi/ServicetStatus') //, {Post:Name}
    .then(function (response) {
      //console.log(response);

      //old school
      //document.getElementById('ServerStatus').innerHTML = response.data.ServerStatus;
      
      //new school
      me.setState({ServerStatus: response.data.ServerStatus});
      me.setButtonStates(response.data.ServerStatus);
    })
    .catch(function (error) {
      console.log(error);
      me.setState({ServerStatus: "Error!"});
    });
  }

  setButtonStates(status){
    switch(status){
      default:
        break;
      case "Running":
        this.setState({
          btnStartDisabled: true,
          btnStopDisabled: false,
          btnRestartDisabled: false
        });
        break;
      case "Stopped":
        this.setState({
          btnStartDisabled: false,
          btnStopDisabled: true,
          btnRestartDisabled: true
        });
        break;
    }
  }

  initConsoleWebsocket(){
    // example from https://revel.github.io/examples/chat.html
    // var socket = new WebSocket('ws://127.0.0.1:9000/websocket/room/socket?user={{.user}}');
  
    // // Message received on the socket
    // socket.onmessage = function(event) {
    //     display(JSON.parse(event.data));
    // }
  }

  componentDidMount() {
    this.getServerStatus();
    this.initConsoleWebsocket();
}

  render() {
    return (
      <div className="App">
        <header className="App-header">
          <img src={logo} className="App-logo" alt="logo" />
          <h1 className="App-title">MineFrame Server</h1>
        </header>
        <p className="ServerStatus">Server Status: {this.state.ServerStatus}</p>
        <p>
        <button disabled={this.state.btnStartDisabled} className="btn-service-control" onClick={(e) => {this.minecraftServerAction('start',e)}}>
          Start
        </button>
        <button disabled={this.state.btnStopDisabled} className="btn-service-control" onClick={(e) => {this.minecraftServerAction('stop',e)}}>
          Stop
        </button>
        <button disabled={this.state.btnRestartDisabled} className="btn-service-control" onClick={(e) => {this.minecraftServerAction('restart',e)}}>
          Restart
        </button>
        </p>
        <ServerConfigEditor />
      </div>
    );
  }
}

export default App;
