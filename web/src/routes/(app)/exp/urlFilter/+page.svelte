<script>
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
</script>

1. {page.url.searchParams}
2. {page.url.hash}
