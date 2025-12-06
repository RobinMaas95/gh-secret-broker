<script lang="ts">
    import { untrack } from "svelte";
    import * as Card from "$lib/components/ui/card";
    import * as Avatar from "$lib/components/ui/avatar";
    import { Button } from "$lib/components/ui/button";
    import { Label } from "$lib/components/ui/label";
    import { Separator } from "$lib/components/ui/separator";

    interface User {
        AvatarURL: string;
        Name: string;
        NickName: string;
        Email: string;
        Location: string;
        Description: string;
        UserID: string;
        Provider: string;
    }

    let user = $state<User | null>(null);
    let error = $state<string | null>(null);

    async function loadUser() {
        try {
            const res = await fetch("/api/me");
            if (!res.ok) {
                throw new Error("Unauthorized");
            }
            user = await res.json();
        } catch (e) {
            error = (e as Error).message;
        }
    }

    $effect(() => {
        untrack(() => {
            loadUser();
        });
    });
</script>

<div class="flex min-h-screen items-center justify-center p-4">
    <Card.Root class="w-full max-w-lg">
        {#if error}
            <Card.Content class="pt-6">
                <div class="text-destructive text-center mb-4">{error}</div>
                <div class="text-center">
                    <Button variant="link" href="/">Go back to login</Button>
                </div>
            </Card.Content>
        {:else if user}
            <Card.Header class="flex flex-row items-center gap-4">
                <Avatar.Root class="h-20 w-20">
                    <Avatar.Image src={user.AvatarURL} alt={user.Name} />
                    <Avatar.Fallback
                        >{user.Name.slice(0, 2).toUpperCase()}</Avatar.Fallback
                    >
                </Avatar.Root>
                <div class="flex flex-col gap-1">
                    <Card.Title class="text-2xl">{user.Name}</Card.Title>
                    <Card.Description>@{user.NickName}</Card.Description>
                </div>
            </Card.Header>
            <Separator />
            <Card.Content class="grid gap-4 py-6">
                <div class="grid grid-cols-[100px_1fr] items-center gap-4">
                    <Label class="text-right text-muted-foreground">Email</Label
                    >
                    <div class="font-medium">{user.Email}</div>
                </div>
                <div class="grid grid-cols-[100px_1fr] items-center gap-4">
                    <Label class="text-right text-muted-foreground"
                        >Location</Label
                    >
                    <div class="font-medium">
                        {user.Location || "Not specified"}
                    </div>
                </div>
                <div class="grid grid-cols-[100px_1fr] items-start gap-4">
                    <Label class="text-right text-muted-foreground pt-1"
                        >Bio</Label
                    >
                    <div class="font-medium">
                        {user.Description || "No bio"}
                    </div>
                </div>
                <div class="grid grid-cols-[100px_1fr] items-center gap-4">
                    <Label class="text-right text-muted-foreground"
                        >User ID</Label
                    >
                    <div class="font-mono text-xs text-muted-foreground">
                        {user.UserID}
                    </div>
                </div>
            </Card.Content>
            <Card.Footer class="justify-center">
                <Button variant="destructive" href="/logout/{user.Provider}"
                    >Logout</Button
                >
            </Card.Footer>
        {:else}
            <Card.Content class="pt-6 text-center text-muted-foreground">
                Loading profile...
            </Card.Content>
        {/if}
    </Card.Root>
</div>
