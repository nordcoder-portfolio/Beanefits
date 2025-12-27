import React from "react";

type Variant = "primary" | "outline" | "soft";

export function Button(
    props: React.ButtonHTMLAttributes<HTMLButtonElement> & { variant?: Variant; fullWidth?: boolean }
) {
    const { variant = "primary", fullWidth, className = "", ...rest } = props;

    const base =
        "inline-flex items-center justify-center rounded-full px-6 py-4 text-base transition active:scale-[0.99] disabled:opacity-50 disabled:pointer-events-none";
    const w = fullWidth ? "w-full" : "";
    const v =
        variant === "primary"
            ? "bg-primary2 text-white"
            : variant === "outline"
                ? "bg-white text-ink border border-border"
                : "bg-peach text-ink border border-border";

    return <button className={`${base} ${w} ${v} ${className}`} {...rest} />;
}
