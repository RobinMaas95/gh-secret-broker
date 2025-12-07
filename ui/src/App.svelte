<script lang="ts">
  import { ModeWatcher } from "mode-watcher";
  import Login from "./lib/Login.svelte";
  import User from "./lib/User.svelte"; // This is now the Dashboard
  import Profile from "./lib/Profile.svelte";
  import Header from "$lib/components/Header.svelte";

  let page = $state<"login" | "dashboard" | "profile">("login");

  function handleHashChange() {
    const hash = window.location.hash;
    if (hash === "#/dashboard") {
      page = "dashboard";
    } else if (hash === "#/profile") {
      page = "profile";
    } else {
      page = "login";
    }
  }

  // Initialize on mount
  handleHashChange();
</script>

<ModeWatcher />
<svelte:window onhashchange={handleHashChange} />

<div class="min-h-screen bg-background">
  <Header />
  <main>
    {#if page === "login"}
      <Login />
    {:else if page === "dashboard"}
      <User />
    {:else if page === "profile"}
      <Profile />
    {/if}
  </main>
</div>

<style>
</style>
