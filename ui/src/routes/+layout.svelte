<script lang="ts">
    import { ModeWatcher } from "mode-watcher";
    import Header from "$lib/components/Header.svelte";
    import { Toaster } from "$lib/components/ui/sonner";
    import "../app.css";

    import { userState } from "$lib/user_state.svelte";
    import * as Tooltip from "$lib/components/ui/tooltip";

    let { data, children } = $props<{
        data: { user: typeof userState.current } | undefined;
        children: any;
    }>();

    // Hydrate shared user state from server data
    $effect(() => {
        if (data?.user) {
            userState.setLoaded(data.user, null);
        }
    });
</script>

<svelte:head>
    <title>GH Secret Broker</title>
</svelte:head>

<ModeWatcher />
<Tooltip.Provider>
    <div class="min-h-screen bg-background">
        <Header />
        <main>
            {#if userState.loading}
                <div
                    class="flex items-center justify-center py-16 text-muted-foreground"
                >
                    Loading your data...
                </div>
            {:else}
                {@render children()}
            {/if}
        </main>
        <Toaster
            position="top-center"
            richColors
            expand
            toastOptions={{ class: "max-w-xl w-full" }}
        />
    </div>
</Tooltip.Provider>
