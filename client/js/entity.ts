export interface EntityInfo {
    x(value: number): void;
    y(value: number): void;
    sprite: PIXI.Sprite;
}

class Entity {
    private readonly _id: number;
    private _x: number;
    private _y: number;
    private _width: number;
    private _height: number;
    private _entityInfo: EntityInfo;

    constructor(id: number, x: number, y: number, width: number, height: number, entityInfo: EntityInfo) {
        this._id = id;
        this.entityInfo = entityInfo;
        this.x = x;
        this.y = y;
        this.width = width;
        this.height = height;
    }

    get id(): number {
        return this._id;
    }

    get x(): number {
        return this._x;
    }

    set x(value: number) {
        this._x = value;
        this.entityInfo.x(value);
    }

    get y(): number {
        return this._y;
    }

    set y(value: number) {
        this._y = value;
        this.entityInfo.y(value);
    }

    get width(): number {
        return this._width;
    }

    set width(value: number) {
        this._width = value;
    }

    get height(): number {
        return this._height;
    }

    set height(value: number) {
        this._height = value;
    }

    get entityInfo(): EntityInfo {
        return this._entityInfo;
    }

    set entityInfo(value: EntityInfo) {
        this._entityInfo = value;
    }
}

export default Entity;
