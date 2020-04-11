import Socket from "./socket";

class Game {
    private app: PIXI.Application;
    private sock: Socket;

    constructor(app: PIXI.Application, sock: Socket) {
        this.app = app;
        this.sock = sock;
    }
}

export default Game;
