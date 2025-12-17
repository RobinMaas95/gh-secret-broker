<script lang="ts">
    import { cn } from "$lib/utils.js";
    import type { HTMLAttributes } from "svelte/elements";

    let {
        class: className,
        children,
        ...rest
    }: HTMLAttributes<HTMLDivElement> = $props();
</script>

<div class={cn("flex relative items-center w-full", className)} {...rest}>
    {@render children?.()}
</div>
