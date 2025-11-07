<script lang="ts" module>
	import { onMount } from 'svelte';

	import { items } from './data';

	import { goto } from '$app/navigation';
</script>

<script lang="ts">
	onMount(() => {
		const defaultItem = items.find((item) => item.default);
		const [firstItem] = items;
		// eslint-disable-next-line svelte/no-navigation-without-resolve
		goto(defaultItem ? defaultItem.url : firstItem.url);
	});
</script>
