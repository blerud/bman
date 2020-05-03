import * as PIXI from "pixi.js";
import {EntityInfo} from "./entity";
import {Const} from "./gamescreen";

export enum Direction {
    UP,
    DOWN,
    LEFT,
    RIGHT
}

export class PlayerConsts {
    HEIGHT = 1;
    WIDTH = 1;
}

class Player implements EntityInfo {
    private _xSpeed: number;
    private _ySpeed: number;
    private _direction: Direction;
    private _hasKick: boolean;
    private _numBombs: number;
    private _numBombsAvailable: number;
    private _numFire: number;
    private _sprite: PIXI.Sprite;
    private readonly _heightOffset: number;

    constructor() {
        this._xSpeed = 0;
        this._ySpeed = 0;

        this._direction = Direction.DOWN;

        this._hasKick = false;
        this._numBombs = 1;
        this._numBombsAvailable = this._numBombs;
        this._numFire = 1;

        this._sprite = PIXI.Sprite.from('assets/res/player0.png');
        this._sprite.width = Const.BLOCK_SIZE;
        this._sprite.height = 60 / 36 * Const.BLOCK_SIZE;
        this._heightOffset = this._sprite.height - Const.BLOCK_SIZE;
    }

    get xSpeed(): number {
        return this._xSpeed;
    }

    set xSpeed(value: number) {
        this._xSpeed = value;
    }

    get ySpeed(): number {
        return this._ySpeed;
    }

    set ySpeed(value: number) {
        this._ySpeed = value;
    }

    get direction(): Direction {
        return this._direction;
    }

    set direction(value: Direction) {
        this._direction = value;
    }

    get hasKick(): boolean {
        return this._hasKick;
    }

    set hasKick(value: boolean) {
        this._hasKick = value;
    }

    get numBombs(): number {
        return this._numBombs;
    }

    set numBombs(value: number) {
        this._numBombs = value;
    }

    get numBombsAvailable(): number {
        return this._numBombsAvailable;
    }

    set numBombsAvailable(value: number) {
        this._numBombsAvailable = value;
    }

    get numFire(): number {
        return this._numFire;
    }

    set numFire(value: number) {
        this._numFire = value;
    }

    get sprite(): PIXI.Sprite {
        return this._sprite;
    }

    set sprite(value: PIXI.Sprite) {
        this._sprite = value;
    }

    x(value: number): void {
        this.sprite.x = value * Const.BLOCK_SIZE;
    }

    y(value: number): void {
        this.sprite.y = value * Const.BLOCK_SIZE - this._heightOffset;
    }
}

export default Player;