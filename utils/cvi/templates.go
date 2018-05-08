package cvi

const (
	// 心跳包
	Xml_heart_beat = "<?xml version=\"1.0\" encoding=\"UTF-8\"?><root:ROOT xmlns:root=\"http://xmldefs.vw-group.com/KAP/station/V2.0/root\" xmlns:common=\"http://xmldefs.vw-group.com/KAP/station/V2.0/common\" xmlns:msl_msg=\"http://xmldefs.vw-group.com/KAP/station/V2.0/msl_msg\" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" xsi:schemaLocation=\"http://xmldefs.vw-group.com/KAP/station/V2.0/root ..\\xsd\\root.xsd\"><MSL_MSG><PNR>21</PNR><RST/></MSL_MSG></root:ROOT>"

	// 订阅数据(日期，时间)
	Xml_subscribe = "<?xml version=\"1.0\" encoding=\"UTF-8\"?><root:ROOT xmlns:root=\"http://xmldefs.vw-group.com/KAP/station/V2.0/root\" xmlns:common=\"http://xmldefs.vw-group.com/KAP/station/V2.0/common\" xmlns:msl_msg=\"http://xmldefs.vw-group.com/KAP/station/V2.0/msl_msg\" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" xsi:schemaLocation=\"http://xmldefs.vw-group.com/KAP/station/V2.0/root ..\\xsd\\root.xsd\"><MSL_MSG><PNR>21</PNR><LSN>0</LSN><SYN><DAT>%s</DAT><TIM>%s</TIM></SYN><SDR><SMG>1</SMG><SER>1</SER><SPC>1</SPC><SAL>0</SAL><SDS>0</SDS><SNU>1</SNU></SDR><PDC><CIO>5</CIO><CNO>5</CNO></PDC><UDR><UCN/><UID/><UST/><UNT/></UDR></MSL_MSG></root:ROOT>"

	// pset程序设定(日期，时间, sn, workorder_id, 程序号)
	Xml_pset = "<?xml version=\"1.0\" encoding=\"UTF-8\"?><root:ROOT xmlns:root=\"http://xmldefs.vw-group.com/KAP/station/V2.0/root\" xmlns:common=\"http://xmldefs.vw-group.com/KAP/station/V2.0/common\" xmlns:msl_msg=\"http://xmldefs.vw-group.com/KAP/station/V2.0/msl_msg\" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" xsi:schemaLocation=\"http://xmldefs.vw-group.com/KAP/station/V2.0/root ..\\xsd\\root.xsd\"><MSL_MSG><PNR>21</PNR><SYN><DAT>%s</DAT><TIM>%s</TIM></SYN><PID><PRT>%s</PRT><PI1>%d</PI1><PI2>PI2</PI2><STC>STC</STC></PID><SID><FAP>FAP</FAP><FNR>FNR</FNR><COM>COM</COM><CNR>CNR</CNR><WID>WID</WID></SID><PRS><PRG>%d</PRG><TAP>TAP</TAP><TNR>TNR</TNR><SIO>1</SIO><MNO>0</MNO><NUT>0</NUT></PRS><TOL><VAL>1</VAL></TOL><STR><VAL>1</VAL></STR></MSL_MSG></root:ROOT>"
)
