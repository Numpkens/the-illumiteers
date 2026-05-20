<script>
  export let type = 'video';
</script>

{#if type === 'video'}
  <div class="animate-pulse bg-[#F1E8D9]/40 border border-black/5 aspect-[2.5/3.5] w-full max-w-[360px] rounded-2xl flex flex-col justify-between p-4 shadow-sm relative">
    <div class="absolute inset-1 border border-black/5 rounded-[12px] pointer-events-none"></div>
    <div class="w-full aspect-[4/3] bg-black/5 rounded-t-lg"></div>
    <div class="h-4 bg-black/5 rounded mt-3 mb-1 w-3/4"></div>
    <div class="h-3 bg-black/5 rounded mb-2 w-1/2"></div>
    <div class="flex-grow bg-black/5 rounded p-3 mb-1 min-h-[60px]"></div>
  </div>
{:else}
  <div class="animate-pulse bg-[#F1E8D9]/40 border border-black/5 w-full sm:w-[260px] md:w-[280px] h-[280px] rounded-xl flex flex-col justify-between p-3 shadow-sm relative">
    <div class="absolute inset-1 border border-black/5 rounded-[10px] pointer-events-none"></div>
    <div class="h-8 bg-black/5 rounded-t-md mb-2"></div>
    <div class="h-3 bg-black/5 rounded mb-2 w-1/2"></div>
    <div class="flex-grow bg-black/5 rounded-b-md p-3 min-h-[120px]"></div>
  </div>
{/if}
