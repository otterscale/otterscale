<script lang="ts" module>
	export function evaluate(abstractStructureTreeNode: any, row: any): any {
		switch (abstractStructureTreeNode.type) {
			case 'UnaryExpression': {
				const argument =
					abstractStructureTreeNode.argument.type === 'Identifier'
						? row[abstractStructureTreeNode.argument.name]
						: evaluate(abstractStructureTreeNode.argument, row);
				switch (abstractStructureTreeNode.operator) {
					case '!':
						return !argument;
					default:
						return false;
				}
			}
			case 'CallExpression': {
				const callee = abstractStructureTreeNode.callee;
				const relation = callee.type === 'Identifier' ? callee.name : null;
				const nodeArguments = abstractStructureTreeNode.arguments.map((argument: any) =>
					evaluate(argument, row)
				);
				const [nodeArgument] = nodeArguments;
				switch (relation) {
					case 'floor':
						return Math.floor(Number(nodeArgument));
					case 'ceil':
						return Math.ceil(Number(nodeArgument));
					case 'round':
						return Math.round(Number(nodeArgument));
					default:
						return false;
				}
			}
			case 'BinaryExpression': {
				const left =
					abstractStructureTreeNode.left.type === 'Identifier'
						? row[abstractStructureTreeNode.left.name]
						: evaluate(abstractStructureTreeNode.left, row);
				const right =
					abstractStructureTreeNode.right.type === 'Literal'
						? abstractStructureTreeNode.right.value
						: evaluate(abstractStructureTreeNode.right, row);

				switch (abstractStructureTreeNode.operator) {
					case '|':
						return Boolean(left) || Boolean(right);
					case '&':
						return Boolean(left) && Boolean(right);

					case '>':
						if (!isNaN(new Date(left).getTime()) && !isNaN(new Date(right).getTime())) {
							return new Date(left).getTime() > new Date(right).getTime();
						}
						return Number(left) > Number(right);
					case '>=':
						if (!isNaN(new Date(left).getTime()) && !isNaN(new Date(right).getTime())) {
							return new Date(left).getTime() >= new Date(right).getTime();
						}
						return Number(left) >= Number(right);
					case '<':
						if (!isNaN(new Date(left).getTime()) && !isNaN(new Date(right).getTime())) {
							return new Date(left).getTime() < new Date(right).getTime();
						}
						return Number(left) < Number(right);
					case '<=':
						if (!isNaN(new Date(left).getTime()) && !isNaN(new Date(right).getTime())) {
							return new Date(left).getTime() <= new Date(right).getTime();
						}
						return Number(left) <= Number(right);

					case '=':
						return left == right;
					case '~':
						return String(left).includes(String(right));
					case '^~':
						return String(left).startsWith(String(right));
					case '~$':
						return String(left).endsWith(String(right));
					case '^~$':
						return String(left).startsWith(String(right)) && String(left).endsWith(String(right));

					case '@': {
						if (abstractStructureTreeNode.right.type === 'ArrayExpression') {
							const right = abstractStructureTreeNode.right.elements.map(
								(element: any) => element.value
							);
							return right.includes(left);
						}
						return false;
					}
					default:
						return false;
				}
			}
			case 'Identifier':
				return row[abstractStructureTreeNode.name];
			case 'Literal':
				return abstractStructureTreeNode.value;
			default:
				return false;
		}
	}
</script>

<script lang="ts">
	import SearchIcon from '@lucide/svelte/icons/search';
	import X from '@lucide/svelte/icons/x';
	import type { JsonValue } from '@openfeature/server-sdk';
	import type { Table } from '@tanstack/table-core';
	import jsep from 'jsep';

	import { Button } from '$lib/components/ui/button/index.js';
	import * as ButtonGroup from '$lib/components/ui/button-group/index.js';
	import * as InputGroup from '$lib/components/ui/input-group/index.js';

	let {
		expression = $bindable(),
		table
	}: { expression: string; table: Table<Record<string, JsonValue>> } = $props();

	jsep.addUnaryOp('!');

	jsep.addBinaryOp('|', 5);
	jsep.addBinaryOp('&', 5);

	jsep.addBinaryOp('>', 6);
	jsep.addBinaryOp('>=', 6);
	jsep.addBinaryOp('<', 6);
	jsep.addBinaryOp('<=', 6);

	jsep.addBinaryOp('=', 7);
	jsep.addBinaryOp('~', 7);
	jsep.addBinaryOp('^~', 7);
	jsep.addBinaryOp('~$', 7);
	jsep.addBinaryOp('^~$', 7);

	jsep.addBinaryOp('@', 8);

	jsep.hooks.add('gobble-token', function (this: any, environment: any) {
		const character = this.code;

		if (character === 96) {
			this.index = this.index + 1;

			let value = '';
			while (this.index < this.expr.length) {
				const code = this.expr.charCodeAt(this.index);

				if (code === 96) {
					this.index = this.index + 1;
					break;
				}

				value = value + this.expr[this.index];
				this.index = this.index + 1;
			}

			environment.node = {
				type: 'Literal',
				value: value,
				raw: '`' + value + '`'
			};
		}
	});

	function handleEnterDown(event: KeyboardEvent) {
		if (event.key === 'Enter') {
			table.setGlobalFilter(expression);
		}
	}
	function handleSearch() {
		table.setGlobalFilter(expression);
	}
	function handleClear() {
		expression = '';
	}
</script>

<ButtonGroup.Root class="w-full">
	<InputGroup.Root>
		<InputGroup.Input
			placeholder="Search via Query Language"
			bind:value={expression}
			class="w-full"
			onkeydown={handleEnterDown}
		/>
		<InputGroup.Addon>
			<SearchIcon />
		</InputGroup.Addon>
	</InputGroup.Root>
	<Button
		class={expression ? 'flex' : 'hidden'}
		variant="outline"
		size="icon"
		aria-label="Search"
		onclick={handleClear}
	>
		<X class="text-destructive" />
	</Button>
	<Button variant="outline" size="icon" aria-label="Search" onclick={handleSearch}>
		<SearchIcon />
	</Button>
</ButtonGroup.Root>
