import React, { FunctionComponent, useEffect, useRef } from "react";
import Tile from "./Tile";
import { useDrop } from "react-dnd";
import { DraggedTile } from "./types";

const Discards: FunctionComponent<{
  discards: string[];
  canDiscard: boolean;
  discardTile: (tile: string, index: number) => void;
}> = ({ discards, canDiscard, discardTile }) => {
  const prevCount = useRef(discards.length);
  const audioCtxRef = useRef<AudioContext | null>(null);

  const playDiscardSound = () => {
    try {
      const AudioCtx = (window as any).AudioContext || (window as any).webkitAudioContext;
      if (!audioCtxRef.current) {
        audioCtxRef.current = new AudioCtx();
      }
  const ctx = audioCtxRef.current!;
      // Some browsers require resume() on a user gesture
      if (ctx.state === "suspended") {
        ctx.resume();
      }
      const now = ctx.currentTime;
      const o = ctx.createOscillator();
      const g = ctx.createGain();
      o.type = "square";
      o.frequency.setValueAtTime(800, now);
      g.gain.setValueAtTime(0.0001, now);
      g.gain.exponentialRampToValueAtTime(0.4, now + 0.01);
      g.gain.exponentialRampToValueAtTime(0.0001, now + 0.15);
      o.connect(g);
      g.connect(ctx.destination);
      o.start(now);
      o.stop(now + 0.18);
    } catch (e) {
      // ignore audio errors
      // console.warn("Audio not available", e);
    }
  };

  useEffect(() => {
    if (discards.length > prevCount.current) {
      playDiscardSound();
    }
    prevCount.current = discards.length;
  }, [discards]);
  const [{ isOver }, drop] = useDrop({
    accept: "tile",
    collect: (monitor) => ({
      isOver: monitor.isOver(),
    }),
    drop: ({ tile, index }: DraggedTile) => {
      if (canDiscard) {
        discardTile(tile, index);
      }
    },
  });
  return (
    <div
      ref={drop}
      className="centre"
      style={{
        backgroundColor: canDiscard && isOver ? "#f1f1f1" : "#205539",
      }}>
      {discards.map((tile, index) => (
        <Tile tile={tile} key={tile + index} />
      ))}
    </div>
  );
};

export default Discards;
