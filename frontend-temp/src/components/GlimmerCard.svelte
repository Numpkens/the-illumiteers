<script>
  import { activeVideoId } from '../stores/videoStore.js';
  export let video;
  export let index = 0;

  function handleSelect() {
    activeVideoId.set(video.id);
    // Smooth scroll up to player
    const playerEl = document.getElementById('top-workspace');
    if (playerEl) {
      playerEl.scrollIntoView({ behavior: 'smooth' });
    }
  }

  function handleKeyDown(e) {
    if (e.key === 'Enter' || e.key === ' ') {
      e.preventDefault();
      handleSelect();
    }
  }

  $: isActive = $activeVideoId === video.id;

  // Define the 6 official Lorcana ink colors schema mapping
  const themes = [
    { bgText: 'bg-[#F5D583] text-[#5C4017]', border: 'border-[#C99C33]' }, // Amber
    { bgText: 'bg-[#D1B3DB] text-[#42214D]', border: 'border-[#835294]' }, // Amethyst
    { bgText: 'bg-[#A3D4BC] text-[#143B27]', border: 'border-[#378C62]' }, // Emerald
    { bgText: 'bg-[#EBB0B5] text-[#521318]', border: 'border-[#B8323D]' }, // Ruby
    { bgText: 'bg-[#A3C7EB] text-[#102D4C]', border: 'border-[#336FA8]' }, // Sapphire
    { bgText: 'bg-[#C2CBD1] text-[#333E45]', border: 'border-[#6D7B85]' }  // Steel
  ];

  $: theme = themes[index % 6];
</script>

<!-- svelte-ignore a11y_no_noninteractive_element_to_interactive_role -->
<article
  tabindex="0"
  role="button"
  on:click={handleSelect}
  on:keydown={handleKeyDown}
  class="aspect-[2.5/3.5] w-full max-w-[360px] flex flex-col justify-between rounded-2xl p-2 bg-black border-4 border-black shadow-xl select-none text-white relative hover:-translate-y-1 hover:scale-[1.02] hover:shadow-2xl transition-transform duration-300 ease-out group active:scale-[0.98] text-left cursor-pointer focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[#F5D583]
    {isActive ? 'ring-4 ring-amber-500' : ''}"
>
  <!-- Dynamic 2px Inner Border Line just inside black frame -->
  <div class="absolute inset-1 border-2 rounded-[12px] pointer-events-none z-10 {theme.border}"></div>

  <!-- Holographic Foil Hover Sweep -->
  <div class="absolute inset-0 rounded-[12px] bg-gradient-to-tr from-transparent via-white/5 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300 ease-out pointer-events-none z-30"></div>

  <!-- THE TOP ZONE: ART & STATS -->
  <div class="w-full aspect-[4/3] rounded-t-lg bg-black overflow-hidden relative shadow-md border-b-2 border-black z-0">
    {#if video.thumbnailUrl}
      <img
        src={video.thumbnailUrl}
        alt="Video thumbnail for {video.title}"
        class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300 rounded-t-lg border-b-2 border-black"
      />
    {:else}
      <div class="w-full h-full flex flex-col items-center justify-center text-[10px] text-gray-500 bg-neutral-800 rounded-t-lg border-b-2 border-black">
        <span class="text-lg">✦</span>
        <span>NO IMAGE GLIMMER</span>
      </div>
    {/if}
    <!-- Play Overlay on Hover -->
    <div class="absolute inset-0 bg-black/20 opacity-0 group-hover:opacity-100 flex items-center justify-center transition-opacity duration-200">
      <div class="w-9 h-9 bg-amber-500 border-2 border-black rounded-full flex items-center justify-center shadow-md">
        <span class="text-black text-xs ml-0.5">▶</span>
      </div>
    </div>
  </div>

  <!-- THE MIDDLE ZONE: TITLE & IDENTIFIER (Banner Stripe) -->
  <h3 class="w-full border-b border-black z-10 {theme.bgText} px-3 py-1 font-serif text-sm font-bold tracking-tight line-clamp-1 text-shadow-sm">
    {video.title}
  </h3>
  <div class="w-full bg-[#EFE6D5] text-[#3A2E1A] text-[10px] font-mono tracking-widest uppercase px-3 py-0.5 border-b border-black font-semibold z-10">
    GLIMMER • COMMUNITY • VIDEO
  </div>

  <!-- THE BOTTOM ZONE: PARCHMENT TEXT PANEL -->
  <div class="bg-[#F1E8D9] p-3 flex-grow flex flex-col justify-between rounded-b-lg text-[#3A2E1A] z-10">
    <p class="font-serif text-xs leading-relaxed opacity-95 line-clamp-3 italic text-left">
      {video.description || 'No detailed description was gathered for this video feed in Liam\'s logs.'}
    </p>
  </div>
</article>
