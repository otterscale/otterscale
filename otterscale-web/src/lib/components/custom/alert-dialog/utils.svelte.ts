export class StateController {
	state = $state(false);

	constructor(initialState = false) {
		this.state = initialState;
	}

	close() {
		this.state = false;
	}
}
