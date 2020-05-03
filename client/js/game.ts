import Socket, {Message} from "./socket";
import * as PIXI from "pixi.js";
import {Keys, MessageTypes} from "./consts";
import Player from "./player";
import {HardWallMessage, PlayerMessage} from "./messages";
import Entity from "./entity";
import HardWall from "./hardwall";

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
    private entities: Map<number, Entity>;

    constructor(app: PIXI.Application, sock: Socket, initInfo: InitInfo) {
        this.app = app;
        this.sock = sock;
        this.entities = new Map<number, Entity>();

        let loader = PIXI.Loader.shared;
        loader.load(this.setup.bind(this));
    }

    setup() {
        console.log("setup");

        this.sock.registerMessageHandler(this.processMessage.bind(this));
        this.setupKeyboard();
        this.sock.init();
    }

    processMessage(message: Message): boolean {
        switch (message.type) {
            case MessageTypes.SERVER_HEARTBEAT:
                this.sock.sendClientHeartbeat(message.timestamp);
                return true;
            case MessageTypes.CREATED:
                return this.processCreated(message);
            case MessageTypes.UPDATED:
                return this.processUpdated(message);
            case MessageTypes.DELETED:
                return this.processDeleted(message);
        }

        return false;
    }

    processCreated(message: Message): boolean {
        let contentView = new DataView(message.content);
        let numUpdated = contentView.getUint8(0);
        let offset = 1;
        for (let i = 0; i < numUpdated; i++) {
            let entityType = contentView.getUint8(offset);
            if (entityType == PlayerMessage.TYPE) {
                let playerMessage = PlayerMessage.fromBytes(new DataView(message.content, offset));
                offset += PlayerMessage.LENGTH;

                let playerInfo = new Player();
                let player = new Entity(playerMessage.id, playerMessage.posX, playerMessage.posY, 1, 1, playerInfo);

                this.entities.set(playerMessage.id, player);
                this.app.stage.addChild(playerInfo.sprite);
            } else if (entityType == HardWallMessage.TYPE) {
                let hardWallMessage = HardWallMessage.fromBytes(new DataView(message.content, offset));
                offset += HardWallMessage.LENGTH;

                let hardWallInfo = new HardWall();
                let hardWall = new Entity(hardWallMessage.id, hardWallMessage.posX, hardWallMessage.posY, 1, 1, hardWallInfo);

                this.entities.set(hardWallMessage.id, hardWall);
                this.app.stage.addChild(hardWallInfo.sprite);
            }
        }
        return true;
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

                let player = this.entities.get(playerMessage.id);
                if (player == undefined) {
                    continue;
                }

                player.x = playerMessage.posX;
                player.y = playerMessage.posY;
                let playerInfo = player.entityInfo as Player;
                playerInfo.direction = playerMessage.direction;
            }
        }
        return true;
    }

    processDeleted(message: Message): boolean {
        let contentView = new DataView(message.content);
        let numUpdated = contentView.getUint8(0);
        let offset = 1;
        for (let i = 0; i < numUpdated; i++) {
            let entityType = contentView.getUint8(offset);
            if (entityType == PlayerMessage.TYPE) {
                let playerMessage = PlayerMessage.fromBytes(new DataView(message.content, offset));
                offset += PlayerMessage.LENGTH;

                let player = this.entities.get(playerMessage.id);
                if (player == undefined) {
                    continue;
                }

                this.app.stage.removeChild(player.entityInfo.sprite);
                this.entities.delete(player.id);
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
    }
}

export default Game;
