<script>
  import { onMount } from 'svelte';
  import { marked } from 'marked';
  import { blogList, isBlogLoading, blogError, loadBlogPosts } from '../stores/blogStore.js';
  import SkeletonCard from './SkeletonCard.svelte';

  // Seed store data on mount
  onMount(() => {
    loadBlogPosts();
  });

  let selectedPost = null;

  // Define the 6 official Lorcana ink colors schema mapping
  const themes = [
    { bgText: 'bg-[#F5D583] text-[#5C4017]', border: 'border-[#C99C33]' }, // Amber
    { bgText: 'bg-[#D1B3DB] text-[#42214D]', border: 'border-[#835294]' }, // Amethyst
    { bgText: 'bg-[#A3D4BC] text-[#143B27]', border: 'border-[#378C62]' }, // Emerald
    { bgText: 'bg-[#EBB0B5] text-[#521318]', border: 'border-[#B8323D]' }, // Ruby
    { bgText: 'bg-[#A3C7EB] text-[#102D4C]', border: 'border-[#336FA8]' }, // Sapphire
    { bgText: 'bg-[#C2CBD1] text-[#333E45]', border: 'border-[#6D7B85]' }  // Steel
  ];

  function getSnippet(content) {
    if (!content) return "";
    let clean = content
      .replace(/[#*`_-]/g, "")
      .replace(/\[([^\]]+)\]\([^)]+\)/g, "$1")
      .trim();
    return clean;
  }

  function handleKeyDown(e, post) {
    if (e.key === 'Enter' || e.key === ' ') {
      e.preventDefault();
      selectedPost = post;
    }
  }
</script>

<div class="w-full">
  {#if $isBlogLoading}
    <div class="max-w-4xl mx-auto w-full flex flex-col items-center justify-center bg-transparent border-4 border-black rounded-2xl shadow-[6px_6px_0px_0px_rgba(0,0,0,1)] p-12 text-center mb-6">
      <div class="w-14 h-14 bg-purple-100 border-2 border-black rounded-full flex items-center justify-center animate-spin mb-4">
        <span class="text-purple-950 font-black text-xl">✦</span>
      </div>
      <h3 class="font-serif font-bold text-[clamp(1rem,4.5vw,1.125rem)] text-purple-950 uppercase leading-[1.3]">Unrolling Strategy Scrolls...</h3>
      <p class="text-xs font-sans text-gray-500 mt-1">Drawing guides from the Illumiteers vault</p>
    </div>
    <div class="grid grid-cols-1 sm:grid-cols-[repeat(auto-fit,minmax(280px,1fr))] gap-6 p-4 max-w-7xl mx-auto relative z-10 bg-transparent w-full">
      {#each Array(3) as _}
        <SkeletonCard type="blog" />
      {/each}
    </div>
  {:else if $blogError && $blogList.length === 0}
    <div class="max-w-4xl mx-auto p-8 bg-red-50 border-4 border-black rounded-2xl shadow-[6px_6px_0px_0px_rgba(0,0,0,1)] text-center">
      <h3 class="font-serif font-bold text-red-700 text-[clamp(1rem,4.5vw,1.125rem)] uppercase leading-[1.3]">Summoning Failed</h3>
      <p class="text-sm font-sans text-red-600 mt-1">{$blogError}</p>
    </div>
  {:else}
    <!-- Strategy Guide Collectible Card Grid -->
    <div class="grid grid-cols-1 sm:grid-cols-[repeat(auto-fit,minmax(280px,1fr))] gap-6 p-4 max-w-7xl mx-auto relative z-10 bg-transparent w-full">
      {#each $blogList as post, index}
        {@const theme = themes[index % 6]}
        <!-- svelte-ignore a11y_no_noninteractive_element_to_interactive_role -->
        <article
          tabindex="0"
          role="button"
          on:click={() => selectedPost = post}
          on:keydown={(e) => handleKeyDown(e, post)}
          class="w-full min-h-[280px] h-auto flex flex-col justify-between rounded-xl p-1.5 pb-3 bg-black border-2 border-black shadow-lg text-white relative hover:scale-[1.02] transition-transform duration-300 ease-out group active:scale-[0.98] text-left cursor-pointer focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[#F5D583]"
        >
          <!-- Dynamic 2px Inner Border Line just inside black frame -->
          <div class="absolute inset-1 border-2 rounded-[10px] pointer-events-none z-10 {theme.border}"></div>

          <!-- Holographic Foil Hover Sweep -->
          <div class="absolute inset-0 rounded-[10px] bg-gradient-to-tr from-transparent via-white/5 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300 ease-out pointer-events-none z-30"></div>

          <!-- The Top Banner Stripe (Ink Tone) -->
          <h3 class="w-full {theme.bgText.split(' ')[0]} px-3 py-2 font-serif text-[clamp(0.7rem,3.5vw,0.75rem)] font-bold line-clamp-2 border-b border-black rounded-t-md z-10 text-[#1C1C1C] leading-[1.3]">
            {post.title}
          </h3>

          <!-- The Sub-Category Tag Line -->
          <div class="w-full bg-[#EFE6D5] text-[#3A2E1A] text-[9px] font-mono tracking-wider uppercase px-3 py-0.5 border-b border-black font-semibold z-10">
            {post.category}
          </div>

          <!-- The Aged Parchment Text Block -->
          <div class="bg-[#F1E8D9] p-3 flex-grow flex flex-col justify-between rounded-b-md text-[#3A2E1A] z-10">
            <p class="font-serif text-[11px] leading-relaxed opacity-95 line-clamp-5 italic text-left">
              {getSnippet(post.content)}
            </p>

            <div class="flex justify-between items-center mt-1 pt-1 border-t border-[#3A2E1A]/10">
              <span class="text-[9px] font-mono opacity-65">{post.date}</span>
              <span class="text-[10px] font-serif font-bold text-purple-950 group-hover:text-purple-700 transition-colors uppercase tracking-wider flex items-center gap-1">
                Read Insights →
              </span>
            </div>
          </div>
        </article>
      {/each}
    </div>
  {/if}
</div>

<!-- Premium Full-Screen Reading Modal -->
{#if selectedPost}
  {@const theme = themes[$blogList.indexOf(selectedPost) % 6]}
  <!-- Modal Backdrop -->
  <div class="fixed inset-0 bg-black/60 z-50 flex items-center justify-center p-4 backdrop-blur-sm">
    <!-- Modal Container -->
    <div class="bg-[#F4EFE6] border-4 border-black rounded-2xl max-w-2xl w-full max-h-[85vh] flex flex-col shadow-[8px_8px_0px_0px_rgba(0,0,0,1)] relative overflow-hidden">
      <!-- Inner Border Trim -->
      <div class="absolute inset-1 border-2 border-[#3A2E1A]/20 rounded-[12px] pointer-events-none z-10"></div>

      <!-- Header Banner (styled like Title Stripe using deterministic color theme) -->
      <div class="w-full border-b-2 border-black z-20 {theme.bgText} px-6 py-4 font-serif text-base sm:text-lg font-bold flex items-center justify-between text-shadow-sm">
        <div class="flex flex-col text-left">
          <span class="text-[10px] uppercase tracking-widest opacity-80 leading-normal mb-1">{selectedPost.category}</span>
          <span class="leading-[1.3]">{selectedPost.title}</span>
        </div>
        <button
          type="button"
          on:click={() => selectedPost = null}
          class="w-8 h-8 rounded-full border border-current flex items-center justify-center hover:bg-black/10 transition-colors shrink-0 ml-4"
          aria-label="Close modal"
        >
          ✕
        </button>
      </div>

      <!-- Sub-Type Bar -->
      <div class="w-full bg-[#EFE6D5] text-[#3A2E1A] text-[10px] font-mono tracking-widest uppercase px-6 py-1.5 border-b border-black font-semibold z-20 text-left">
        ILLUMITEER INSIGHTS • ARCHIVE DATE: {selectedPost.date}
      </div>

      <!-- Scrollable Parchment Content -->
      <div class="p-6 overflow-y-auto z-20 flex-grow font-serif text-[#3A2E1A] bg-[#FCF9F2] prose prose-amber max-w-none text-left">
        {@html marked.parse(selectedPost.content)}
      </div>

      <!-- Footer close -->
      <div class="px-6 py-3 border-t border-black/10 bg-[#EFE6D5] z-20 flex justify-end">
        <button
          type="button"
          on:click={() => selectedPost = null}
          class="px-4 py-2 border-2 border-black bg-white hover:bg-neutral-50 text-black text-xs font-serif font-bold uppercase tracking-wider rounded shadow-[2px_2px_0px_0px_rgba(0,0,0,1)] active:translate-y-[1px] active:shadow-[1px_1px_0px_0px_rgba(0,0,0,1)] transition-all"
        >
          Close Scroll
        </button>
      </div>
    </div>
  </div>
{/if}
