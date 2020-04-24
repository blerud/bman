import * as React from 'react';
import {Screen} from "./consts";
import GameScreen from "./gamescreen";
import Socket from "./socket";

interface State {
    name: string
    serverid: string
    currentScreen: Screen
}

class App extends React.Component<{}, State> {
    private socket: WebSocket;
    private initInfo: { gameid: number; userid: number; username: string };
    private sock: Socket;

    constructor(props: any) {
        super(props);

        this.state = {
            name: '',
            serverid: '',
            currentScreen: Screen.MENU
        };

        this.handleClick = this.handleClick.bind(this);
        this.handleSubmit = this.handleSubmit.bind(this);
        this.handleServeridChange = this.handleServeridChange.bind(this);
        this.handleNameChange = this.handleNameChange.bind(this);
    }

    async handleClick(event: React.MouseEvent) {
        let req = new Request('/creategame', {
            method: 'POST'
        });
        let res = await fetch(req);
        let json = await res.json();
        console.log(json['gameId']);
        let username = "lend";
        let userid = "6969";
        let gameid = json['gameId'].toString();
        this.initInfo = {username: username, userid: Number(userid), gameid: Number(gameid)};
        this.sock = new Socket(this.initInfo);
        this.setState({currentScreen: Screen.GAME});
    }

    handleServeridChange(event: React.FormEvent<HTMLInputElement>) {
        this.setState({serverid: event.currentTarget.value});
    }

    handleNameChange(event: React.FormEvent<HTMLInputElement>) {
        this.setState({name: event.currentTarget.value});
    }

    async handleSubmit(event: React.FormEvent) {
        this.initInfo = {username: this.state.name, userid: 1234, gameid: Number(this.state.serverid)};
        this.sock = new Socket(this.initInfo);
        this.setState({currentScreen: Screen.GAME});
    }

    render() {
        switch (this.state.currentScreen) {
            case Screen.MENU:
                return (<div>
                    <button onClick={this.handleClick}>create server</button>
                    <br />
                    <input type="text" onChange={this.handleServeridChange} placeholder="server id" />
                    <input type="text" onChange={this.handleNameChange} placeholder="name" />
                    <button onClick={this.handleSubmit}>join server</button>
                </div>);
            case Screen.LOBBY:
            case Screen.GAME:
                return <GameScreen sock={this.sock} initInfo={this.initInfo} />;
            default:
                return <h1>bad</h1>;
        }
    }
}

export default App;
