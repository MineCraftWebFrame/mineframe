import React, { Component } from 'react';
import logo from './block.png';
import './App.css';
import $ from 'jquery';
import Ajax from './Ajax.js';
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
        this.sendWebsocketMessage("StartService");
        //this.startService();
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

  Ajax({
    url:'/MfApi/ServiceStart',
    success:function (data) {
      //console.log(data);
      
      this.setState({ServerStatus: data.ServerStatus});
      this.setButtonStates(data.ServerStatus);
    },
    failure:function () {
      this.setState({ServerStatus: "Error!"});
    },
    scope:this
  });

}

  getServerStatus(){

    Ajax({
      url:'/MfApi/ServicetStatus',
      success:function (data) {
        console.log('status success');
        console.log(data);

        this.setState({ServerStatus: data.ServerStatus});
        this.setButtonStates(data.ServerStatus);
      },
      failure:function () {
        this.setState({ServerStatus: "Error!"});
      },
      scope:this
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
  sendWebsocketMessage(message){
    if(!this.socket){
      console.log('Websocket not initialized');
      return false;
    }

    if(this.socket.readyState !== WebSocket.OPEN){
      console.log('Websocket not connected');
      return false;
    }

    this.socket.send(message);
  }
  initConsoleWebsocket(){
    var me=this;
    // example from https://revel.github.io/examples/chat.html
    me.socket = new WebSocket('ws://127.0.0.1:9000/ServiceSocket');
    me.socket.onerror = function(){
      me.display("error with websocket!");
      return;
    }
    me.socket.onclose = function(){
      me.display("websocket closed!");
      return;
    }
    // Message received on the socket
    me.socket.onmessage = function(event) {
      me.display(JSON.parse(event.data));
      
    }
  }
  socketInputKey(event){
    if(event.key == 'Enter'){
       this.socketInput();
    }
  }
  socketInput(){
    var cmd = $('#wsInput');
    var val = cmd.val();
    //console.log(val);
    this.sendWebsocketMessage(val);
    cmd.val('');
  }
  display(event) {
    console.log('Socket Event!')
    console.log(event);
    var sc = $('#serviceConsole')
    sc.append('<div class="message">'+ event + '</div>');
    //sc.animate({ scrollTop: sc[0].scrollHeight - sc.height() }, "easeOutQuint");
    sc.scrollTop( sc[0].scrollHeight - sc.height() );
    //$('#serviceConsole').scrollTo('max')
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
        <div id="serviceConsole"></div>
        <p>
          Command: <input id="wsInput" onKeyPress={(e) => {this.socketInputKey(e)}} /> <button onClick={() => {this.socketInput()}}>Send</button>
        </p>
        <ServerConfigEditor />
      </div>
    );
  }
}

export default App;
