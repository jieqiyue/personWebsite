<template>
  <div class="comments-section">
    <hr class="comments-separator">
    <h2 class="comments-title">评论</h2>
    <div ref="giscusContainer"></div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch, nextTick } from 'vue';
import { useRoute } from 'vue-router'; // Optional: if theme needs to react to system changes

const props = defineProps({
  pageTerm: {
    type: String,
    required: true,
  },
});

const giscusContainer = ref(null);
const route = useRoute(); // Optional: for theme reactivity

// Function to load the Giscus script
const loadGiscus = (term) => {
  if (!giscusContainer.value) return;

  // Clear previous Giscus instance if exists
  while (giscusContainer.value.firstChild) {
    giscusContainer.value.removeChild(giscusContainer.value.firstChild);
  }

  const script = document.createElement('script');
  script.src = 'https://giscus.app/client.js';
  script.async = true;
  script.crossOrigin = 'anonymous';

  // --- YOUR GISCUS CONFIGURATION BELOW ---
  script.setAttribute('data-repo', 'jieqiyue/personWebsite'); // <-- Replace with your value if needed
  script.setAttribute('data-repo-id', 'R_kgDOOV5N5Q');      // <-- Replace with your value if needed
  script.setAttribute('data-category', 'Announcements');     // <-- Replace with your value if needed
  script.setAttribute('data-category-id', 'DIC_kwDOOV5N5c4Co_3Z'); // <-- Replace with your value if needed
  script.setAttribute('data-mapping', 'pathname');          // <-- Replace with your value if needed
  script.setAttribute('data-term', term); // Use the dynamic term passed via prop
  script.setAttribute('data-strict', '0');               // <-- Replace with your value if needed
  script.setAttribute('data-reactions-enabled', '1');  // <-- Replace with your value if needed
  script.setAttribute('data-emit-metadata', '0');        // <-- Replace with your value if needed
  script.setAttribute('data-input-position', 'bottom');    // <-- Replace with your value if needed
  script.setAttribute('data-theme', 'preferred_color_scheme'); // <-- Replace with your value if needed
  script.setAttribute('data-lang', 'zh-CN');             // <-- Replace with your value if needed
  // --- END OF YOUR GISCUS CONFIGURATION ---

  giscusContainer.value.appendChild(script);
};

// Load Giscus when the component mounts
onMounted(() => {
  loadGiscus(props.pageTerm);
});

// Reload Giscus when the pageTerm prop changes (e.g., navigating between articles)
watch(() => props.pageTerm, (newTerm) => {
  // Use nextTick to ensure the container is ready after potential DOM updates
  nextTick(() => {
    loadGiscus(newTerm);
  });
});

// Optional: Watch for system theme changes if using preferred_color_scheme
// This might require more complex setup depending on how theme is managed globally
// watch(() => route.meta.theme, (newTheme) => { // Example: if theme is in route meta
//   if (giscusContainer.value && giscusContainer.value.querySelector('iframe.giscus-frame')) {
//      const iframe = giscusContainer.value.querySelector('iframe.giscus-frame');
//      iframe.contentWindow.postMessage({ giscus: { setConfig: { theme: newTheme } } }, 'https://giscus.app');
//   }
// });

</script>

<style scoped>
.comments-section {
  margin-top: 3rem; /* Adjust overall top margin for the section */
}

.comments-separator {
  border: none;
  border-top: 1px solid #eee; /* A simple light gray line */
  margin: 2.5rem 0;      /* Add vertical space around the separator */
}

.comments-title {
  font-size: 1.6em;      /* Adjust title size as needed */
  font-weight: 600;
  margin-bottom: 1.5rem; /* Space below the title */
}

/* Ensure the Giscus container itself doesn't add extra margin */
/* No specific style needed for giscusContainer div unless overriding defaults */
</style> 