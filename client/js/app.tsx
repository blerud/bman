import * as React from 'react';
import {Screen} from "./consts";
import GameScreen from "./gamescreen";
import Socket from "./socket";

interface State {
    text: string
    currentScreen: Screen
}

class App extends React.Component<{}, State> {
    private socket: WebSocket;
    private initInfo: { gameid: number; userid: number; username: string };
    private sock: Socket;

    constructor(props: any) {
        super(props);

        this.state = {
            text: '',
            currentScreen: Screen.MENU
        };

        this.handleClick = this.handleClick.bind(this);
        this.handleSubmit = this.handleSubmit.bind(this);
        this.handleChange = this.handleChange.bind(this);
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

    handleChange(event: React.FormEvent<HTMLInputElement>) {
        this.setState({text: event.currentTarget.value});
    }

    async handleSubmit(event: React.FormEvent) {
        let arr = new ArrayBuffer(5);
        let view = new DataView(arr);
        view.setUint8(0, 0x01);
        view.setUint32(1, this.state.text.length, false);
        this.socket.send(arr);
        this.socket.send(this.state.text);
    }

    render() {
        switch (this.state.currentScreen) {
            case Screen.MENU:
                return (<div>
                    <button onClick={this.handleClick}>post</button>
                    <input type="text" onChange={this.handleChange} />
                    <button onClick={this.handleSubmit}>asdf</button>
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
