<script lang="ts">
	import { resolve } from '$app/paths';

	let { display, hyperlink }: { display: string; hyperlink: string } = $props();
</script>

<a href={resolve(hyperlink)} class="hover:underline">
	<p class="max-w-3xs truncate">{display}</p>
</a>
