<script lang="ts">
  import { ModeWatcher } from "mode-watcher";
  import Login from "./lib/Login.svelte";
  import User from "./lib/User.svelte";
  import ModeToggle from "$lib/components/ModeToggle.svelte";

  let page = $state<"login" | "user">("login");

  function handleHashChange() {
    const hash = window.location.hash;
    if (hash === "#/userpage") {
      page = "user";
    } else {
      page = "login";
    }
  }

  // Initialize on mount
  handleHashChange();
</script>

<ModeWatcher />
<svelte:window onhashchange={handleHashChange} />

<div class="absolute right-4 top-4">
  <ModeToggle />
</div>

<main>
  {#if page === "login"}
    <Login />
  {:else if page === "user"}
    <User />
  {/if}
</main>

<style>
</style>
