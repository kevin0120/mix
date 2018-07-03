package fis

import (
	"fmt"
	"time"
)

const (
	FIS_VER        = "V001"
	FIS_TAG        = "ERGEBNIS"
	FIS_RESULT_VER = "002"
	FIS_RESULT_NUM = 1
	FIS_RESULT_KEY = "RESULT"
	FIS_VALUE_KEY  = "VALUE"
	FIS_RESULT_LEN = 182
	FIS_UNIT_NM    = "Nm"
	FIS_UNIT_DEG   = "DEG"
	FIS_ID_NM      = 1
	FIS_ID_DEG     = 2

	// 1. 设备名(4位）
	// 2. 协议版本（4位）
	// 3. 工厂代码（2位）
	// 4. 订单年份（4位）
	// 5. pin（7位）
	// 6. pin_check_code（1位）
	// 7. 装配流水线（2位）
	// 8. 序列号（6位）
	// 9. TAG（8位）
	// 10. 结果版本（3位）
	// 11. 结果总长（5位）
	// 12. 结果数量（3位）
	// 13. 结果Key（6位）
	// 14. 自定义结果id（20位，左对齐，右补空格）
	// 15. 系统类型（10位，左对齐，右补空格）
	// 16. 软件版本（6位，左对齐，右补空格）
	// 17. 拧接结果（4位，左对齐，右补下划线）
	// 18. 结果时间（14位）
	// 19. 操作模式（4位）
	// 20. 结果值key（6位，左对齐，右补空格）
	// 21. 结果数量（3位）
	// 22. 扭矩结果序号（6位）
	// 23. 扭矩结果值（20位，左对齐，右补空格）
	// 24. 扭矩单位（10位，左对齐，右补空格）
	// 25. 扭矩测量结果（1位）
	// 26. 角度结果序号（6位）
	// 27. 角度结果值（20位，左对齐，右补空格）
	// 28. 角度单位（10位，左对齐，右补空格）
	// 29. 角度测量结果（1位）
	//				  1   2   3  4  5  6  7 8     9 1011    12   13 14    15    16   17 18 19 20   21   22  23   24   25 26  27   28   29
	resultTemplate = "%s--%s--%s-%d-%d=%d-%s%06s**%s%s%05d##%03d*%s*%-20s*%-10s*%-6s*%s*%s*%s*%-6s*%03d*%06d%-20f%-10s%d*%06d%-20f%-10s%d***"
)

func (fr *FisResult) Serialize() string {

	resultTime := fr.Dat.Format("20060102150405")

	l := len(fr.Values)
	sFisResult := fmt.Sprintf(resultTemplate,
		fr.EquipemntName,
		FIS_VER,
		fr.FactoryName,
		fr.Year,
		fr.Pin,
		fr.PinCheckCode,
		fr.AssemblyLine,
		fr.Lnr,
		FIS_TAG,
		FIS_RESULT_VER,
		FIS_RESULT_LEN,
		FIS_RESULT_NUM,
		FIS_RESULT_KEY,
		fr.ResultID,
		fr.SystemType,
		fr.SoftwareVersion,
		fr.ResultValue,
		resultTime,
		fr.Mode,
		FIS_VALUE_KEY,
		l,
		fr.Values[0].ID,
		fr.Values[0].Value,
		fr.Values[0].Unit,
		fr.Values[0].Measure,
		fr.Values[1].ID,
		fr.Values[1].Value,
		fr.Values[1].Unit,
		fr.Values[1].Measure)

	return sFisResult
}

type FisResultValue struct {
	ID      int64   //结果值序号
	Value   float32 //结果值(扭矩/角度)
	Unit    string  //结果单位(Nm/DEG)
	Measure int     //检测结果(0: OK 1: NOK)
}

type FisResult struct {
	EquipemntName   string    //设备名
	Ver             string    //协议版本（V001）
	FactoryName     string    //工厂代码
	Year            int64     //订单年份
	Pin             int64     //车辆装配代码
	PinCheckCode    int64     //校验位
	AssemblyLine    string    //流水线
	Lnr             string    //流水号
	Seq             int64     //序列号
	Tag             string    //结果Tag（ERGEBNIS）
	ResultVer       string    //结果版本(002)
	ResultLen       int64     //结果长度（从Tag开始算）
	ResultNum       int64     //结果数量（001）
	ResultKey       string    //结果Key（RESULT）
	ResultID        string    //自定义的结果id
	SystemType      string    //系统类型
	SoftwareVersion string    //软件版本
	ResultValue     string    //拧接结果(IO/NIO)
	Dat             time.Time //时间日期（YYYYMMDDHHMMSS）
	Mode            string    //操作模式（AUTO/MANU）
	ValueKey        string    //结果key（VALUE）
	Values          []FisResultValue
}

func (fr *FisResult) Init() {
	fr.Ver = FIS_VER
	fr.Tag = FIS_TAG
	fr.ResultVer = FIS_RESULT_VER
	fr.ResultNum = FIS_RESULT_NUM
	fr.ResultKey = FIS_RESULT_KEY
	fr.ValueKey = FIS_VALUE_KEY
}
