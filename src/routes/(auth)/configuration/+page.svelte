<script lang="ts" module>
	import { onMount } from 'svelte';

	import { goto } from '$app/navigation';

	import { items } from './data';
</script>

<script lang="ts">
	onMount(() => {
		const defaultItem = items.find((item) => item.default);
		const [firstItem] = items;
		// eslint-disable-next-line svelte/no-navigation-without-resolve
		goto(defaultItem ? defaultItem.url : firstItem.url);
	});
</script>
