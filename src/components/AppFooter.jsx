import React from 'react';

import { Footer } from 'rsuite';

export default function AppFooter() {
	return (
		<Footer
			style={{
				textAlign: 'center',
				padding: 15,
				background: '#00000006',
				borderTop: '1px solid rgb(229, 229, 234)',
			}}
		>
			Copyright © 2020 - {new Date().getFullYear()} <b>HaXiBiao Developer.</b>
		</Footer>
	);
}
