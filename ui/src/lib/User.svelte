<script lang="ts">
    import { untrack } from "svelte";

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

<div class="container">
    <h1>User Profile</h1>
    <div class="card">
        {#if error}
            <p class="error">{error}</p>
            <p><a href="/">Go back to login</a></p>
        {:else if user}
            <div class="profile-header">
                <img src={user.AvatarURL} class="avatar" alt="Avatar" />
                <div>
                    <h2>{user.Name}</h2>
                    <p>{user.NickName}</p>
                </div>
            </div>

            <div class="info-grid">
                <div class="label">Email:</div>
                <div>{user.Email}</div>
                <div class="label">Location:</div>
                <div>{user.Location}</div>
                <div class="label">Description:</div>
                <div>{user.Description}</div>
                <div class="label">UserID:</div>
                <div>{user.UserID}</div>
            </div>

            <div class="logout-link">
                <a href="/logout/{user.Provider}">Logout</a>
            </div>
        {:else}
            <p>Loading profile...</p>
        {/if}
    </div>
</div>

<style>
    .container {
        max-width: 800px;
        margin: 0 auto;
        padding: 2rem;
    }
    h1 {
        font-size: 2.5em;
        line-height: 1.1;
        text-align: center;
    }
    .card {
        padding: 2em;
        background-color: #1a1a1a;
        border-radius: 8px;
        margin-top: 1em;
        display: flex;
        flex-direction: column;
        gap: 1em;
    }
    .profile-header {
        display: flex;
        align-items: center;
        gap: 1em;
        border-bottom: 1px solid #333;
        padding-bottom: 1em;
        margin-bottom: 1em;
    }
    .avatar {
        width: 100px;
        height: 100px;
        border-radius: 50%;
        object-fit: cover;
    }
    .info-grid {
        display: grid;
        grid-template-columns: auto 1fr;
        gap: 0.5em 1em;
    }
    .label {
        font-weight: bold;
        color: #aaa;
    }
    .logout-link {
        text-align: center;
        margin-top: 2em;
    }
    .error {
        color: red;
        text-align: center;
    }
</style>
