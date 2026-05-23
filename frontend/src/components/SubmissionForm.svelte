<script>
  import { API_BASE } from '../lib/api';
  export let currentForm = 'tip';

  // State management
  let loading = false;
  let success = false;
  let errorMessage = '';

  // Form fields: Tip
  let tipName = '';
  let tipContact = '';
  let tipContent = '';
  let tipReference = '';

  // Form fields: Content Request
  let requestRequester = '';
  let requestTitle = '';
  let requestDetails = '';
  let requestReference = '';

  // Form fields: Share a Moment
  let momentSubmitter = '';
  let momentDetails = '';
  let momentMedia = '';

  function resetForm() {
    success = false;
    errorMessage = '';
    
    // Clear all fields
    tipName = '';
    tipContact = '';
    tipContent = '';
    tipReference = '';

    requestRequester = '';
    requestTitle = '';
    requestDetails = '';
    requestReference = '';

    momentSubmitter = '';
    momentDetails = '';
    momentMedia = '';
  }

  // Reactive listener to reset validation states when toggling tabs
  $: {
    if (currentForm) {
      errorMessage = '';
      success = false;
    }
  }

  async function handleSubmit(e) {
    e.preventDefault();
    loading = true;
    errorMessage = '';

    const payload = { type: currentForm };

    // Set payloads and local verification lengths
    if (currentForm === 'tip') {
      if (!tipName.trim() || !tipContent.trim()) {
        errorMessage = 'Name and Tip details are required.';
        loading = false;
        return;
      }
      payload.name = tipName;
      payload.contact = tipContact;
      payload.content = tipContent;
      payload.referenceLink = tipReference;
    } else if (currentForm === 'request') {
      if (!requestRequester.trim() || !requestTitle.trim() || !requestDetails.trim()) {
        errorMessage = 'Name, Suggested Title, and Details are required.';
        loading = false;
        return;
      }
      payload.requester = requestRequester;
      payload.suggestedTitle = requestTitle;
      payload.requestDetails = requestDetails;
      payload.referenceLink = requestReference;
    } else if (currentForm === 'moment') {
      if (!momentSubmitter.trim() || !momentDetails.trim()) {
        errorMessage = 'Your Name and Moment details are required.';
        loading = false;
        return;
      }
      payload.submitter = momentSubmitter;
      payload.eventDescription = momentDetails;
      payload.mediaUrl = momentMedia;
    }

    try {
      const response = await fetch(`${API_BASE}/submit`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(payload),
      });

      if (!response.ok) {
        const errText = await response.text();
        throw new Error(errText || 'Server error occurred during submission.');
      }

      success = true;
    } catch (err) {
      console.error(err);
      errorMessage = err.message || 'Connection to server failed. Please try again.';
    } finally {
      loading = false;
    }
  }
</script>

<div class="w-full h-full flex flex-col justify-center">
  {#if success}
    <div class="flex flex-col items-center justify-center text-center py-10 px-4 transition-all duration-300">
      <!-- TCG Card Game themed success illustration -->
      <div class="relative w-28 h-40 bg-amber-100 border-4 border-black rounded-2xl flex flex-col items-center justify-center rotate-6 shadow-[5px_5px_0px_0px_rgba(0,0,0,1)] mb-8">
        <!-- Inside Card Cost Badge -->
        <div class="absolute -top-3 -right-3 w-8 h-8 rotate-45 border-2 border-black bg-white flex items-center justify-center shadow-sm">
          <span class="-rotate-45 block font-serif font-black text-sm text-purple-950">✓</span>
        </div>
        <!-- Card Art Frame Mockup -->
        <div class="w-[85%] h-20 bg-purple-900 border-2 border-black rounded-lg flex items-center justify-center mb-2 overflow-hidden">
          <span class="text-white text-3xl">✨</span>
        </div>
        <!-- Monospace text footer -->
        <div class="text-[8px] font-mono tracking-widest text-gray-800 uppercase mt-2">
          SUCCESS
        </div>
      </div>

      <div class="bg-[#F1E8D9] border-2 border-[#3A2E1A]/20 shadow-[inset_0_1px_3px_rgba(0,0,0,0.05)] rounded-xl p-6 max-w-xs sm:max-w-sm mx-auto text-center mb-8">
        <p class="font-serif italic text-base text-[#3A2E1A] leading-relaxed">
          "Your scroll has reached the Great Illuminary. Thank you, Illumineer!"
        </p>
      </div>
      
      <button 
        type="button"
        on:click={resetForm}
        class="px-6 py-3 bg-amber-100 hover:bg-amber-200 text-purple-950 text-xs font-semibold tracking-wider uppercase border-2 border-black rounded-xl shadow-[3px_3px_0px_0px_rgba(0,0,0,1)] hover:translate-x-[1px] hover:translate-y-[1px] hover:shadow-[2px_2px_0px_0px_rgba(0,0,0,1)] transition-all">
        Make Another Submission
      </button>
    </div>
  {:else}
    <form on:submit={handleSubmit} class="space-y-4">
      {#if errorMessage}
        <div class="p-3 bg-red-50 border-2 border-red-500 text-red-700 text-xs rounded-lg font-sans">
          <strong>Error:</strong> {errorMessage}
        </div>
      {/if}

      {#if currentForm === 'tip'}
        <div>
          <label for="tip-name" class="block text-xs font-bold uppercase text-purple-950 mb-1">Your Name / Handle *</label>
          <input 
            id="tip-name"
            type="text" 
            bind:value={tipName} 
            placeholder="e.g. IllumiteerPro" 
            required
            class="w-full text-sm font-sans px-3 py-2 border-2 border-[#4A3565] rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-900 bg-white"
          />
        </div>

        <div>
          <label for="tip-contact" class="block text-xs font-bold uppercase text-purple-950 mb-1">Contact Info (Optional)</label>
          <input 
            id="tip-contact"
            type="text" 
            bind:value={tipContact} 
            placeholder="Discord ID or Email (for credit/questions)" 
            class="w-full text-sm font-sans px-3 py-2 border-2 border-[#4A3565] rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-900 bg-white"
          />
        </div>

        <div>
          <label for="tip-content" class="block text-xs font-bold uppercase text-purple-950 mb-1">Tip Content *</label>
          <textarea 
            id="tip-content"
            bind:value={tipContent} 
            placeholder="Provide specific details about leak, strategy, or news..." 
            rows="4"
            required
            class="w-full text-sm font-sans px-3 py-2 border-2 border-[#4A3565] rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-900 bg-white resize-none"
          ></textarea>
        </div>

        <div>
          <label for="tip-reference" class="block text-xs font-bold uppercase text-purple-950 mb-1">Reference Link (Optional)</label>
          <input 
            id="tip-reference"
            type="url" 
            bind:value={tipReference} 
            placeholder="https://twitter.com/..." 
            class="w-full text-sm font-sans px-3 py-2 border-2 border-[#4A3565] rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-900 bg-white"
          />
        </div>

      {:else}
        <!-- request -->
        {#if currentForm === 'request'}
          <div>
            <label for="req-name" class="block text-xs font-bold uppercase text-purple-950 mb-1">Your Name *</label>
            <input 
              id="req-name"
              type="text" 
              bind:value={requestRequester} 
              placeholder="Your handle or Nickname" 
              required
              class="w-full text-sm font-sans px-3 py-2 border-2 border-[#4A3565] rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-900 bg-white"
            />
          </div>

          <div>
            <label for="req-title" class="block text-xs font-bold uppercase text-purple-950 mb-1">Suggested Video Title *</label>
            <input 
              id="req-title"
              type="text" 
              bind:value={requestTitle} 
              placeholder="e.g. Ruby/Steel Deck Analysis" 
              required
              class="w-full text-sm font-sans px-3 py-2 border-2 border-[#4A3565] rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-900 bg-white"
            />
          </div>

          <div>
            <label for="req-details" class="block text-xs font-bold uppercase text-purple-950 mb-1">Request Details / Topic Description *</label>
            <textarea 
              id="req-details"
              bind:value={requestDetails} 
              placeholder="What specifically do you want Liam to cover or analyze?" 
              rows="4"
              required
              class="w-full text-sm font-sans px-3 py-2 border-2 border-[#4A3565] rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-900 bg-white resize-none"
            ></textarea>
          </div>

          <div>
            <label for="req-reference" class="block text-xs font-bold uppercase text-purple-950 mb-1">Reference Link (Optional)</label>
            <input 
              id="req-reference"
              type="url" 
              bind:value={requestReference} 
              placeholder="https://dreamborn.ink/decks/..." 
              class="w-full text-sm font-sans px-3 py-2 border-2 border-[#4A3565] rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-900 bg-white"
            />
          </div>

        <!-- moment -->
        {:else}
          <div>
            <label for="moment-name" class="block text-xs font-bold uppercase text-purple-950 mb-1">Submitter Name *</label>
            <input 
              id="moment-name"
              type="text" 
              bind:value={momentSubmitter} 
              placeholder="Your Name" 
              required
              class="w-full text-sm font-sans px-3 py-2 border-2 border-[#4A3565] rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-900 bg-white"
            />
          </div>

          <div>
            <label for="moment-details" class="block text-xs font-bold uppercase text-purple-950 mb-1">Milestone Details & Event Description *</label>
            <textarea 
              id="moment-details"
              bind:value={momentDetails} 
              placeholder="Describe your moment! e.g. 'Won my first local store championship with Liam's Ruby Control list!'" 
              rows="5"
              required
              class="w-full text-sm font-sans px-3 py-2 border-2 border-[#4A3565] rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-900 bg-white resize-none"
            ></textarea>
          </div>

          <div>
            <label for="moment-media" class="block text-xs font-bold uppercase text-purple-950 mb-1">Media URL (Optional)</label>
            <input 
              id="moment-media"
              type="url" 
              bind:value={momentMedia} 
              placeholder="Link to image, YouTube clip, or tweet" 
              class="w-full text-sm font-sans px-3 py-2 border-2 border-[#4A3565] rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-900 bg-white"
            />
          </div>
        {/if}
      {/if}

      <div class="pt-2">
        <button 
          type="submit" 
          disabled={loading}
          class="w-full py-3 bg-purple-900 text-white font-sans text-xs uppercase font-bold tracking-widest border-2 border-black rounded-xl shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] hover:translate-x-[1px] hover:translate-y-[1px] hover:shadow-[3px_3px_0px_0px_rgba(0,0,0,1)] disabled:bg-gray-400 disabled:shadow-none disabled:translate-y-0 disabled:cursor-not-allowed transition-all">
          {#if loading}
            <span class="inline-block animate-spin mr-2">✦</span> Processing...
          {:else}
            Cast Submission ✦
          {/if}
        </button>
      </div>
    </form>
  {/if}
</div>
