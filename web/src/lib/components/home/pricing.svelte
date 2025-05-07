<script lang="ts">
	import * as Card from '$lib/components/ui/card';
	import Icon from '@iconify/svelte';
	import Badge from '../ui/badge/badge.svelte';
	import Button from '../ui/button/button.svelte';
	const pricingList = [
		{
			title: 'Free',
			popular: 0,
			price: 0,
			description: 'Lorem ipsum dolor sit, amet ipsum consectetur adipisicing elit.',
			buttonText: 'Get Started',
			benefitList: [
				'1 Team member',
				'2 GB Storage',
				'Upto 4 pages',
				'Community support',
				'lorem ipsum dolor'
			]
		},
		{
			title: 'Premium',
			popular: 1,
			price: 5,
			description: 'Lorem ipsum dolor sit, amet ipsum consectetur adipisicing elit.',
			buttonText: 'Start Free Trial',
			benefitList: [
				'4 Team member',
				'4 GB Storage',
				'Upto 6 pages',
				'Priority support',
				'lorem ipsum dolor'
			]
		},
		{
			title: 'Enterprise',
			popular: 0,
			price: 40,
			description: 'Lorem ipsum dolor sit, amet ipsum consectetur adipisicing elit.',
			buttonText: 'Contact US',
			benefitList: [
				'10 Team member',
				'8 GB Storage',
				'Upto 10 pages',
				'Priority support',
				'lorem ipsum dolor'
			]
		}
	];
</script>

<section id="pricing" class="container py-24 sm:py-32">
	<h2 class="text-center text-3xl font-bold md:text-4xl">
		Get
		<span class="bg-gradient-to-b from-primary/60 to-primary bg-clip-text text-transparent">
			{' '}
			Unlimited{' '}
		</span>
		Access
	</h2>
	<h3 class="pb-8 pt-4 text-center text-xl text-muted-foreground">
		Lorem ipsum dolor sit amet consectetur adipisicing elit. Alias reiciendis.
	</h3>
	<div class="grid gap-8 md:grid-cols-2 lg:grid-cols-3">
		{#each pricingList as pricing}
			<Card.Root
				class={pricing.popular === 1 ? 'shadow-black/10 drop-shadow-xl dark:shadow-white/10' : ''}
			>
				<Card.Header>
					<Card.Title class="item-center flex justify-between">
						{pricing.title}
						{#if pricing.popular === 1}
							<Badge variant="secondary" class="text-sm text-primary">Most popular</Badge>
						{/if}
					</Card.Title>
					<div>
						<span class="text-3xl font-bold">${pricing.price}</span>
						<span class="text-muted-foreground"> /month</span>
					</div>

					<Card.Description>{pricing.description}</Card.Description>
				</Card.Header>

				<Card.Content>
					<Button class="w-full">{pricing.buttonText}</Button>
				</Card.Content>

				<hr class="m-auto mb-4 w-4/5" />

				<Card.Footer class="flex">
					<div class="space-y-4">
						{#each pricing.benefitList as benefit}
							<span class="flex">
								<Icon icon="lucide:check" class="text-green-500" />
								<h3 class="ml-2">{benefit}</h3>
							</span>
						{/each}
					</div>
				</Card.Footer>
			</Card.Root>
		{/each}
	</div>
</section>
