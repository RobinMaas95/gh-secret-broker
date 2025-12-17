<script lang="ts">
    import { Button } from "$lib/components/ui/button";
    import { cn } from "$lib/utils.js";
    import type { HTMLButtonAttributes } from "svelte/elements";

    let {
        class: className,
        children,
        ...rest
    }: HTMLButtonAttributes = $props();
</script>

<Button
    class={cn("pointer-events-auto h-7 w-7", className)}
    variant="ghost"
    size="icon"
    {...rest as any}
>
    {@render children?.()}
</Button>
