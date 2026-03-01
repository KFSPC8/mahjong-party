import React, { MouseEventHandler, CSSProperties } from "react";

export default function Tile({
  tile,
  onClick,
  selected,
  inline,
  style,
}: {
  tile: string;
  onClick?: MouseEventHandler;
  selected?: boolean;
  inline?: boolean;
  style?: CSSProperties;
}) {
  const classes = ["tile"];
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
