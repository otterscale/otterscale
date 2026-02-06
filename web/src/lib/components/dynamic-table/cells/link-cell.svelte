<script lang="ts">
	let { display, hyperlink }: { display: string; hyperlink: string } = $props();
</script>

<!-- eslint-disable-next-line svelte/no-navigation-without-resolve -->
<a href={hyperlink} class="hover:underline">
	<p class="max-w-3xs truncate">{display}</p>
</a>
