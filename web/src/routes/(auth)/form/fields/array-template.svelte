<script lang="ts">
	import { type ComponentProps, getComponent, getFormContext } from '@sjsf/form';

	const { children, buttons, config, errors }: ComponentProps['arrayItemTemplate'] = $props();

	const ctx = getFormContext();

	const Layout = $derived(getComponent(ctx, 'layout', config));
</script>

<div class="relative">
	<Layout type="array-item" {config} {errors}>
		<Layout type="array-item-content" {config} {errors}>
			{@render children()}
		</Layout>
		<div class="absolute -top-1 right-0">
			{#if buttons}
				<Layout type="array-item-controls" {config} {errors}>
					{@render buttons()}
				</Layout>
			{/if}
		</div>
	</Layout>
</div>
