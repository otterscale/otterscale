<script lang="ts">
	import { resolve } from '$app/paths';

	let { display, hyperlink }: { display: string; hyperlink: string } = $props();
</script>

<a href={resolve(hyperlink)} class="hover:underline"> {display}</a>
