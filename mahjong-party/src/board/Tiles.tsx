import Tile from "./Tile";
import React from "react";

export default function Tiles({ tiles, own }: { tiles: string[]; own?: boolean }) {
  const elements = tiles.map((tile, index) => (
    <Tile tile={tile} key={tile + index} own={own} />
  ));
  return <div className="rack">{elements}</div>;
}
