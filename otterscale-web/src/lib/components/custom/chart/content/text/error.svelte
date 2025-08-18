<script lang="ts">
	let {
		value,
		unit
	}: {
		value?: string;
		unit?: string;
	} = $props();
</script>

<div class="flex aspect-auto h-[250px] w-full items-center justify-center">
	<p class="fill-foreground text-4xl! font-bold">Unknown error</p>
</div>
