package audi_vw

const (
	// 心跳包
	Xml_heart_beat = "<?xml version=\"1.0\" encoding=\"UTF-8\"?><root:ROOT xmlns:root=\"http://xmldefs.vw-group.com/KAP/station/V2.0/root\" xmlns:common=\"http://xmldefs.vw-group.com/KAP/station/V2.0/common\" xmlns:msl_msg=\"http://xmldefs.vw-group.com/KAP/station/V2.0/msl_msg\" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" xsi:schemaLocation=\"http://xmldefs.vw-group.com/KAP/station/V2.0/root \"><MSL_MSG><PNR>21</PNR></MSL_MSG></root:ROOT>"

	// 订阅数据(日期，时间)
	Xml_subscribe = "<?xml version=\"1.0\" encoding=\"UTF-8\"?><root:ROOT xmlns:root=\"http://xmldefs.vw-group.com/KAP/station/V2.0/root\" xmlns:common=\"http://xmldefs.vw-group.com/KAP/station/V2.0/common\" xmlns:msl_msg=\"http://xmldefs.vw-group.com/KAP/station/V2.0/msl_msg\" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" xsi:schemaLocation=\"http://xmldefs.vw-group.com/KAP/station/V2.0/root \"><MSL_MSG><PNR>21</PNR><LSN>0</LSN><SYN><DAT>%s</DAT><TIM>%s</TIM></SYN><SDR><SMG>1</SMG><SER>1</SER><SPC>1</SPC><SAL>1</SAL><SDS>1</SDS><SNU>1</SNU></SDR><PDC><CIO>5</CIO><CNO>5</CNO></PDC><UDR></UDR></MSL_MSG></root:ROOT>"

	// pset程序设定(日期，时间, sn, workorder_id, result_id, 拧紧次数, 用户id，程序号, knr, knr, knr)
	Xml_pset = "<?xml version=\"1.0\" encoding=\"UTF-8\"?><root:ROOT xmlns:root=\"http://xmldefs.vw-group.com/KAP/station/V2.0/root\" xmlns:common=\"http://xmldefs.vw-group.com/KAP/station/V2.0/common\" xmlns:msl_msg=\"http://xmldefs.vw-group.com/KAP/station/V2.0/msl_msg\" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" xsi:schemaLocation=\"http://xmldefs.vw-group.com/KAP/station/V2.0/root \"><MSL_MSG><PNR>21</PNR><SYN><DAT>%s</DAT><TIM>%s</TIM></SYN><PID><PRT>%s</PRT><PI1>%d</PI1><PI2>%d</PI2><STC>%d</STC></PID><SID><FAP>%d</FAP><FNR>FNR</FNR><COM>COM</COM><CNR>CNR</CNR><WID>WID</WID></SID><PRS><PRG>%d</PRG><TAP>TAP</TAP><TNR>TNR</TNR><SIO>1</SIO><MNO>0</MNO><NUT>0</NUT>%s</PRS><TOL><VAL>1</VAL>%s</TOL><STR><VAL>1</VAL>%s</STR></MSL_MSG></root:ROOT>"
)
