<script lang="ts">
	import * as Card from '$lib/components/ui/card';

	const features = [
		{
			icon: '',
			title: 'Accessibility',
			description:
				'Lorem ipsum dolor sit amet, consectetur adipisicing elit. Illum quas provident cum'
		},
		{
			icon: '',
			title: 'Community',
			description:
				'Lorem ipsum dolor sit amet, consectetur adipisicing elit. Illum quas provident cum'
		},
		{
			icon: '',
			title: 'Scalability',
			description:
				'Lorem ipsum dolor sit amet, consectetur adipisicing elit. Illum quas provident cum'
		},
		{
			icon: '',
			title: 'Gamification',
			description:
				'Lorem ipsum dolor sit amet, consectetur adipisicing elit. Illum quas provident cum'
		}
	];
</script>

<section id="howItWorks" class="container py-24 text-center sm:py-32">
	<h2 class="text-3xl font-bold md:text-4xl">
		How It{' '}
		<span class="bg-gradient-to-b from-primary/60 to-primary bg-clip-text text-transparent">
			Works{' '}
		</span>
		Step-by-Step Guide
	</h2>
	<p class="mx-auto mb-8 mt-4 text-xl text-muted-foreground md:w-3/4">
		Lorem ipsum dolor sit amet consectetur, adipisicing elit. Veritatis dolor pariatur sit!
	</p>

	<div class="grid grid-cols-1 gap-8 md:grid-cols-2 lg:grid-cols-4">
		{#each features as { icon, title, description }}
			<Card.Root class="bg-muted/50">
				<Card.Header>
					<Card.Title class="grid place-items-center gap-4">
						{icon}
						{title}
					</Card.Title>
				</Card.Header>
				<Card.Content>{description}</Card.Content>
			</Card.Root>
		{/each}
	</div>
</section>
