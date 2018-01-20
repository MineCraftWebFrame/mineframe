import React, { Component } from 'react';
import Ajax from './Ajax.js';
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
        
        var config = this.state.config;
        if(config == ""){
            console.log("Blank config!");
            return false;
        }
        
        this.setState({btnSaveConfigDisabled:true});

        Ajax({
            url:'/MfApi/MinecraftConfigUpdate',
            params:{config:this.state.config},
            success:function (data) {
              
                this.setState({
                    btnSaveConfigDisabled:false
                });
            },
            failure:function () {
                this.setState({
                    btnSaveConfigDisabled:false,  
                    serverConfig: "Error Reading Config From Server!"
                });
            },
            scope:this
          });

    }
    ServerConfigRead(){

        Ajax({
            url:'/MfApi/MinecraftConfigRead',
            params:{config:this.state.config},
            success:function (data) {
                
                this.setState({
                    btnSaveConfigDisabled:false, 
                    configFile:data.configFile,
                    config:data.config
                });
            },
            failure:function (error) {
                this.setState({serverConfig: "Error Reading Config From Server!"});
            },
            scope:this
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
