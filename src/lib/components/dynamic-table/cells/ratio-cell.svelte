<script lang="ts">
	let {
		numerator,
		denominator
	}: {
		numerator: any;
		denominator: any;
	} = $props();
</script>

{#if denominator !== undefined}
	{numerator ?? '?'}/{denominator}
{/if}
