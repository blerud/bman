import Socket, {Message} from "./socket";
import * as PIXI from "pixi.js";
import {Keys} from "./consts";
import Player from "./player";

export interface InitInfo {
    username: string;
    userid: number;
    gameid: number;
}

class Game {
    private app: PIXI.Application;
    private sock: Socket;
    private keys: Map<number, boolean>;
    private player: Player;

    constructor(app: PIXI.Application, sock: Socket, initInfo: InitInfo) {
        this.app = app;
        this.sock = sock;
        this.player = new Player(0, 0, initInfo.userid);

        let loader = PIXI.Loader.shared;
        loader.add('bomb1', 'assets/res/bomb1.png');
        loader.load(this.setup.bind(this));
    }

    setup() {
        console.log("setup");

        this.sock.registerMessageHandler(this.processMessage.bind(this));
        this.setupKeyboard();
    }

    processMessage(message: Message): boolean {
        // console.log("event messageType: " + message.type + " messageLength: " + message.length + " timestamp: " + message.timestamp + " content: " + new Uint8Array(message.content));
        // for now we're assuming only update message, and only 1 obj in it on id 6969
        let contentView = new DataView(message.content);
        let numUpdated = contentView.getUint8(0);
        let entityType = contentView.getUint8(1)
        let id = contentView.getUint32(2);
        let posX = contentView.getFloat32(6);
        let posY = contentView.getFloat32(10);
        let action = contentView.getUint8(14);
        let direction = contentView.getUint8(15);
        this.player.x = posX;
        this.player.y = posY;
        this.player.sprite.x = posX;
        this.player.sprite.y = posY;
        // console.log(new Uint8Array(message.content), numUpdated, id, posX, posY, action, direction);
        return true;
    }

    setupKeyboard() {
        this.keys = new Map<number, boolean>();
        this.keys.set(Keys.UP, false);
        this.keys.set(Keys.DOWN, false);
        this.keys.set(Keys.LEFT, false);
        this.keys.set(Keys.RIGHT, false);
        this.keys.set(Keys.BOMB, false);

        window.addEventListener("keydown", this.keyEvent.bind(this), false);
        window.addEventListener("keyup", this.keyEvent.bind(this), false);
    }

    keyEvent(event: KeyboardEvent) {
        if (event.defaultPrevented) {
            return;
        }

        let key = -1;

        switch (event.key) {
            case "w":
            case "W":
                key = Keys.UP;
                break;
            case "a":
            case "A":
                key = Keys.LEFT;
                break;
            case "s":
            case "S":
                key = Keys.DOWN;
                break;
            case "d":
            case "D":
                key = Keys.RIGHT;
                break;
            case " ":
                key = Keys.BOMB;
                break;
            case "Escape":
                break;
            default:
                return;
        }

        if (event.type == "keydown") {
            this.keyPressed(key);
        } else if (event.type == "keyup") {
            this.keyReleased(key);
        }

        this.sock.sendPlayerAction(this.keys);

        event.preventDefault();
    }

    keyPressed(key: number) {
        console.log("keydown " + key);
        this.keys.set(key, true);
    }

    keyReleased(key: number) {
        console.log("keyup " + key);
        this.keys.set(key, false);
    }

    render() {

    }

    renderBackground() {
        let background = new PIXI.Graphics()
            .beginFill(0xffffff)
            .drawRect(0, 0, this.app.view.width, this.app.view.height)
            .endFill();
        this.app.stage.addChild(background);
    }

    renderPlayer() {
        this.player.sprite.x = this.player.x;
        this.player.sprite.y = this.player.y;
        this.app.stage.addChild(this.player.sprite);
    }
}

export default Game;
