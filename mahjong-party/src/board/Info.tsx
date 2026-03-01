import React, { FunctionComponent } from "react";
import { Player, Round } from "../mahjong";
import Tile from "./Tile";

const WINDS = ["East", "South", "West", "North"];
const WIND_TILES = ["40东风", "41南风", "42西风", "43北风"];

const Info: FunctionComponent<{ players: Player[]; round: Round }> = ({
  players,
  round,
  children,
}) => {
  const { wind, dealer, draws_left } = round;
  return (
    <div className="info">
      <span className="infobox">
        Prevailing wind: {WINDS[wind]} 
        <Tile
          tile={WIND_TILES[wind]}
          inline={true}
          style={{
            width: "20px",
            height: "22px",
            marginLeft: "6px",
            verticalAlign: "-20%",
          }}
        />
      </span>
      <span className="infobox">Dealer: {players[dealer].name}</span>
      <span className="infobox">Draws left: {draws_left}</span>
      <span style={{ marginRight: "auto" }} />
      {children}
    </div>
  );
};

export default Info;
