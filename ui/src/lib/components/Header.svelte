<script lang="ts">
    import { userState } from "$lib/user_state.svelte";
    import * as Avatar from "$lib/components/ui/avatar";
    import * as Button from "$lib/components/ui/button";
    import ModeToggle from "$lib/components/ModeToggle.svelte";
</script>

<header
    class="flex h-14 items-center gap-4 border-b bg-muted/40 px-6 lg:h-[60px]"
>
    <div class="flex-1">
        <a href="/" class="hover:opacity-80 transition-opacity">
            <h1 class="font-semibold text-lg">GH Secret Broker</h1>
        </a>
    </div>
    <div class="flex items-center gap-4">
        <ModeToggle />
        {#if userState.current}
            <a href="/profile">
                <Avatar.Root
                    class="h-9 w-9 cursor-pointer hover:opacity-80 transition-opacity"
                >
                    <Avatar.Image
                        src={userState.current.AvatarURL}
                        alt={userState.current.Name}
                    />
                    <Avatar.Fallback>
                        {userState.current.Name.slice(0, 2).toUpperCase()}
                    </Avatar.Fallback>
                </Avatar.Root>
            </a>
        {/if}
    </div>
</header>
