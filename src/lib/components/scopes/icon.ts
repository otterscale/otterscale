const SCOPE_ICONS = [
	'ph:airplane-tilt',
	'ph:cactus',
	'ph:cherries',
	'ph:piggy-bank',
	'ph:flower',
	'ph:joystick',
	'ph:clover',
	'ph:cube',
	'ph:gavel'
];

export function scopeIcon(index: number): string {
	return index !== -1 ? SCOPE_ICONS[index % SCOPE_ICONS.length] : SCOPE_ICONS[0];
}
