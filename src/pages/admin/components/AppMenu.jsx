import React from 'react';
import { Sidenav, Nav, Icon } from 'rsuite';

export default function AppMenu(props) {
	return (
		<Sidenav {...props}>
			<Sidenav.Body>
				<Nav>
					<Nav.Item eventKey="1" icon={<Icon icon="dashboard" />}>
						<p className="app-menu-text">系统状态</p>
					</Nav.Item>
					<Nav.Item eventKey="2" icon={<Icon icon="peoples" />}>
						<p className="app-menu-text">用户管理</p>
					</Nav.Item>
					<Nav.Item eventKey="5" icon={<Icon icon="gear2" />}>
						<p className="app-menu-text">系统设置</p>
					</Nav.Item>
				</Nav>
			</Sidenav.Body>
		</Sidenav>
	);
}
