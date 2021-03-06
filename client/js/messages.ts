import {Direction} from "./player";

export class PlayerMessage {
    static readonly TYPE = 0;
    static readonly LENGTH = 15;
    id: number;
    posX: number;
    posY: number;
    action: number;
    direction: Direction;

    static fromBytes(bytes: DataView): PlayerMessage {
        let message = new PlayerMessage();
        message.id = bytes.getUint32(1);
        message.posX = bytes.getFloat32(5);
        message.posY = bytes.getFloat32(9);
        message.action = bytes.getUint8(13);
        message.direction = bytes.getUint8(14);
        return message;
    }
}

export class BombMessage {
    static readonly TYPE = 1;
    static readonly LENGTH = 14;
    id: number;
    posX: number;
    posY: number;
    state: number;

    static fromBytes(bytes: DataView): BombMessage {
        let message = new BombMessage();
        message.id = bytes.getUint32(1);
        message.posX = bytes.getFloat32(5);
        message.posY = bytes.getFloat32(9);
        message.state = bytes.getUint8(13);
        return message;
    }
}

export class HardWallMessage {
    static readonly TYPE = 3;
    static readonly LENGTH = 13;
    id: number;
    posX: number;
    posY: number;

    static fromBytes(bytes: DataView): HardWallMessage {
        let message = new HardWallMessage();
        message.id = bytes.getUint32(1);
        message.posX = bytes.getFloat32(5);
        message.posY = bytes.getFloat32(9);
        return message;
    }
}
