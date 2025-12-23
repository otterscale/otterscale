import type { UnitType } from './types';

const INPUT_CLASSNAME =
	'border-input placeholder:text-muted-foreground focus-visible:ring-ring flex h-9 w-full rounded-md border bg-transparent text-base shadow-sm transition-colors file:border-0 file:bg-transparent file:text-sm file:font-medium focus-visible:outline-none focus-visible:ring-1 disabled:cursor-not-allowed disabled:opacity-50 md:text-sm';
const typeToIcon: Record<string, string> = {
	color: 'ph:palette',
	'datetime-local': 'ph:clock',
	date: 'ph:calendar',
	time: 'ph:clock',
	url: 'ph:link',
	email: 'ph:mailbox',
	tel: 'ph:phone',
	switch: 'ph:toggle-left',
	checkbox: 'ph:check-square',
	text: 'ph:textbox',
	number: 'ph:list-numbers',
	search: 'ph:magnifying-glass',
	password: 'ph:password'
};

function getInputMeasurementUnitByValue(
	value: number | undefined,
	units: UnitType[]
): { value: number | undefined; unit: UnitType | undefined } {
	const sortedUnits = units.sort((p, n) => p.value - n.value);

	if (value === undefined) {
		return { value: undefined, unit: sortedUnits[0] };
	}

	const rawValue = Number(value);

	let temporaryUnit = sortedUnits[0];
	let temporaryValue = rawValue / sortedUnits[0].value;
	for (const unit of sortedUnits) {
		if (rawValue / unit.value >= 1) {
			temporaryValue = rawValue / unit.value;
			temporaryUnit = unit;
		}
	}
	return { value: temporaryValue, unit: temporaryUnit };
}

class PasswordManager {
	isVisible = $state<boolean>(false);

	enable() {
		this.isVisible = true;
	}

	disable() {
		this.isVisible = false;
	}
}

export { getInputMeasurementUnitByValue, INPUT_CLASSNAME, PasswordManager, typeToIcon };
