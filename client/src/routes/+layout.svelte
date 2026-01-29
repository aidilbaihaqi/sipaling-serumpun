<script lang="ts">
  import '../app.css';
  import Navbar from '$lib/components/Navbar.svelte';
  import Footer from '$lib/components/Footer.svelte';
</script>

<div class="app-wrapper">
  <Navbar />
  
  <main class="main-content">
    <slot />
  </main>
  
  <Footer />
</div>

<style>
  .app-wrapper {
    display: flex;
    flex-direction: column;
    min-height: 100vh;
  }
  
  .main-content {
    flex: 1;
    display: flex;
    flex-direction: column;
  }
</style>
