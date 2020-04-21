export enum Screen {
    MENU,
    LOBBY,
    GAME,
}

export class Keys {
    static UP = 87;
    static DOWN = 83;
    static LEFT = 65;
    static RIGHT = 68;
    static BOMB = 32;
    static ESC = 27;
}

export class MessageTypes {
    static SERVER_HEARTBEAT = 0;
    static GAME_INIT = 2;
    static GAME_START = 3;
    static GAME_END = 6;
    static CREATED = 10;
    static UPDATED = 11;
    static DELETED = 12;

    static CLIENT_HEARTBEAT = 1;
    static CLIENT_ACTION = 20;
}

export class ObjectTypes {
    static PLAYER = 0;
}
