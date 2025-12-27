import React from "react";

export function Segmented({
                              left,
                              right,
                              value,
                              onChange,
                          }: {
    left: string;
    right: string;
    value: "left" | "right";
    onChange: (v: "left" | "right") => void;
}) {
    return (
        <div className="flex gap-4">
            <button
                type="button"
                onClick={() => onChange("left")}
                className={[
                    "flex-1 rounded-full px-6 py-4 text-lg border",
                    value === "left" ? "bg-primary text-white border-primary" : "bg-white text-ink border-border",
                ].join(" ")}
            >
                {left}
            </button>
            <button
                type="button"
                onClick={() => onChange("right")}
                className={[
                    "flex-1 rounded-full px-6 py-4 text-lg border",
                    value === "right" ? "bg-primary text-white border-primary" : "bg-white text-ink border-border",
                ].join(" ")}
            >
                {right}
            </button>
        </div>
    );
}
