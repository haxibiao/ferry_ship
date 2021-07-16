import React, { useState, useEffect } from 'react';
import {
	Header,
	Navbar,
	Avatar,
	Nav,
	Dropdown,
	Whisper,
	Popover,
	Modal,
	Notification,
	Form,
	FormControl,
	FormGroup,
	ControlLabel,
	Button,
} from 'rsuite';

import useAxios from 'axios-hooks';
import axios from 'axios';

export default function AppHeader(props) {
	const { user } = props;
	const triggerRef = React.createRef();

	function handleSelectMenu(eventKey, event) {
		switch (eventKey) {
			case 1:
				setshowUpdateUser(true);
				break;
			case 2:
				__unLogin();
				break;
		}
		triggerRef.current.hide();
	}

	// 注销登陆
	const __unLogin = () => {
		// localStorage.removeItem("user");
		document.cookie = 'u_token=; expires=Thu, 01 Jan 1970 00:00:00 GMT';
		window.location.href = '/login';
	};

	// 修改用户密码操作
	const [showUpdateUser, setshowUpdateUser] = useState(false);
	const [updateUserLoading, setupdateUserLoading] = useState(false);
	const [updateUser, setupdateUser] = useState({
		id: user?.id,
		name: user?.name,
		passwd: '',
	});
	const onUpdateUser = (value) => {
		setupdateUser({ ...updateUser, ...value });
	};
	const APIUpdateUser = (id, passwd) => {
		if (!id || !passwd) {
			Notification.error({
				title: '新密码不得为空！',
			});
			return;
		}

		setupdateUserLoading(true);

		const params = new URLSearchParams();
		params.append('id', id);
		params.append('password', passwd);

		axios
			.post('/api/user/update', params, {})
			.then((res) => {
				setupdateUserLoading(false);

				const { data } = res;
				const { code } = data;

				if (code < 1) {
					Notification.error({
						title: data?.msg || '修改失败，请稍后重试！',
					});
				} else {
					const user = data?.data;

					Notification.success({
						title: `修改密码成功！`,
						onClose: () => {
							__unLogin(); // 退出登陆
						},
					});

					// 修改用户成功，关闭弹窗，刷新列表数据，清空编辑框数据
					setshowUpdateUser(false);
				}
			})
			.catch((error) => {
				Notification.error({
					title: '修改失败，' + error || '修改失败，请稍后重试！',
				});
				setupdateUserLoading(false);
			});
	};
	useEffect(() => {
		setupdateUser({ ...updateUser, id: user?.id, name: user?.name });
	}, [user]);

	return (
		<>
			<Header>
				<Navbar appearance="inverse" style={{ padding: '0 15px' }}>
					<Navbar.Header justify="center">
						<a className="navbar-brand logo" style={{ textDecoration: 'none' }}>
							<h3 style={{ lineHeight: '56px', color: '#FFFB' }}>Ferry ship 控制台</h3>
						</a>
					</Navbar.Header>
					{user && (
						<Nav pullRight justify="center" style={{ lineHeight: '56px' }}>
							<Whisper
								placement="bottomEnd"
								trigger="click"
								triggerRef={triggerRef}
								speaker={
									<Popover full>
										<Dropdown.Menu onSelect={handleSelectMenu}>
											<Dropdown.Item eventKey={1}>修改密码</Dropdown.Item>
											<Dropdown.Item eventKey={2}>退出登陆</Dropdown.Item>
										</Dropdown.Menu>
									</Popover>
								}
							>
								<Avatar circle>{user?.name}</Avatar>
							</Whisper>
						</Nav>
					)}
				</Navbar>
			</Header>
			<Modal
				show={showUpdateUser}
				onHide={() => {
					setshowUpdateUser(false);
				}}
				backdrop="static"
			>
				<Modal.Header>
					<Modal.Title>修改用户密码</Modal.Title>
				</Modal.Header>
				<Modal.Body>
					<Form fluid onChange={onUpdateUser} formValue={updateUser}>
						<FormGroup>
							<ControlLabel>账号：</ControlLabel>
							<FormControl name="name" disabled="true" />
						</FormGroup>

						<FormGroup>
							<ControlLabel>密码：</ControlLabel>
							<FormControl name="passwd" />
						</FormGroup>
					</Form>
				</Modal.Body>
				<Modal.Footer>
					<Button
						onClick={() => APIUpdateUser(updateUser.id, updateUser.passwd)}
						style={{ color: '#FFF' }}
						appearance="primary"
						loading={updateUserLoading}
					>
						修改密码
					</Button>
				</Modal.Footer>
			</Modal>
		</>
	);
}
