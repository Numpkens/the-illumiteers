import { writable } from 'svelte/store';

export const blogList = writable([]);
export const isBlogLoading = writable(true);
export const blogError = writable(null);

export async function loadBlogPosts() {
  isBlogLoading.set(true);
  blogError.set(null);
  try {
    const response = await fetch('http://localhost:8080/api/blog');
    if (!response.ok) {
      throw new Error(`Failed to load blog scrolls. Server status: ${response.status}`);
    }
    const data = await response.json();
    blogList.set(data || []);
  } catch (err) {
    console.error('Error in loadBlogPosts store:', err);
    blogError.set(err.message || 'Could not retrieve strategy guides.');
  } finally {
    isBlogLoading.set(false);
  }
}
