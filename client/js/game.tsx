import Socket from "./socket";
import * as PIXI from "pixi.js";
import {Keys} from "./consts";

export interface InitInfo {
    username: string;
    userid: number;
    gameid: number;
}

class Game {
    private app: PIXI.Application;
    private sock: Socket;
    private keys: Map<number, boolean>;

    constructor(app: PIXI.Application, sock: Socket) {
        this.app = app;
        this.sock = sock;

        let loader = PIXI.Loader.shared;
        loader.add('bomb1', 'assets/res/bomb1.png');
        loader.load(this.setup.bind(this));
    }

    setup() {
        console.log("setup");

        this.setupKeyboard();
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

        if (event.type == "keyup") {
            this.keyPressed(key);
        } else if (event.type == "keydown") {
            this.keyReleased(key);
        }

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
        let sprite = PIXI.Sprite.from('assets/res/bomb1.png');
        this.app.stage.addChild(sprite);
    }
}

export default Game;
