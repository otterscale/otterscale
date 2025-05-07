<script lang="ts">
	import Icon from '@iconify/svelte';

	const sponsors = [
		{
			icon: '',
			name: 'Sponsor 1'
		},
		{
			icon: '',
			name: 'Sponsor 2'
		},
		{
			icon: '',
			name: 'Sponsor 3'
		},
		{
			icon: '',
			name: 'Sponsor 4'
		},
		{
			icon: '',
			name: 'Sponsor 5'
		},
		{
			icon: '',
			name: 'Sponsor 6'
		}
	];
</script>

<section id="sponsors" class="container pt-24 sm:py-32">
	<h2 class="text-md mb-8 text-center font-bold text-primary lg:text-xl">Investors and founders</h2>

	<div class="flex flex-wrap items-center justify-center gap-4 md:gap-8">
		{#each sponsors as { icon, name }}
			<div class="flex items-center gap-1 text-muted-foreground/60">
				<span><Icon {icon} /></span>
				<h3 class="text-xl font-bold">{name}</h3>
			</div>
		{/each}
	</div>
</section>
