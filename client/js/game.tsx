import Socket from "./socket";
import * as PIXI from "pixi.js";

export interface InitInfo {
    username: string;
    userid: number;
    gameid: number;
}

class Game {
    private app: PIXI.Application;
    private sock: Socket;

    constructor(app: PIXI.Application, sock: Socket) {
        this.app = app;
        this.sock = sock;
    }

    renderBackground() {
        let background = new PIXI.Graphics()
            .beginFill(0x000000)
            .drawRect(0, 0, this.app.view.width, this.app.view.height)
            .endFill();
        this.app.stage.addChild(background);
    }

    renderPlayer() {
        let sprite = PIXI.Sprite.from('assets/res/bomb1.png');
        this.app.stage.addChild(sprite);
    }
}

export default Game;
