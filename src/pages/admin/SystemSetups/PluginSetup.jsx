import React, { useState, useEffect } from 'react';
import {
	Panel,
	Placeholder,
	Form,
	FormGroup,
	FormControl,
	ControlLabel,
	HelpBlock,
	ButtonToolbar,
	Button,
	Row,
	Col,
	Avatar,
	Tag,
	Notification,
} from 'rsuite';

import useAxios from 'axios-hooks';
import axios from 'axios';

const { Paragraph } = Placeholder;

const OppositeMessageItem = ({ text, user = {}, style = {}, ...props }) => {
	if (!text || text == '') return null;

	return (
		<div
			style={{ display: 'flex', justifyContent: 'flex-start', paddingRight: 40, marginTop: 20, ...style }}
			{...props}
		>
			<div style={{ marginRight: 10 }}>
				<Avatar circle src={user?.avatar || 'https://q2.qlogo.cn/headimg_dl?spec=100&dst_uin=843369923'} />
			</div>
			<div>
				<div
					style={{
						display: 'flex',
						justifyContent: 'flex-start',
						textAlign: 'center',
						marginBottom: 5,
					}}
				>
					<p style={{ color: '#0006', fontSize: 12 }}>{user?.name || '无敌可爱的我'}</p>
				</div>
				<div>
					<div style={{ background: '#FFF', borderRadius: 8, padding: 10 }}>
						<p style={{ whiteSpace: 'pre-line' }}>{text}</p>
					</div>
				</div>
			</div>
		</div>
	);
};

export default function PluginSetup() {
	const [messageTemplate, setmessageTemplate] = useState({
		success: '',
		empty: '',
		fail: '',
		movie: {
			keywords: '罗小黑',
			total: 1,
			list: '1，《罗小黑战记》，立即观看：https://xiaocaihong.tv/movie/64596',
		},
		netLoding: false,
	});
	// const [messageTemplate, setmessageTemplate] = useState({
	// 	success:
	// 		'小彩虹视频（xiaocaihong.tv）帮您搜索到 ${movie.total} 条《${movie.keywords}》相关内容：\n\n${movie.list}',
	// 	empty: '小彩虹视频（xiaocaihong.tv）很遗憾暂时没有搜索到相关内容，资源马上就上线下载 APP 看看？立即下载：https://xiaocaihong.tv/app',
	// 	fail: '小彩虹视频（xiaocaihong.tv）好像不知道到您想要搜索的关键词，试试热门搜索：\n1，流浪地球\n2，你的名字\n3，我和我的祖国永远在一起\n\n下载APP高清资源无限免费看：https://xiaocaihong.tv/app',
	// 	movie: {
	// 		keywords: '罗小黑',
	// 		total: 1,
	// 		list: '1，《罗小黑战记》，立即观看：https://xiaocaihong.tv/movie/64596',
	// 	},
	// });

	const replaceTemplateCharacters = (template, data) => {
		let callback = template;
		const regex = /\$\{movie\.(.*?)\}/g;
		const replaceList = template.match(regex) || [];
		replaceList?.map((item, index) => {
			const itemR = item?.match(/\$\{movie\.(.*?)\}/) || [];
			if (itemR && itemR.length > 1) {
				const itemKey = itemR[1];
				callback = callback.replace(item, data[itemKey] || '');
			}
		});
		return callback;
	};

	// 获取消息模版数据
	const [{ data: messageTemData, loading: messageTemLoading, error: messageTemError }, messageTemRefetch] = useAxios({
		url: '/api/system/message_template/get',
	});
	useEffect(() => {
		if (messageTemData && messageTemData['data']) {
			const data = messageTemData['data'];
			setmessageTemplate({
				...messageTemplate,
				...data,
			});
		}
	}, [messageTemData]);

	// 提交保存配置
	const APISaveMessageTemplate = (success, empty, fail) => {
		if (!success || !empty || !fail) {
			Notification.error({
				title: '参数异常，消息模版不得为空！',
			});
			return;
		}

		// 请求中，不能重复发起请求
		if (messageTemplate.netLoding) {
			return;
		}

		setmessageTemplate({
			...messageTemplate,
			netLoding: true,
		});

		const params = new URLSearchParams();
		params.append('success', success);
		params.append('empty', empty);
		params.append('fail', fail);

		axios
			.post('/api/system/message_template/save', params, {})
			.then((res) => {
				setmessageTemplate({
					...messageTemplate,
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
					Notification.success({
						title: '保存配置成功',
					});
					messageTemRefetch();
				}
			})
			.catch((error) => {
				Notification.error({
					title: '失败，' + error || '失败，请稍后重试！',
				});
				setmessageTemplate({
					...messageTemplate,
					netLoding: false,
				});
			});
	};

	return (
		<Panel header="插件设定" bordered>
			<Row style={{ display: 'flex' }}>
				<Col style={{ flex: 1 }}>
					<Form
						layout="horizontal"
						formValue={messageTemplate}
						onChange={(value) => {
							setmessageTemplate({
								...messageTemplate,
								...value,
							});
						}}
					>
						<FormGroup>
							<ControlLabel>成功消息模版</ControlLabel>
							<FormControl name="success" rows={5} componentClass="textarea" />
							<HelpBlock>机器人搜索电影成功时的回复模版</HelpBlock>
						</FormGroup>
						<FormGroup>
							<ControlLabel>为空消息模版</ControlLabel>
							<FormControl name="empty" rows={5} componentClass="textarea" />
							<HelpBlock>机器人搜索电影为空时的回复模版</HelpBlock>
						</FormGroup>
						<FormGroup>
							<ControlLabel>失败消息模版</ControlLabel>
							<FormControl name="fail" rows={5} componentClass="textarea" />
							<HelpBlock>机器人搜索电影失败时的回复模版</HelpBlock>
						</FormGroup>
						<FormGroup>
							<ButtonToolbar>
								<Button
									appearance="primary"
									style={{ color: '#FFF' }}
									loading={messageTemplate?.netLoding}
									disabled={
										!messageTemplate?.success || !messageTemplate?.empty || !messageTemplate?.fail
									}
									onClick={() =>
										APISaveMessageTemplate(
											messageTemplate?.success,
											messageTemplate?.empty,
											messageTemplate?.fail
										)
									}
								>
									保存设定
								</Button>
							</ButtonToolbar>
						</FormGroup>
					</Form>
				</Col>
				<Col md={10}>
					<div style={{ background: '#f4f7fe', padding: '25px 15px', borderRadius: 8 }}>
						<div style={{ display: 'flex', justifyContent: 'flex-end', paddingLeft: 40 }}>
							<div>
								<div
									style={{
										display: 'flex',
										justifyContent: 'flex-end',
										textAlign: 'center',
										marginBottom: 5,
									}}
								>
									<p style={{ color: '#0006', fontSize: 12 }}>是天真呀</p>
								</div>
								<div>
									<div style={{ background: '#34b7ff', borderRadius: 8, padding: 10 }}>
										<p style={{ color: '#FFF', whiteSpace: 'pre-line' }}>搜索 罗小黑</p>
									</div>
								</div>
							</div>
							<div style={{ marginLeft: 10 }}>
								<Avatar circle src="https://haxibiao.com/images/avatar-2.jpg" />
							</div>
						</div>

						<OppositeMessageItem
							text={replaceTemplateCharacters(messageTemplate?.success, messageTemplate.movie || {})}
						/>
						<OppositeMessageItem
							text={replaceTemplateCharacters(messageTemplate?.empty, messageTemplate.movie || {})}
						/>
						<OppositeMessageItem
							text={replaceTemplateCharacters(messageTemplate?.fail, messageTemplate.movie || {})}
						/>
					</div>
					<HelpBlock style={{ marginTop: 5 }}>消息回复示例</HelpBlock>
				</Col>
			</Row>
		</Panel>
	);
}
