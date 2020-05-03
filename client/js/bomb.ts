import * as PIXI from "pixi.js";
import {EntityInfo} from "./entity";
import {Const} from "./gamescreen";

class Bomb implements EntityInfo {
    private static readonly TIMER_LENGTH = 120;

    private readonly _sprite: PIXI.Sprite;
    private _state: number;

    constructor() {
        this._sprite = PIXI.Sprite.from('assets/res/bomb0.png');
        this._sprite.width = Const.BLOCK_SIZE;
        this._sprite.height = Const.BLOCK_SIZE;
        this._state = Bomb.TIMER_LENGTH;
    }

    get sprite(): PIXI.Sprite {
        return this._sprite;
    }

    x(value: number): void {
        this._sprite.x = value * Const.BLOCK_SIZE;
    }

    y(value: number): void {
        this._sprite.y = value * Const.BLOCK_SIZE;
    }
}

export default Bomb;
