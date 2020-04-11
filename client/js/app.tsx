import * as React from 'react';

interface State {
    text: string
}

class App extends React.Component<{}, State> {
    private socket: WebSocket;
    constructor(props: any) {
        super(props);

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
        let params = {
            userId: userid,
            name: username,
            gameId: json['gameId'].toString()
        };
        let paramsStr = new URLSearchParams(params).toString();
        console.log(paramsStr);
        let url = new URL('/ws?' + paramsStr, window.location.href);
        url.protocol = url.protocol.replace('http', 'ws');
        this.socket = new WebSocket(url.toString());
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
        return (
        <div>
            <button onClick={this.handleClick}>post</button>
            <input type="text" onChange={this.handleChange} />
            <button onClick={this.handleSubmit}>asdf</button>
        </div>
        );
    }
}

export default App;
