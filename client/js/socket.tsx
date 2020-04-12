import {InitInfo} from "./game";

class Socket {
    private username: string;
    private userid: number;
    private gameid: number;
    private sock: WebSocket;

    constructor(initInfo: InitInfo) {
        this.username = initInfo.username;
        this.userid = initInfo.userid;
        this.gameid = initInfo.gameid;
        this.sock = null;

        this.connect();
    }

    connect(): boolean {
        let params = {
            userId: this.userid.toString(),
            name: this.username,
            gameId: this.gameid.toString(),
        };
        let paramsStr = new URLSearchParams(params).toString();
        let url = new URL('/ws?' + paramsStr, window.location.href);
        url.protocol = url.protocol.replace('http', 'ws');
        this.sock = new WebSocket(url.toString());
        return true; // todo return false if connection failure
    }

    connected(): boolean {
        return this.sock == null;
    }
}

export default Socket;
