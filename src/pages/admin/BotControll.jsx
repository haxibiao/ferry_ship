import React from 'react';

import useAxios from 'axios-hooks';
import axios from 'axios';

import {
	Table,
	Panel,
	Loader,
	FlexboxGrid,
	Button,
	Modal,
	Notification,
	Form,
	FormControl,
	FormGroup,
	ControlLabel,
	Checkbox,
	Tag,
	Message,
	Grid,
	Row,
	Col,
	Input,
	InputGroup,
} from 'rsuite';
import BotLoginErrorView from './components/BotLoginErrorView';

const { useState } = React;
const { Column, HeaderCell, Cell, Pagination } = Table;

export default function BotControll() {
	const [{ data, loading, error }, refetch] = useAxios({
		url: '/api/account/list',
	});

	// 添加账号操作
	const [showAddAccount, setshowAddAccount] = useState(false);
	const [addAccountLoading, setaddAccountLoading] = useState(false);
	const [addAccount, setaddAccount] = useState({
		account: '',
		password: '',
		auto: 1,
	});
	const onAddAccount = (value) => {
		setaddAccount({
			...addAccount,
			...value,
		});
	};
	const APIAddAccount = () => {
		const { account, password, auto } = addAccount;

		if (!account || !password) {
			Notification.error({
				title: '账号密码不得为空！',
			});
			return;
		}

		setaddAccountLoading(true);

		const params = new URLSearchParams();
		params.append('account', account);
		params.append('password', password);
		params.append('auto', auto);

		axios
			.post('/api/account/add', params, {})
			.then((res) => {
				setaddAccountLoading(false);

				const { data } = res;
				const { code } = data;

				if (code < 1) {
					Notification.error({
						title: data?.msg || '添加失败，请稍后重试！',
					});
				} else {
					const bot = data?.data;

					Notification.success({
						title: `添加账号 ${bot.account} 成功！`,
					});

					// 创建用户成功，关闭弹窗，刷新列表数据，清空编辑框数据
					setshowAddAccount(false);
					setaddAccount({ account: '', password: '' });
					refetch();
				}
			})
			.catch((error) => {
				Notification.error({
					title: '添加失败，' + error || '添加失败，请稍后重试！',
				});
				setaddAccountLoading(false);
			});
	};

	// 用户列表数据处理
	let accounts = [];

	if (data) {
		let { data: data_array } = data;

		const da = new Date();
		const year = da.getFullYear() + '年';
		const month = da.getMonth() + 1 + '月';
		const date = da.getDate() + '日';
		const timeStr = [year, month, date].join('');

		data_array = data_array?.map((item, index) => {
			return {
				...item,
				time: timeStr,
			};
		});
		accounts = data_array;
	}

	// 禁用用户
	const APIDisableUser = (id, isDisable) => {
		const msgStrTitle = isDisable ? '禁用' : '启用';
		const params = new URLSearchParams();
		params.append('id', id);

		axios
			.post('/api/user/upstatus', params, {})
			.then((res) => {
				const { data } = res;
				const { code } = data;

				if (code < 1) {
					Notification.error({
						title: data?.msg || msgStrTitle + '失败，请稍后重试！',
					});
				} else {
					const user = data?.data;

					Notification.success({
						title: `${msgStrTitle}用户 ${user.name} 成功！`,
					});

					// 禁用用户成功，关闭弹窗，刷新列表数据，清空编辑框数据
					refetch();
				}
			})
			.catch((error) => {
				Notification.error({
					title: msgStrTitle + '失败，' + error || msgStrTitle + '失败，请稍后重试！',
				});
			});
	};

	// 登陆账号
	const [botLoginConfig, setbotLoginConfig] = useState({
		id: 0, // id
		account: 0, // 账号
		showModal: false, // 显示弹窗
		netLoding: false, // 标识请求中状态
		callback: {}, // 登陆回调数据
	});
	const APIBotLogin = () => {
		const { id } = botLoginConfig;

		if (!id) {
			Notification.error({
				title: 'id 异常',
			});
			return;
		}

		// 请求中，不能重复发起请求
		if (botLoginConfig.netLoding) {
			return;
		}

		setbotLoginConfig({
			...botLoginConfig,
			netLoding: true,
		});

		const params = new URLSearchParams();
		params.append('id', id);

		axios
			.post('/api/account/bot/login', params, {})
			.then((res) => {
				setbotLoginConfig({
					...botLoginConfig,
					netLoding: false,
				});

				const { data } = res;
				const { code } = data;

				if (code < 1) {
					Notification.error({
						title: data?.msg || '失败，请稍后重试！',
					});
				} else {
					const callback = data?.data;

					console.log('回调', callback);

					if (callback?.error) {
						// 需要二次认证
						setbotLoginConfig({
							...botLoginConfig,
							callback: callback,
						});
					} else {
						setbotLoginConfig({
							...botLoginConfig,
							id: 0,
							account: 0,
							showModal: false,
						});
						refetch();
					}
				}
			})
			.catch((error) => {
				Notification.error({
					title: '失败，' + error || '失败，请稍后重试！',
				});
				setbotLoginConfig({
					...botLoginConfig,
					netLoding: false,
				});
			});
	};

	// 退出登陆机器人账号
	const [botLogoutConfig, setbotLogoutConfig] = useState({
		id: 0, // id
		account: 0, // 账号
		showModal: false, // 显示弹窗
		netLoding: false, // 标识请求中状态
	});
	const APIBotLogout = () => {
		const { id } = botLogoutConfig;

		if (!id) {
			Notification.error({
				title: 'id 异常',
			});
			return;
		}

		// 请求中，不能重复发起请求
		if (botLogoutConfig.netLoding) {
			return;
		}

		setbotLogoutConfig({
			...botLogoutConfig,
			netLoding: true,
		});

		const params = new URLSearchParams();
		params.append('id', id);

		axios
			.post('/api/account/bot/logout', params, {})
			.then((res) => {
				setbotLogoutConfig({
					...botLogoutConfig,
					netLoding: false,
				});

				const { data } = res;
				const { code } = data;

				if (code < 1) {
					Notification.error({
						title: data?.msg || '失败，请稍后重试！',
					});
				} else {
					const callback = data?.data;

					console.log('回调', callback);

					if (callback?.error) {
						// 需要二次认证
						setbotLogoutConfig({
							...botLogoutConfig,
							callback: callback,
						});
					} else {
						setbotLogoutConfig({
							...botLogoutConfig,
							id: 0,
							account: 0,
							showModal: false,
						});
						refetch();
					}
				}
			})
			.catch((error) => {
				Notification.error({
					title: '失败，' + error || '失败，请稍后重试！',
				});
				setbotLogoutConfig({
					...botLogoutConfig,
					netLoding: false,
				});
			});
	};

	// 重设机器人账号密码
	const [repasswoedConfig, setrepasswoedConfig] = useState({
		id: 0, // id
		account: 0, // 账号
		password: '', // 密码
		showModal: false, // 显示弹窗
		netLoding: false, // 标识请求中状态
	});
	const onRepasswoedConfig = (value) => {
		setrepasswoedConfig({
			...repasswoedConfig,
			...value,
		});
	};
	const APIBotRepassword = () => {
		const { id, password } = repasswoedConfig;

		if (!id || !password) {
			Notification.error({
				title: '新密码不得为空！',
			});
			return;
		}

		// 请求中，不能重复发起请求
		if (repasswoedConfig.netLoding) {
			return;
		}

		setrepasswoedConfig({
			...repasswoedConfig,
			netLoding: true,
		});

		const params = new URLSearchParams();
		params.append('id', id);
		params.append('password', password);

		axios
			.post('/api/account/bot/repassword', params, {})
			.then((res) => {
				setrepasswoedConfig({
					...repasswoedConfig,
					netLoding: false,
				});

				const { data } = res;
				const { code } = data;

				if (code < 1) {
					Notification.error({
						title: '错误，请稍后重试！',
						description: data?.msg,
					});
				} else {
					const callback = data?.data;
					// console.log('回调', callback);

					Notification.success({
						title: '机器人密码重新设定成功！',
					});

					setrepasswoedConfig({
						id: 0,
						password: '',
						showModal: false,
					});
					refetch(); // 刷新列表数据
				}
			})
			.catch((error) => {
				Notification.error({
					title: '失败，请稍后重试！',
					description: error,
				});
				setrepasswoedConfig({
					...repasswoedConfig,
					netLoding: false,
				});
			});
	};

	if (loading) return <Loader backdrop content="loading..." vertical />;

	return (
		<div className="page_botControll" style={{ paddingTop: 25, paddingBottom: 25 }}>
			<p>机器人账号管理页面</p>
			<br />

			<FlexboxGrid style={{ marginBottom: 15 }} justify="end">
				<FlexboxGrid.Item>
					<Button appearance="ghost" onClick={() => setshowAddAccount(true)}>
						添加账号
					</Button>
				</FlexboxGrid.Item>
			</FlexboxGrid>

			<Panel bordered bodyFill>
				<Table
					autoHeight
					data={accounts}
					onRowClick={(data) => {
						console.log(data);
					}}
				>
					<Column width={70} align="center" fixed>
						<HeaderCell>ID</HeaderCell>
						<Cell dataKey="id" />
					</Column>
					<Column width={70} align="center" fixed>
						<HeaderCell>头像</HeaderCell>
						<Cell>
							{(rowData) => (
								<img
									style={{ width: 25, height: 25, borderRadius: 25 }}
									src={rowData?.avatar || 'https://haxibiao.com/images/avatar-2.jpg'}
								/>
							)}
						</Cell>
					</Column>

					<Column width={140} fixed>
						<HeaderCell>昵称</HeaderCell>
						<Cell dataKey="name" />
					</Column>

					<Column width={160}>
						<HeaderCell>账号</HeaderCell>
						<Cell dataKey="account" />
					</Column>

					<Column width={120}>
						<HeaderCell>状态</HeaderCell>

						<Cell>
							{(rowData) => (
								<span>
									{rowData.status == 1 ? <Tag color="green">在线</Tag> : <Tag color="red">离线</Tag>}
								</span>
							)}
						</Cell>
					</Column>

					<Column width={260}>
						<HeaderCell>登陆时间</HeaderCell>
						<Cell dataKey="updated" />
					</Column>

					<Column width={360} fixed="right">
						<HeaderCell>操作</HeaderCell>

						<Cell>
							{(rowData) => {
								// 登陆机器人账号
								function onBotLogin(e) {
									setbotLoginConfig({
										...botLoginConfig,
										id: rowData?.id,
										account: rowData?.account,
										showModal: true,
									});

									// 结束事件分发
									e.stopPropagation();
								}

								// 退出登陆机器人账号
								function onBotLogout(e) {
									setbotLogoutConfig({
										...botLogoutConfig,
										id: rowData?.id,
										account: rowData?.account,
										showModal: true,
									});

									// 结束事件分发
									e.stopPropagation();
								}

								// 重新设定机器人密码
								function onBotRepassword(e) {
									setrepasswoedConfig({
										...botLoginConfig,
										id: rowData?.id,
										account: rowData?.account,
										showModal: true,
									});
									// 结束事件分发
									e.stopPropagation();
								}

								function disableAction(e) {
									// 结束事件分发
									e.stopPropagation();
								}

								function onClick(e) {
									// 结束事件分发
									e.stopPropagation();
								}

								return (
									<span>
										{rowData.status == 1 ? (
											<>
												<a onClick={onBotLogout}> 退出登陆 </a> |
											</>
										) : (
											<>
												<a onClick={onBotLogin}> 立即登陆 </a> |
											</>
										)}
										<a onClick={onBotRepassword}> 重设密码 </a> |<a onClick={onClick}> 删除账号 </a>{' '}
										|
										<a onClick={disableAction}>
											{rowData.auto_login == 1 ? ' 禁用自动登陆 ' : ' 启用自动登陆 '}
										</a>{' '}
									</span>
								);
							}}
						</Cell>
					</Column>
				</Table>
			</Panel>

			{/* 添加机器人账号弹窗 */}
			<Modal
				show={showAddAccount}
				onHide={() => {
					setshowAddAccount(false);
				}}
				backdrop="static"
			>
				<Modal.Header>
					<Modal.Title>添加机器人账号</Modal.Title>
				</Modal.Header>
				<Modal.Body>
					<Form fluid onChange={onAddAccount} formValue={addAccount}>
						<FormGroup>
							<ControlLabel>QQ 账号：</ControlLabel>
							<FormControl name="account" />
						</FormGroup>

						<FormGroup>
							<ControlLabel>QQ 密码：</ControlLabel>
							<FormControl name="password" />
						</FormGroup>
					</Form>

					<Checkbox
						defaultChecked
						style={{ marginTop: 20 }}
						onChange={(value, checked, event) => {
							setaddAccount({
								...addAccount,
								auto: checked ? 1 : 0,
							});
						}}
					>
						启动系统时是否自动登陆
					</Checkbox>
				</Modal.Body>
				<Modal.Footer>
					<Button
						onClick={APIAddAccount}
						style={{ color: '#FFF' }}
						appearance="primary"
						loading={addAccountLoading}
					>
						添加账号
					</Button>
				</Modal.Footer>
			</Modal>

			{/* 重设机器人账号弹窗 */}
			<Modal
				show={repasswoedConfig.showModal}
				onHide={() => {
					setrepasswoedConfig({
						showModal: false,
					});
				}}
				backdrop="static"
			>
				<Modal.Header>
					<Modal.Title>重设机器人账号密码</Modal.Title>
				</Modal.Header>
				<Modal.Body>
					<p style={{ marginBottom: 20 }}>重新设定密码的机器人账号：{repasswoedConfig.account}</p>
					<Form fluid onChange={onRepasswoedConfig} formValue={repasswoedConfig}>
						<FormGroup>
							<ControlLabel>请输入新的密码：</ControlLabel>
							<FormControl name="password" />
						</FormGroup>
					</Form>
				</Modal.Body>
				<Modal.Footer>
					<Button
						onClick={APIBotRepassword}
						style={{ color: '#FFF' }}
						appearance="primary"
						loading={repasswoedConfig.netLoding}
					>
						重设密码
					</Button>
				</Modal.Footer>
			</Modal>

			{/* 立即登陆机器人账号弹窗 */}
			<Modal
				show={botLoginConfig.showModal}
				onHide={() => {
					setbotLoginConfig({
						showModal: false,
					});
				}}
				backdrop="static"
			>
				<Modal.Header>
					<Modal.Title>是否要登陆账号？</Modal.Title>
				</Modal.Header>
				<Modal.Body>
					<p>是否需要登陆机器人账号：{botLoginConfig.account}</p>
					<p>
						登陆成功之后会监听全部群消息，发现有人 @ 并包含 <code>搜索</code>{' '}
						关键词将会调用接口搜索电影并回复结果消息。
					</p>
					{botLoginConfig?.callback?.error ? (
						<BotLoginErrorView
							callbackData={botLoginConfig?.callback}
							onCallBack={(data) => {
								// 需要二次认证
								setbotLoginConfig({
									...botLoginConfig,
									callback: data,
								});
							}}
						/>
					) : null}
				</Modal.Body>
				<Modal.Footer>
					<Button
						onClick={APIBotLogin}
						style={{ color: '#FFF' }}
						appearance="primary"
						disabled={botLoginConfig?.callback?.error ? true : false}
						loading={botLoginConfig.netLoding}
					>
						立即登陆
					</Button>
				</Modal.Footer>
			</Modal>

			{/* 退出登陆机器人账号弹窗 */}
			<Modal
				show={botLogoutConfig.showModal}
				onHide={() => {
					setbotLogoutConfig({
						showModal: false,
					});
				}}
				backdrop="static"
			>
				<Modal.Header>
					<Modal.Title>是否要登陆账号？</Modal.Title>
				</Modal.Header>
				<Modal.Body>
					<p>是否需要退出登陆机器人账号：{botLogoutConfig.account}</p>
					<p>退出登陆之后需要重新登陆才能继续使用机器人搜索电影功能。</p>
				</Modal.Body>
				<Modal.Footer>
					<Button
						onClick={APIBotLogout}
						style={{ color: '#FFF' }}
						appearance="primary"
						disabled={botLogoutConfig?.callback?.error ? true : false}
						loading={botLogoutConfig.netLoding}
					>
						立即登陆
					</Button>
				</Modal.Footer>
			</Modal>
		</div>
	);
}
