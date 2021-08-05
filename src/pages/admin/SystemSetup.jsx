import React from 'react';
import {
	FlexboxGrid,
	Row,
	Col,
	Input,
	Divider,
	Button,
	Avatar,
	Tag,
	InputNumber,
	InputGroup,
	DatePicker,
	SelectPicker,
	Notification,
	Nav,
	Dropdown,
	Icon,
} from 'rsuite';

import PluginSetup from './SystemSetups/PluginSetup';

import useAxios from 'axios-hooks';
import axios from 'axios';

const { useState, useEffect } = React;

const CustomNav = ({ active, onSelect, ...props }) => {
	return (
		<Nav {...props} vertical activeKey={active} onSelect={onSelect} style={{ width: 100, height: '100%' }}>
			<Nav.Item eventKey="plugin">插件设定</Nav.Item>
			<Nav.Item eventKey="api">接口设定</Nav.Item>
			<Nav.Item eventKey="backstage">后台设定</Nav.Item>
		</Nav>
	);
};

export default function SystemSetup() {
	const [navActive, setnavActive] = useState('plugin');

	return (
		<Row style={{ paddingTop: 25, paddingBottom: 25 }}>
			<Col md={4}>
				<CustomNav appearance="subtle" active={navActive} onSelect={setnavActive} />
			</Col>
			<Col style={{ flex: 1 }}>
				<div>{navActive === 'plugin' && <PluginSetup />}</div>
			</Col>
		</Row>
	);
}
