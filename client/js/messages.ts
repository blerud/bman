export class PlayerMessage {
    static readonly TYPE = 0;
    static readonly LENGTH = 15;
    id: number;
    posX: number;
    posY: number;
    action: number;
    direction: number;

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