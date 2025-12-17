import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
    return twMerge(clsx(inputs));
}

export type WithElementRef<T = object> = T & { ref?: HTMLElement | null };

// Relaxing types for Svelte 5 compat
export type WithoutChild<T> = T;
export type WithoutChildren<T> = T;
export type WithoutChildrenOrChild<T> = T;
