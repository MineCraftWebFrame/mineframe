import React, { Component } from 'react';
import axios from 'axios';
import './ServerConfigEditor.css';

class ServerConfigEditor extends Component {
    constructor(props) {
      super(props);
      this.state = {
        btnSaveConfigDisabled: true,
        configFile: "server.properties",
        config: ""
      };
    }
    ServerConfigFieldChange(event){
        this.setState({config: event.target.value});
    }
    ServerConfigUpdate(){
        var me = this;
        //console.log(this.state);
        var config = this.state.config;
        if(config == ""){
            console.log("Blank config!");
            return false;
        }
        
        this.setState({btnSaveConfigDisabled:true});
        axios.post('/MfApi/ServerConfigUpdate',{config:this.state.config})
        .then(function (response) {

            me.setState({
                btnSaveConfigDisabled:false
            });
        })
        .catch(function (error) {
            me.setState({
                btnSaveConfigDisabled:false,  
                serverConfig: "Error Reading Config From Server!"
            });
        });
    }
    ServerConfigRead(){
        var me = this;
        axios.post('/MfApi/ServerConfigRead')
        .then(function (response) {

            me.setState({
                btnSaveConfigDisabled:false, 
                configFile:response.data.configFile,
                config:response.data.config
            });
        })
        .catch(function (error) {
          me.setState({serverConfig: "Error Reading Config From Server!"});
        });
    }
    componentDidMount() {
        this.ServerConfigRead();
    }
    render() {
        return (
          <div className="ServerConfigEditor">
            <hr />
            <h2>Server Config Editor</h2>
            <h3>{this.state.configFile}</h3>
                <textarea id="serverConfigTextarea" value={this.state.config}  onChange={(e)=>{this.ServerConfigFieldChange(e)}}></textarea>
                <br />
                <button disabled={this.state.btnSaveConfigDisabled} className="btn-service-control" onClick={(e)=>{this.ServerConfigUpdate(e)}}>
                Save Changes
                </button>
            <hr />
          </div>
        );
      }
}

export default ServerConfigEditor;
