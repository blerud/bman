import * as React from 'react';
import * as PIXI from 'pixi.js';
import Game, {InitInfo} from "./game";
import Socket from "./socket";

interface Props {
    sock: Socket;
    initInfo: InitInfo;
}

class GameScreen extends React.Component<Props, {}> {
    private game: Game;
    private app: PIXI.Application;
    private initInfo: InitInfo;
    private sock: Socket;

    constructor(props: Props) {
        super(props);

        this.app = new PIXI.Application({width: 800, height: 600});
        this.game = new Game(this.app, props.sock, props.initInfo);
        this.sock = props.sock;
        this.initInfo = props.initInfo;
    }

    componentDidMount() {
        let container = document.getElementById("container");
        container.appendChild(this.app.view);

        this.game.renderBackground();
        this.game.renderPlayer();
    }

    render() {
        return (
            <div id="container">
                <h1>server id: {this.props.initInfo.gameid}</h1>
            </div>
        );
    }
}

export default GameScreen;
