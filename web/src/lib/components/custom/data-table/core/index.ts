import ActionItem from './data-table-action-item.svelte';
import ActionList from './data-table-actions.svelte';
import CellRowExpander from './data-table-cell-row-expander.svelte';
import CellRowPicker from './data-table-cell-row-picker.svelte';
import Empty from './data-table-empty.svelte';
import FilterBooleanMatch from './data-table-filter-boolean-match.svelte';
import FilterColumn from './data-table-filter-column.svelte';
import FilterStringFuzzy from './data-table-filter-string-fuzzy.svelte';
import FilterStringMatch from './data-table-filter-string-match.svelte';
import Footer from './data-table-footer.svelte';
import HeaderRowExpander from './data-table-header-row-expander.svelte';
import HeaderRowPicker from './data-table-header-row-picker.svelte';
import Pagination from './data-table-pagination.svelte';
import Sorter from './data-table-sorter.svelte';
import { getSortingFunction } from './utils.svelte';

const Filters = {
	BooleanMatch: FilterBooleanMatch,
	Column: FilterColumn,
	StringMatch: FilterStringMatch,
	StringFuzzy: FilterStringFuzzy
};

const Cells = {
	RowPicker: CellRowPicker,
	RowExpander: CellRowExpander
};

const Headers = {
	RowPicker: HeaderRowPicker,
	RowExpander: HeaderRowExpander
};

const Actions = {
	ActionList: ActionList,
	ActionItem: ActionItem
};

export { Actions, Cells, Empty, Filters, Footer, getSortingFunction, Headers, Pagination, Sorter };
