import React, { MouseEventHandler, CSSProperties } from "react";

export default function Tile({
  tile,
  onClick,
  selected,
  inline,
  style,
  own,
}: {
  tile: string;
  onClick?: MouseEventHandler;
  selected?: boolean;
  inline?: boolean;
  style?: CSSProperties;
  // when true, use the bottom-board tile styles
  own?: boolean;
}) {
  const base = own ? "tile-own" : "tile";
  const classes = [base];
  selected && classes.push("selected");
  inline && classes.push("inline");
  return (
    <span
      className={classes.join(" ")}
      data-tile={tile}
      onClick={onClick}
      style={style}
    />
  );
}
