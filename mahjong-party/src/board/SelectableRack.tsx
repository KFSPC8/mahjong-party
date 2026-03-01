import React from "react";

import { TileClickCallback } from "./types";
import Tile from "./Tile";

export default function SelectableRack({
  tiles,
  onTileClick,
  selecting,
  selected,
  own,
}: {
  tiles: Array<{ tile: string; id: number }>;
  selecting?: boolean;
  selected?: number[];
  onTileClick?: TileClickCallback;
  own?: boolean;
}) {
  const elements = tiles.map(({ tile, id }, index) => (
    <Tile
      tile={tile}
      key={id}
      selected={selected?.includes(index)}
      onClick={onTileClick ? onTileClick(tile, index) : undefined}
      own={own}
    />
  ));
  return (
    <div
      className={selecting ? "rack selectable selecting" : "rack selectable"}>
      {elements}
    </div>
  );
}
