import * as React from 'react';
import * as PIXI from 'pixi.js';
import Game from "./game";
import Socket from "./socket";

interface Props {
    sock: Socket;
    name: string;
}

interface State {
    sock: Socket;
}

class GameScreen extends React.Component<Props, State> {
    private game: Game;
    private app: PIXI.Application;

    constructor(props: Props) {
        super(props);

        this.app = new PIXI.Application({width: 800, height: 600});
        this.game = new Game(this.app, props.sock);
    }

    render() {
        return this.app.view;
    }
}

export default GameScreen;
