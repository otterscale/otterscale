class RequestManager<Type> {
	INITIAL_VALUE: Type;
	request: Type = $state({} as Type);

	constructor(initialValue: Type) {
		this.INITIAL_VALUE = initialValue;
		this.request = this.INITIAL_VALUE;
	}

	reset(resetter?: () => void) {
		this.request = this.INITIAL_VALUE;
		if (resetter) {
			resetter();
		}
	}
}

export { RequestManager };
