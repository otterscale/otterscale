import { driver } from 'driver.js';

import { m } from '$lib/paraglide/messages';

export function startTour() {
	const driverObj = driver({
		showProgress: true,
		animate: true,
		allowClose: true,
		nextBtnText: m.next(),
		prevBtnText: m.back(),
		doneBtnText: m.done(),
		steps: [
			{
				element: '#cluster-guide-step',
				popover: {
					title: m.console_guide_cluster(),
					description: m.console_guide_cluster_description(),
					side: 'left',
					align: 'start'
				}
			},
			{
				element: '#workspace-guide-step',
				popover: {
					title: m.console_guide_workspace(),
					description: m.console_guide_workspace_description(),
					side: 'right',
					align: 'start'
				}
			},
			{
				element: '#sidebar-guide-step',
				popover: {
					title: m.console_guide_sidebar(),
					description: m.console_guide_sidebar_description(),
					side: 'right',
					align: 'start'
				}
			}
		]
	});

	driverObj.drive();
}
