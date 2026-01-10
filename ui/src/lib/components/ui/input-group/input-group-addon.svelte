<script lang="ts">
    import { cn } from "$lib/utils.js";
    import type { HTMLAttributes } from "svelte/elements";

    let {
        class: className,
        children,
        align = "inline-end",
        ...rest
    } = $props<
        HTMLAttributes<HTMLDivElement> & {
            align?: "inline-start" | "inline-end";
        }
    >();
</script>

<div
    class={cn(
        "absolute flex items-center justify-center pointer-events-none px-3",
        align === "inline-start" ? "left-0" : "right-0",
        className,
    )}
    {...rest}
>
    {@render children?.()}
</div>
