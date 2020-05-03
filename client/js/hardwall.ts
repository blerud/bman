import * as PIXI from "pixi.js";
import {EntityInfo} from "./entity";
import {Const} from "./gamescreen";

class HardWall implements EntityInfo {
    private readonly _sprite: PIXI.Sprite;

    constructor() {
        this._sprite = PIXI.Sprite.from('assets/res/wall.png');
        this._sprite.width = Const.BLOCK_SIZE;
        this._sprite.height = Const.BLOCK_SIZE
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

export default HardWall;
