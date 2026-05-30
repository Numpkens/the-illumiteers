<script>
  import { onMount } from 'svelte';
  import { loadVideos, videoList, isLoading, error } from './stores/videoStore.js';
  import HeroSection from './components/HeroSection.svelte';
  import GlimmerCard from './components/GlimmerCard.svelte';
  import BlogSection from './components/BlogSection.svelte';
  import SkeletonCard from './components/SkeletonCard.svelte';

  // Navigation tab switcher state: 'dashboard' or 'blog'
  let currentTab = 'dashboard';

  // Seed the video lists on mount
  onMount(() => {
    loadVideos();
  });

</script>

<svelte:head>
  {#if currentTab === 'dashboard'}
    <title>The Illumiteers | A Lorcana Community Hub</title>
  {:else}
    <title>Lore & Strategy Guides | The Illumiteers</title>
  {/if}
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <meta name="description" content="Explore community guides, meta analysis, and tactical insights from across the Inklands." />
</svelte:head>

<!-- Structured Wallpaper Matrix Layer (z-0) -->
<div aria-hidden="true" class="fixed inset-0 z-0 overflow-hidden pointer-events-none bg-[#F4EFE6] bg-cover bg-fixed select-none grid grid-cols-6 md:grid-cols-10 lg:grid-cols-12 gap-x-8 gap-y-16 p-4">
  {#each Array(72) as _, i}
    {#if i % 4 === 0}
      <div class="w-10 h-14 rounded-sm bg-teal-800/10 font-serif text-[10px] select-none text-black/5 flex flex-col justify-between p-1 shadow-sm transition-all duration-300 pointer-events-none {i % 3 === 0 ? 'rotate-6' : '-rotate-12'}">
        <span>i</span>
        <span class="text-right">i</span>
      </div>
    {:else if i % 4 === 1}
      <div class="w-10 h-14 rounded-sm bg-amber-800/10 font-serif text-[10px] select-none text-black/5 flex flex-col justify-between p-1 shadow-sm transition-all duration-300 pointer-events-none {i % 3 === 0 ? 'rotate-12' : '-rotate-6'}">
        <span>i</span>
        <span class="text-right">i</span>
      </div>
    {:else if i % 4 === 2}
      <div class="w-10 h-14 rounded-sm bg-purple-800/10 font-serif text-[10px] select-none text-black/5 flex flex-col justify-between p-1 shadow-sm transition-all duration-300 pointer-events-none {i % 3 === 0 ? 'rotate-6' : '-rotate-6'}">
        <span>i</span>
        <span class="text-right">i</span>
      </div>
    {:else}
      <div class="w-10 h-14 rounded-sm bg-rose-800/10 font-serif text-[10px] select-none text-black/5 flex flex-col justify-between p-1 shadow-sm transition-all duration-300 pointer-events-none {i % 3 === 0 ? 'rotate-12' : '-rotate-12'}">
        <span>i</span>
        <span class="text-right">i</span>
      </div>
    {/if}
  {/each}
</div>

<!-- Global parchment wrapper -->
<div class="min-h-screen w-full flex flex-col bg-transparent font-sans text-purple-950 relative">

  <!-- Header Banner Section (global branding) -->
  <header class="w-full min-w-full block relative z-50 bg-[#1E0F33] border-b-4 border-black">
    <!-- Star Sparkle indicators -->
    <div class="absolute top-2 right-4 text-white/20 text-xs select-none pointer-events-none font-serif">✦</div>
    <div class="absolute bottom-2 left-1/3 text-white/10 text-lg select-none pointer-events-none font-serif">✦</div>
    <div class="absolute top-4 left-10 text-white/15 text-sm select-none pointer-events-none font-serif">✦</div>

    <div class="w-full max-w-7xl mx-auto flex flex-col md:flex-row items-center justify-between px-4 sm:px-6 lg:px-8 py-4 md:py-0 min-h-[4rem] md:h-20 gap-4">
      <!-- Title & Branding -->
      <div class="flex items-center space-x-3">
        <img
          src="/assets/logo.png"
          alt="The Illumiteers Logo"
          class="w-10 h-10 object-contain rounded-full border border-white/20 bg-black/40"
        />
        <div class="flex flex-col text-left">
          <h1 class="font-serif font-black text-[clamp(1.15rem,4.5vw,1.5rem)] text-white tracking-tight uppercase leading-[1.3]">
            The Illumiteers
          </h1>
          <span class="text-[10px] font-mono tracking-widest text-[#FAF5E6] uppercase mt-1" style="text-shadow: 1px 1px 2px rgba(0, 0, 0, 0.8);">
            A Lorcana Community Hub
          </span>
        </div>
      </div>

      <!-- Navigation Tabs -->
      <nav aria-label="Main Navigation" class="flex flex-col sm:flex-row items-center justify-center gap-3 sm:gap-6 md:gap-10 mx-auto px-4">
        <button
          type="button"
          on:click={() => currentTab = 'dashboard'}
          class={currentTab === 'dashboard'
            ? 'font-serif font-bold text-xs md:text-sm uppercase tracking-widest text-white transition-colors duration-200 relative py-2 border-b-2 border-[#F5D583] pb-1 shadow-[0_1px_4px_rgba(245,213,131,0.4)]'
            : 'font-serif font-bold text-xs md:text-sm uppercase tracking-widest text-[#EFE6D5]/90 hover:text-white transition-colors duration-200 relative py-2'}
        >
          Community Dashboard
        </button>
        <button
          type="button"
          on:click={() => currentTab = 'blog'}
          class={currentTab === 'blog'
            ? 'font-serif font-bold text-xs md:text-sm uppercase tracking-widest text-white transition-colors duration-200 relative py-2 border-b-2 border-[#F5D583] pb-1 shadow-[0_1px_4px_rgba(245,213,131,0.4)]'
            : 'font-serif font-bold text-xs md:text-sm uppercase tracking-widest text-[#EFE6D5]/90 hover:text-white transition-colors duration-200 relative py-2'}
        >
          Lore & Strategy Guides
        </button>
      </nav>
    </div>
  </header>

  <!-- Container Wrapper -->
  <main class="max-w-7xl w-full mx-auto px-4 sm:px-6 lg:px-8 pt-8 flex-grow">
    
    {#if currentTab === 'dashboard'}
      <!-- Hero Splitting Zone (Top Area) -->
      <section id="top-workspace" class="mb-12">
        {#if $isLoading}
          <div class="w-full min-h-[320px] md:min-h-[440px] flex flex-col items-center justify-center bg-transparent border-4 border-black rounded-2xl shadow-[6px_6px_0px_0px_rgba(0,0,0,1)] p-8">
            <div class="w-16 h-16 bg-purple-100 border-2 border-black rounded-full flex items-center justify-center animate-spin mb-4">
              <span class="text-purple-950 font-black text-2xl">✦</span>
            </div>
            <h3 class="font-serif font-bold text-[clamp(1rem,4.5vw,1.125rem)] text-purple-950 uppercase leading-[1.3]">Summoning Community Hub...</h3>
            <p class="text-xs font-sans text-gray-500 mt-1">Drawing video glimmers and environment keys</p>
          </div>
        {:else}
          <HeroSection />
        {/if}
      </section>

      <!-- Recent Glimmers Gallery Zone (Bottom Area) -->
      <section aria-label="Recent Glimmers Video Feed" class="border-t-4 border-dashed border-purple-900/30 pt-10">
        <div class="mb-8">
          <div class="flex items-center space-x-2">
            <div class="w-3 h-3 bg-purple-900 rotate-45"></div>
            <h2 class="font-serif font-black text-[clamp(1.35rem,6vw,1.875rem)] text-purple-950 uppercase tracking-tight leading-[1.3]">
              Recent Glimmers
            </h2>
          </div>
          <p class="font-sans text-xs sm:text-sm text-gray-600 mt-1">
            Select any lore-bound video card below to load it into the top active media console.
          </p>
        </div>

        {#if $isLoading}
          <!-- Grid list skeleton cards -->
          <div class="grid grid-cols-1 sm:grid-cols-[repeat(auto-fit,minmax(280px,1fr))] gap-6">
            {#each Array(4) as _}
              <SkeletonCard type="video" />
            {/each}
          </div>
        {:else if $error && $videoList.length === 0}
          <div class="p-6 bg-red-50 border-4 border-black rounded-2xl shadow-[6px_6px_0px_0px_rgba(0,0,0,1)] text-center">
            <h3 class="font-serif font-bold text-red-700 text-[clamp(1rem,4.5vw,1.125rem)] uppercase leading-[1.3]">Summoning Failed</h3>
            <p class="text-sm font-sans text-red-600 mt-1">{$error}</p>
          </div>
        {:else}
          <!-- Grid list mapping cards -->
          <div class="grid grid-cols-1 sm:grid-cols-[repeat(auto-fit,minmax(280px,1fr))] gap-6">
            {#each $videoList as video, index (video.id)}
              <GlimmerCard {video} {index} />
            {/each}
          </div>
        {/if}
      </section>
    {:else}
      <!-- Lore & Strategy Section -->
      <section aria-label="Lore and Strategy Guides" class="mb-12">
        <div class="mb-8">
          <div class="flex items-center space-x-2">
            <div class="w-3 h-3 bg-purple-900 rotate-45"></div>
            <h2 class="font-serif font-black text-[clamp(1.35rem,6vw,1.875rem)] text-purple-950 uppercase tracking-tight leading-[1.3]">
              Lore & Strategy Guides
            </h2>
          </div>
          <p class="text-xs md:text-sm text-[#4A3565]/80 font-normal mt-1 tracking-wide">
            Explore community guides, meta analysis, and tactical insights from across the Inklands.
          </p>
        </div>

        <BlogSection />
      </section>
    {/if}

  </main>
</div>
