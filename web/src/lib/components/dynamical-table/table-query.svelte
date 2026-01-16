<script lang="ts">
	import { getCoreRowModel, getFilteredRowModel } from '@tanstack/table-core';
	import jsep from 'jsep';

	import { createSvelteTable, FlexRender } from '$lib/components/ui/data-table/index.js';
	import * as Table from '$lib/components/ui/table/index.js';

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

	function evaluate(abstractStructureTreeNode: any, row: any): any {
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

	const data = [
		{ column1: '1', column2: 1, column3: '2000-01-01', column4: true },
		{ column1: '2', column2: 2, column3: '2001-05-10', column4: false },
		{ column1: 'abc', column2: 3, column3: '1999-12-31', column4: true },
		{ column1: 'xyxy', column2: 4, column3: '2024-06-01', column4: false },
		{ column1: `a'b"a`, column2: 2.3, column3: '1979-08-16', column4: false }
	];

	const columns = [
		{ accessorKey: 'column1', header: 'Column 1' },
		{ accessorKey: 'column2', header: 'Column 2' },
		{ accessorKey: 'column3', header: 'Column 3' },
		{ accessorKey: 'column4', header: 'Column 4' }
	];

	let expression = $state('column1 = "1"');
	let globalFilter = $state('');

	function globalFilterFunction(row: any, columnId: string, filterValue: string) {
		if (!filterValue) return true;
		try {
			const ast = jsep(filterValue);
			console.log('AST:', JSON.stringify(ast, null, 2));
			return evaluate(ast, row.original);
		} catch (error) {
			console.error('Parse error:', error);
			return true;
		}
	}

	const table = createSvelteTable({
		get data() {
			return data;
		},
		columns,
		state: {
			get globalFilter() {
				return globalFilter;
			}
		},
		getCoreRowModel: getCoreRowModel(),
		getFilteredRowModel: getFilteredRowModel(),
		globalFilterFn: globalFilterFunction,
		onGlobalFilterChange: (updater) => {
			if (typeof updater === 'function') {
				globalFilter = updater(globalFilter);
			} else {
				globalFilter = updater;
			}
		}
	});
</script>

<div class="flex gap-2">
	<input
		bind:value={expression}
		class="w-full"
		onkeydown={(e) => {
			if (e.key === 'Enter') {
				table.setGlobalFilter(expression);
			}
		}}
	/>
</div>

<div class="rounded-md border">
	<Table.Root>
		<Table.Header>
			{#each table.getHeaderGroups() as headerGroup (headerGroup.id)}
				<Table.Row>
					{#each headerGroup.headers as header (header.id)}
						<Table.Head colspan={header.colSpan}>
							{#if !header.isPlaceholder}
								<FlexRender
									content={header.column.columnDef.header}
									context={header.getContext()}
								/>
							{/if}
						</Table.Head>
					{/each}
				</Table.Row>
			{/each}
		</Table.Header>
		<Table.Body>
			{#each table.getRowModel().rows as row (row.id)}
				<Table.Row data-state={row.getIsSelected() && 'selected'}>
					{#each row.getVisibleCells() as cell (cell.id)}
						<Table.Cell>
							<FlexRender content={cell.column.columnDef.cell} context={cell.getContext()} />
						</Table.Cell>
					{/each}
				</Table.Row>
			{:else}
				<Table.Row>
					<Table.Cell colspan={columns.length} class="h-24 text-center">No results.</Table.Cell>
				</Table.Row>
			{/each}
		</Table.Body>
	</Table.Root>
</div>
