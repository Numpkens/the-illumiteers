import { writable } from 'svelte/store';
import { API_BASE } from '../lib/api';

export const videoList = writable([]);
export const activeVideoId = writable('');
export const isLoading = writable(true);
export const error = writable(null);

export async function loadVideos() {
  isLoading.set(true);
  error.set(null);
  try {
    const res = await fetch(`${API_BASE}/youtube`);
    if (!res.ok) {
      throw new Error(`HTTP error! status: ${res.status}`);
    }
    const data = await res.json();
    videoList.set(data);
    if (data && data.length > 0) {
      activeVideoId.set(data[0].id);
    }
  } catch (err) {
    console.error('Error fetching videos from API:', err);
    error.set(err.message || 'Could not load video data. Check backend status.');
  } finally {
    isLoading.set(false);
  }
}
