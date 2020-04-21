import Socket, {Message} from "./socket";
import * as PIXI from "pixi.js";
import {Keys, MessageTypes} from "./consts";
import Player from "./player";
import {PlayerMessage} from "./messages";

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
        switch (message.type) {
            case MessageTypes.SERVER_HEARTBEAT:
                this.sock.sendClientHeartbeat(message.timestamp);
                return true;
            case MessageTypes.CREATED:
                return true;
            case MessageTypes.UPDATED:
                return this.processUpdated(message);
            case MessageTypes.DELETED:
                return true;
        }

        return false;
    }

    processUpdated(message: Message): boolean {
        let contentView = new DataView(message.content);
        let numUpdated = contentView.getUint8(0);
        let offset = 1;
        for (let i = 0; i < numUpdated; i++) {
            let entityType = contentView.getUint8(offset);
            if (entityType == PlayerMessage.TYPE) {
                let playerMessage = PlayerMessage.fromBytes(new DataView(message.content, offset));
                offset += PlayerMessage.LENGTH;

                this.player.x = playerMessage.posX;
                this.player.y = playerMessage.posY;
                this.player.sprite.x = playerMessage.posX;
                this.player.sprite.y = playerMessage.posY;
            }
        }
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
