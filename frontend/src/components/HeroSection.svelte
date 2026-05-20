<script>
  import { activeVideoId, videoList } from '../stores/videoStore.js';
  import SubmissionForm from './SubmissionForm.svelte';

  let currentForm = 'tip'; // 'tip', 'request', 'moment'

  // Reactive subscription to active video details
  let activeVideo = null;
  $: {
    activeVideo = $videoList.find(v => v.id === $activeVideoId) || null;
  }
</script>

<div class="grid grid-cols-1 lg:grid-cols-12 gap-6">
  <!-- Active Video Workspace (Left Pane - span 7) -->
  <div class="lg:col-span-7 flex flex-col justify-start">
    <div class="w-full aspect-video rounded-lg border-4 border-amber-500 overflow-hidden bg-black shadow-lg">
      {#if $activeVideoId}
        <iframe
          class="w-full h-full"
          src="https://www.youtube.com/embed/{$activeVideoId}"
          title="The Illumiteers - Active Video player"
          frameborder="0"
          allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share"
          allowfullscreen
        ></iframe>
      {:else}
        <div class="w-full h-full flex flex-col items-center justify-center text-white bg-neutral-900 p-6 relative overflow-hidden">
          <div class="absolute inset-0 bg-gradient-to-tr from-purple-950/30 to-teal-950/30 z-0"></div>
          <img 
            src="/assets/logo.png" 
            alt="The Illumiteers Logo" 
            class="w-16 h-16 rounded-full border-2 border-amber-400 shadow-[0_0_15px_rgba(245,158,11,0.3)] mb-4 animate-pulse z-10 object-cover" 
          />
          <p class="font-sans text-xs sm:text-sm font-semibold tracking-wide text-gray-300 z-10 uppercase tracking-wider text-center">
            Summoning active stream glimmer...
          </p>
        </div>
      {/if}
    </div>
    
    {#if activeVideo}
      <div class="mt-4 p-1">
        <h2 class="text-xl md:text-2xl font-serif font-black text-purple-950 tracking-tight leading-tight">
          {activeVideo.title}
        </h2>
        <div class="mt-2 flex items-center space-x-3 text-[#4A3565] opacity-80 font-normal tracking-wide text-sm font-sans">
          <span>Duration: {activeVideo.duration}</span>
          <span>•</span>
          <span>Published: {new Date(activeVideo.publishedAt).toLocaleDateString()}</span>
        </div>
      </div>
    {/if}
  </div>

  <!-- Submission Form Suite (Right Pane - span 5) -->
  <div class="lg:col-span-5 flex flex-col">
    <div class="bg-transparent border-4 border-black shadow-[6px_6px_0px_0px_rgba(0,0,0,1)] rounded-2xl overflow-hidden flex flex-col h-full min-h-[440px]">
      <!-- Tab Headers -->
      <div class="grid grid-cols-3 border-b-4 border-black bg-transparent">
        <button
          type="button"
          on:click={() => currentForm = 'tip'}
          class="py-3.5 px-2 text-center text-[10px] sm:text-xs font-bold font-sans uppercase tracking-wider transition-all duration-200 border-r-4 border-black
            {currentForm === 'tip' 
              ? 'bg-purple-900 text-white font-extrabold' 
              : 'text-purple-950 hover:bg-purple-900/10'}"
        >
          Submit a Tip
        </button>
        <button
          type="button"
          on:click={() => currentForm = 'request'}
          class="py-3.5 px-2 text-center text-[10px] sm:text-xs font-bold font-sans uppercase tracking-wider transition-all duration-200 border-r-4 border-black
            {currentForm === 'request' 
              ? 'bg-purple-900 text-white font-extrabold' 
              : 'text-purple-950 hover:bg-purple-900/10'}"
        >
          Content Request
        </button>
        <button
          type="button"
          on:click={() => currentForm = 'moment'}
          class="py-3.5 px-2 text-center text-[10px] sm:text-xs font-bold font-sans uppercase tracking-wider transition-all duration-200
            {currentForm === 'moment' 
              ? 'bg-purple-900 text-white font-extrabold' 
              : 'text-purple-950 hover:bg-purple-900/10'}"
        >
          Share a Moment
        </button>
      </div>
 
      <!-- Tab Content Area -->
      <div class="flex-grow flex flex-col justify-center bg-transparent p-5">
        <SubmissionForm {currentForm} />
      </div>
    </div>
  </div>
</div>
