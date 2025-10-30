import type { AlertType, ValueType } from './types';

class IterationManager {
	alerts: AlertType[];
	value: ValueType;

	private duration: number;

	private interval: ReturnType<typeof setInterval>;

	constructor(alerts: AlertType[], duration: number, value: ValueType) {
		this.alerts = alerts;
		this.duration = duration;

		this.value = value;

		this.interval = this.getInterval();
	}

	getInterval() {
		return setInterval(() => {
			this.value.index = (this.value.index + 1) % this.alerts.length;
		}, this.duration);
	}

	start() {
		this.interval = this.getInterval();
	}

	stop() {
		clearInterval(this.interval);
	}

	previous() {
		this.value.index = (this.value.index - 1 + this.alerts.length) % this.alerts.length;
	}

	next() {
		this.value.index = (this.value.index + 1) % this.alerts.length;
	}
}

export { IterationManager };
